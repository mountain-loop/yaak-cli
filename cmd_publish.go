package yaakcli

import (
	"archive/zip"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
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

		zipPipeReader, zipPipeWriter := io.Pipe()

		zipWriter := zip.NewWriter(zipPipeWriter)

		selected := make(map[string]bool)
		optionalFiles := []string{"README.md"}
		requiredFiles := []string{"package.json", "package-lock.json", "build/index.js"}
		for _, name := range optionalFiles {
			selected[filepath.Clean(name)] = true
		}

		for _, name := range requiredFiles {
			selected[filepath.Clean(name)] = true
			_, err := os.Stat(filepath.Join(pluginDir, name))
			if err != nil {
				pterm.Warning.Printf("Missing required file: %s\n", name)
				os.Exit(1)
			}
		}

		pterm.Info.Println("Creating zip file ", pluginDir)

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

				relPath = filepath.ToSlash(relPath) // Normalize for zip entries

				if !selected[relPath] {
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
			pterm.Info.Println("Zip file created")
		}()

		req := NewAPIRequest("POST", "/plugins/publish", zipPipeReader)
		body := SendAPIRequest(req)
		pterm.Success.Println("Plugin published")
		fmt.Printf("%s\n", body)
	},
}
