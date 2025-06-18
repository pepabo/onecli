package main

import (
	"os"

	"github.com/pepabo/onecli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
