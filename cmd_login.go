package yaakcli

import (
	"context"
	"fmt"
	"github.com/pkg/browser"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
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

		pterm.Info.Println("Starting browser-based login...")

		// Save the token to a config file
		confirm := pterm.DefaultInteractiveConfirm
		confirm.DefaultValue = true
		open, err := confirm.Show("Open default browser")
		CheckError(err)

		if !open {
			os.Exit(0)
			return
		}

		// Create a channel to receive the auth token
		tokenChan := make(chan string, 1)

		// Create a context that we can cancel
		_, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Set up a simple HTTP server to handle the OAuth callback
		server := &http.Server{
			Addr: "localhost:8085",
		}

		// Define the handler for the callback
		http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
			// Get the token from the query parameters
			token := r.URL.Query().Get("token")
			if token == "" {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(w, "Error: No token provided")
				return
			}

			// Send the token to the channel
			tokenChan <- token

			// Return a success message to the browser
			redirectTo := prodStagingDevStr(
				"https://yaak.app/login-cli/success",
				"https://todo.yaak.app/login-cli/success",
				"http://localhost:9444/login-cli/success",
			)
			http.Redirect(w, r, redirectTo, http.StatusFound)
		})

		// Start the server in a goroutine
		go func() {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				pterm.Error.Printf("HTTP server error: %v\n", err)
				os.Exit(1)
			}
		}()

		// Set up a signal handler to gracefully shut down the server
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		// Open the browser to the login page
		redirect := "http://localhost:8085/oauth/callback"
		loginURL := prodStagingDevStr(
			"https://yaak.app/login-cli?redirect=",
			"https://todo.yaak.app/login-cli?redirect=",
			"http://localhost:9444/login-cli?redirect=",
		) + url.QueryEscape(redirect)

		// Open the browser based on the operating system
		err = browser.OpenURL(loginURL)
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
