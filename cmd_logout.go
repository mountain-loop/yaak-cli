package yaakcli

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out of the Yaak CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteAuthToken()
		CheckError(err)
		pterm.Success.Println("Signed out of Yaak")
	},
}
