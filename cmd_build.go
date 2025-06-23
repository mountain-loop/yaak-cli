package yaakcli

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var buildCmd = &cobra.Command{
	Use:   "build entrypoint",
	Short: "Transpile code into a runnable plugin bundle",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginDir, err := os.Getwd()
		CheckError(err)

		if len(args) > 0 {
			pluginDir, err = filepath.Abs(args[0])
			CheckError(err)
		}

		BuildPlugin(pluginDir)
	},
}
