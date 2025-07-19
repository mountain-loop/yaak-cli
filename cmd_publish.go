package yaakcli

import (
	"archive/zip"
	"encoding/json"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a Yaak plugin version to the plugin registry",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginDir, err := os.Getwd()
		CheckError(err)

		if len(args) > 0 {
			pluginDir, err = filepath.Abs(args[0])
			CheckError(err)
		}

		BuildPlugin(pluginDir)

		spinner, _ := pterm.DefaultSpinner.WithDelay(1 * time.Second).Start("Publishing plugin...")

		zipPipeReader, zipPipeWriter := io.Pipe()

		zipWriter := zip.NewWriter(zipPipeWriter)

		selected := make(map[string]bool)
		optionalFiles := []string{"package-lock.json"}
		requiredFiles := []string{"README.md", "package.json", "build/index.js", "src/index.ts"}
		for _, name := range optionalFiles {
			selected[filepath.ToSlash(filepath.Clean(name))] = true
		}

		for _, name := range requiredFiles {
			selected[filepath.ToSlash(filepath.Clean(name))] = true
			_, err := os.Stat(filepath.Join(pluginDir, name))
			if err != nil {
				pterm.Warning.Printf("Missing required file: %s\n", name)
				os.Exit(1)
			}
		}

		spinner.UpdateText("Archiving plugin")

		go func() {
			defer func() {
				CheckError(zipWriter.Close())
				CheckError(zipPipeWriter.Close())
			}()

			err = filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				relPath, err := filepath.Rel(pluginDir, path)
				if err != nil {
					return err
				}

				relPath = filepath.ToSlash(filepath.Clean(relPath)) // Normalize for zip entries

				// Skip non-desired files or files not in src/ or build/ (we want those)
				if !strings.HasPrefix(relPath, "src/") && !strings.HasPrefix(relPath, "build/") && !selected[relPath] {
					return nil
				}

				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer func(file *os.File) {
					err := file.Close()
					CheckError(err)
				}(file)

				writer, err := zipWriter.Create(relPath)
				if err != nil {
					return err
				}

				_, err = io.Copy(writer, file)

				return err
			})
			CheckError(err)
		}()

		spinner.UpdateText("Uploading plugin")
		req := NewAPIRequest("POST", "/plugins/publish", zipPipeReader)
		body := SendAPIRequest(req)

		var response struct {
			Version string `json:"version"`
			URL     string `json:"url"`
		}

		CheckError(json.Unmarshal(body, &response))
		spinner.Success("Plugin published ", response.Version, "\n → ", response.URL)
	},
}
