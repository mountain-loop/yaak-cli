
# Yaak CLI (`yaakcli`)

This is the CLI for developing [Yaak](https://yaak.app) plugins.

## Installation

```shell
npm install -g @yaakapp/cli
```

## Commands

```
$ yaakcli --help

Generate, build, and debug plugins for Yaak, the most intuitive desktop API client

Usage:
  yaakcli [flags]
  yaakcli [command]

Available Commands:
  build       Transpile code into a runnable plugin bundle
  completion  Generate the autocompletion script for the specified shell
  dev         Build plugin bundle continuously when the filesystem changes
  generate    Generate a "Hello World" Yaak plugin
  help        Help about any command
  login       Login to Yaak via web browser
  logout      Sign out of the Yaak CLI
  publish     Publish a Yaak plugin version to the plugin registry
  whoami      Print the current logged-in user's info

Flags:
  -h, --help      help for yaakcli
      --version   Source directory to read from

Use "yaakcli [command] --help" for more information about a command.
```
