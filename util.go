package yaakcli

import (
	"os"
	"path"
	"strings"
)

func writeFile(writePath, contents string) {
	CheckError(os.MkdirAll(path.Dir(writePath), 0755))
	CheckError(os.WriteFile(writePath, []byte(contents), 0755))
}

func readFile(path string) string {
	pkgBytes, err := TemplateFS.ReadFile(path)
	CheckError(err)
	return string(pkgBytes)
}

func copyFile(relPath, dstDir, name string) {
	contents := readFile(path.Join("template", relPath))
	contents = strings.ReplaceAll(contents, "yaak-plugin-name", name)
	writeFile(path.Join(dstDir, relPath), contents)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func prodStagingDevStr(prod, staging, dev string) string {
	if os.Getenv("ENVIRONMENT") == "staging" {
		return staging
	} else if os.Getenv("ENVIRONMENT") == "development" {
		return dev
	} else {
		return prod
	}
}
