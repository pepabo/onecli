package main

import (
	"fmt"
	"os"

	"github.com/pepabo/onecli/cmd"
)

var version string

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("onecli version %s\n", version)
		return
	}
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
