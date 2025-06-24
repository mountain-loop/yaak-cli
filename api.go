package yaakcli

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"net/http"
	"os"
	"time"
)

func NewAPIRequest(method, path string, body io.Reader) *http.Request {
	baseURL := prodStagingDevStr("https://api.yaak.app", "https://todo.yaak.app", "http://localhost:9444")
	req, err := http.NewRequest(method, baseURL+"/api/v1"+path, body)
	if err != nil {
		pterm.Error.Printf("Failed to create API request: %s\n", err)
		os.Exit(1)
	}
	return req
}

func SendAPIRequest(r *http.Request) []byte {
	found, token, err := getAuthToken()
	if err != nil {
		ExitError(err.Error())
	} else if !found {
		pterm.Warning.Println("Not logged in")
		pterm.Info.Println("Please run `yaakcli login`")
		os.Exit(1)
	}

	r.Header.Set("X-Yaak-Session", token)
	r.Header.Set("User-Agent", fmt.Sprintf("YaakCli/%s (%s)", CLIVersion, GetUAPlatform()))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(r)
	CheckError(err)
	defer func(Body io.ReadCloser) {
		CheckError(Body.Close())
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	CheckError(err)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		err = json.Unmarshal(body, &apiErr)
		if err != nil {
			pterm.Error.Printf("API Error %d → %s\n", resp.StatusCode, body)
		} else {
			pterm.Error.Println(apiErr.Message)
		}
		os.Exit(1)
	}

	message := resp.Header.Get("X-Cli-Message")
	if message != "" {
		pterm.Info.Printf(message)
	}

	return body
}

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
