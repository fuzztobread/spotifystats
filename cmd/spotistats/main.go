package main

import (
	"os"

	"spotistats/cmd/spotistats/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
