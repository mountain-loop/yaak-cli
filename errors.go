package yaakcli

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

func CheckError(err error) {
	if err == nil {
		return
	}
	ExitError(err.Error())
}

func ExitError(msg string) {
	pterm.Println(pterm.Red(fmt.Sprintf("Error: %s", msg)))
	os.Exit(1)
}
