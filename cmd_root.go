package yaakcli

import (
	"github.com/spf13/cobra"
	"os"
)

var CLIVersion string

func rootCmd(v string) *cobra.Command {
	CLIVersion = v
	var fVersion bool
	cmd := &cobra.Command{
		Use:   "yaakcli",
		Short: "Develop plugins for Yaak",
		Long:  "Generate, build, and debug plugins for Yaak, the most intuitive desktop API client",
		Run: func(cmd *cobra.Command, args []string) {
			if fVersion {
				println(CLIVersion)
				os.Exit(0)
			}

			CheckError(cmd.Help())
		},
	}
	cmd.AddCommand(devCmd)
	cmd.AddCommand(buildCmd)
	cmd.AddCommand(generateCmd)
	cmd.AddCommand(whoamiCmd)
	cmd.AddCommand(loginCmd)
	cmd.AddCommand(logoutCmd)
	cmd.AddCommand(publishCmd)

	cmd.Flags().BoolVar(&fVersion, "version", false, "Show the current version of the Yaak CLI")

	return cmd
}

func Execute(version string) {
	CheckError(rootCmd(version).Execute())
}
