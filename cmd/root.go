package cmd

import (
	"fmt"
	"io"
	"log"

	"github.com/pepabo/onecli/version"
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "onecli",
	Short: "OneLogin CLI tool",
	Long:  `A CLI tool for interacting with OneLogin API`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// ログ出力の制御
		if !verbose {
			log.SetOutput(io.Discard)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of onecli",
	Long:  `Print the version number of onecli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", version.Version)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
