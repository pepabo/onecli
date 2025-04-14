package cmd

import (
	"io"
	"log"

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

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(userCmd)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
