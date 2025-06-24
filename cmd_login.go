package yaakcli

import (
	"context"
	"errors"
	"fmt"
	"github.com/pkg/browser"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Yaak via web browser",
	Long:  "Open a web browser to authenticate with Yaak. Works with all browsers including Safari.",
	Run: func(cmd *cobra.Command, args []string) {
		CheckError(deleteAuthToken())
		baseURL := prodStagingDevStr("https://yaak.app", "https://todo.yaak.app", "http://localhost:9444")

		// Create a channel to receive the auth token
		tokenChan := make(chan string, 1)

		// Set up a simple HTTP server to handle the OAuth callback
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		CheckError(err)

		addr := listener.Addr().(*net.TCPAddr)
		redirectURL := fmt.Sprintf("http://127.0.0.1:%d/oauth/callback", addr.Port)

		// Open the browser to the login page
		oauthConfig := oauth2.Config{
			ClientID:     "yaak-cli",
			ClientSecret: "",
			Endpoint: oauth2.Endpoint{
				AuthURL:  baseURL + "/login/oauth/authorize",
				TokenURL: baseURL + "/login/oauth/access_token",
			},
			RedirectURL: redirectURL,
			Scopes:      nil,
		}
		loginURL, err := AuthorizationURL(&oauthConfig)
		CheckError(err)

		mux := http.NewServeMux()

		mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
			h := OAuthRedirectHandler{
				State:        loginURL.State,
				CodeVerifier: loginURL.CodeVerifier,
				OAuthConfig:  &oauthConfig,
			}
			token, err := h.ExchangeCode(r)
			// Get the token from the query parameters
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(w, "Failed to get access token: %s", err.Error())
				return
			}

			// Send the token to the channel
			tokenChan <- token

			// Return a success message to the browser
			redirectTo := baseURL + "/login/oauth/success"
			http.Redirect(w, r, redirectTo, http.StatusFound)
		})

		server := &http.Server{Handler: mux}

		// Start the server in a goroutine
		go func() {
			if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
				pterm.Error.Printf("HTTP server error: %v\n", err)
				os.Exit(1)
			}
		}()

		pterm.Info.Println("Initiating login to", loginURL)

		confirm := pterm.DefaultInteractiveConfirm
		confirm.DefaultValue = true
		open, err := confirm.Show("Open default browser")
		CheckError(err)

		if !open {
			os.Exit(0)
			return
		}

		// Set up a signal handler to gracefully shut down the server
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		// Open the browser based on the operating system
		err = browser.OpenURL(loginURL.String())
		if err != nil {
			pterm.Error.Printf("Failed to open browser: %v\n", err)
			pterm.Info.Println("Please open the following URL manually:")
			pterm.Info.Println(loginURL)
		}

		// Wait for either the token, a signal, or a timeout
		pterm.Info.Println("Waiting for authentication...")

		select {
		case token := <-tokenChan:
			pterm.Success.Println("Authentication successful!")

			// set password
			err = storeAuthToken(token)
			CheckError(err)

			// Shutdown the server
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				pterm.Error.Printf("Server shutdown error: %v\n", err)
			}

		case <-sigChan:
			pterm.Warning.Println("Interrupted by user. Shutting down...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				pterm.Error.Printf("Server shutdown error: %v\n", err)
			}

		case <-time.After(5 * time.Minute):
			pterm.Warning.Println("Timeout waiting for authentication. Shutting down...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				pterm.Error.Printf("Server shutdown error: %v\n", err)
			}
		}
	},
}
