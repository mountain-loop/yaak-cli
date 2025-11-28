package yaakcli

import (
	"os"
	"os/exec"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: `Generate a "Hello World" Yaak plugin`,
	Run: func(cmd *cobra.Command, args []string) {
		pluginName, err := pterm.DefaultInteractiveTextInput.WithDefaultText("Plugin name").WithDefaultValue(RandomName()).Show()
		CheckError(err)

		pluginDir, err := pterm.DefaultInteractiveTextInput.WithDefaultText("Plugin dir").WithDefaultValue("./" + pluginName).Show()
		CheckError(err)

		if fileExists(pluginDir) {
			returnError("")
		}

		pterm.Println("Generating plugin to:", pterm.Magenta(pluginDir))

		// Create destination directory
		CheckError(os.MkdirAll(pluginDir, 0755))

		// Copy static files
		copyFile(".gitignore", pluginDir, pluginName)
		copyFile("package.json", pluginDir, pluginName)
		copyFile("tsconfig.json", pluginDir, pluginName)
		copyFile("src/index.ts", pluginDir, pluginName)
		copyFile("src/index.test.ts", pluginDir, pluginName)

		primary := pterm.NewStyle(pterm.FgLightWhite, pterm.BgMagenta, pterm.Bold)

		pterm.DefaultHeader.WithBackgroundStyle(primary).Println("Installing npm dependencies...")
		runCmd(pluginDir, "npm", "install")
		runCmd(pluginDir, "npm", "install", "@yaakapp/api")
		runCmd(pluginDir, "npm", "install", "-D", "@yaakapp/cli")
		runCmd(pluginDir, "npm", "run", "build")
	},
}

func runCmd(dir, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	CheckError(c.Start())
	CheckError(c.Wait())
}

func returnError(msg string) {
	pterm.Println(pterm.Red(msg))
}
