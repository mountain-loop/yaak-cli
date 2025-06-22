package yaakcli

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Print the current logged-in user's info",
	Run: func(cmd *cobra.Command, args []string) {
		req := NewAPIRequest("GET", "/whoami", nil)
		body := SendAPIRequest(req)
		pterm.Info.Printf("%s\n", body)
	},
}
