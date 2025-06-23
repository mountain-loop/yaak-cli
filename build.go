package yaakcli

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/pterm/pterm"
	"path/filepath"
)

func ESLintBuildOptions(pluginDir string) api.BuildOptions {
	srcPath := filepath.Join(pluginDir, "src", "index.ts")
	outPath := filepath.Join(pluginDir, "build", "index.js")
	return api.BuildOptions{
		EntryPoints: []string{srcPath},
		Outfile:     outPath,
		Platform:    api.PlatformNode,
		Bundle:      true, // Inline dependencies
		Write:       true, // Write to disk
		Format:      api.FormatCommonJS,
		LogLevel:    api.LogLevelInfo,
	}
}

func BuildPlugin(dir string) {
	if !fileExists(filepath.Join(dir, "package.json")) {
		ExitError("./package.json does not exist. Ensure that you are in a plugin directory.")
	}

	pterm.Info.Printf("Building plugin %s...\n", dir)

	api.Build(ESLintBuildOptions(dir))
}
