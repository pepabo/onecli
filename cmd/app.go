package cmd

import (
	"fmt"
	"os"

	"github.com/pepabo/onecli/onelogin"
	"github.com/pepabo/onecli/utils"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "App management commands",
	Long:  `Commands for managing OneLogin apps in your organization`,
}

var (
	appQueryParams onelogin.AppQuery
	appOutput      string
)

var appListCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"l", "ls"},
	Short:        "List all apps",
	Long:         `List all apps in your OneLogin organization`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := initClient()
		if err != nil {
			return err
		}

		apps, err := client.GetApps(appQueryParams)
		if err != nil {
			return fmt.Errorf("error getting apps: %v", err)
		}

		if err := utils.PrintOutput(apps, utils.OutputFormat(appOutput), os.Stdout); err != nil {
			return fmt.Errorf("error printing output: %v", err)
		}
		return nil
	},
}

func init() {
	appCmd.AddCommand(appListCmd)

	appListCmd.Flags().StringVarP(&appOutput, "output", "o", "yaml", "Output format (yaml, json)")
	appListCmd.Flags().StringVar(&appQueryParams.Name, "name", "", "Filter apps by name")
}
