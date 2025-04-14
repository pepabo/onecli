package cmd

import (
	"fmt"
	"os"

	"github.com/pepabo/onecli/onelogin"
	"github.com/pepabo/onecli/utils"
	"github.com/spf13/cobra"
)

type userQueryParams struct {
	output    string
	email     string
	firstname string
	lastname  string
	userId    string
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  `Commands for managing OneLogin users in your organization`,
}

var queryParams userQueryParams

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all users",
	Long:    `List all users in your OneLogin organization`,
	Run: func(cmd *cobra.Command, args []string) {
		// OneLoginクライアントの初期化
		client, err := onelogin.New()
		if err != nil {
			fmt.Printf("Error initializing OneLogin client: %v\n", err)
			os.Exit(1)
		}

		// クエリパラメータの設定
		query := onelogin.UserQuery{}
		if queryParams.email != "" {
			query.Email = queryParams.email
		}
		if queryParams.firstname != "" {
			query.Firstname = queryParams.firstname
		}
		if queryParams.lastname != "" {
			query.Lastname = queryParams.lastname
		}
		if queryParams.userId != "" {
			query.ID = queryParams.userId
		}

		// ユーザーを取得
		users, err := client.GetUsers(query)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// 指定された形式で出力
		if err := utils.PrintOutput(users, utils.OutputFormat(queryParams.output), os.Stdout); err != nil {
			fmt.Printf("Error printing output: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	userCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&queryParams.output, "output", "o", "yaml", "Output format (yaml, json)")
	listCmd.Flags().StringVar(&queryParams.email, "email", "", "Filter users by email")
	listCmd.Flags().StringVar(&queryParams.firstname, "firstname", "", "Filter users by first name")
	listCmd.Flags().StringVar(&queryParams.lastname, "lastname", "", "Filter users by last name")
	listCmd.Flags().StringVar(&queryParams.userId, "user-id", "", "Filter users by user ID")
}
