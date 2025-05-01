package main

import (
	yaakcli "github.com/mountain-loop/yaak-cli"
)

var version = "dev"

func main() {
	yaakcli.Execute(version)
}
