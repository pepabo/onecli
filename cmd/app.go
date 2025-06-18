package cmd

import (
	"fmt"
	"os"
	"strconv"

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
	appDetail      bool
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

		var result interface{}
		var err2 error

		if appDetail {
			result, err2 = client.GetAppsDetails(appQueryParams)
		} else {
			result, err2 = client.GetApps(appQueryParams)
		}

		if err2 != nil {
			return fmt.Errorf("error getting apps: %v", err2)
		}

		if err := utils.PrintOutput(result, utils.OutputFormat(appOutput), os.Stdout); err != nil {
			return fmt.Errorf("error printing output: %v", err)
		}
		return nil
	},
}

var appListUsersCmd = &cobra.Command{
	Use:          "list-users <app-id>",
	Aliases:      []string{"users", "lu"},
	Short:        "List users for a specific app",
	Long:         `List all users assigned to a specific app in your OneLogin organization`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid app ID: %v", err)
		}

		client, err := initClient()
		if err != nil {
			return err
		}

		users, err := client.GetAppUsers(appID)
		if err != nil {
			return fmt.Errorf("error getting app users: %v", err)
		}

		if err := utils.PrintOutput(users, utils.OutputFormat(appOutput), os.Stdout); err != nil {
			return fmt.Errorf("error printing output: %v", err)
		}
		return nil
	},
}

func init() {
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appListUsersCmd)

	appListCmd.Flags().StringVarP(&appOutput, "output", "o", "yaml", "Output format (yaml, json)")
	appListCmd.Flags().StringVar(&appQueryParams.Name, "name", "", "Filter apps by name")
	appListCmd.Flags().BoolVar(&appDetail, "detail", false, "Include user details for each app")

	appListUsersCmd.Flags().StringVarP(&appOutput, "output", "o", "yaml", "Output format (yaml, json)")
}
