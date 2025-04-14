package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/pepabo/onecli/utils"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  `Commands for managing OneLogin users`,
}

var output string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `List all users in your OneLogin account`,
	Run: func(cmd *cobra.Command, args []string) {
		// OneLoginクライアントの初期化
		client, err := onelogin.NewOneloginSDK()
		if err != nil {
			fmt.Printf("Error initializing OneLogin client: %v\n", err)
			os.Exit(1)
		}

		// ページネーションを使用して全ユーザーを取得
		users, err := utils.Paginate(func(offset int) ([]models.User, error) {
			query := models.UserQuery{
				Limit: strconv.Itoa(100),
				Page:  strconv.Itoa(offset/100 + 1),
			}
			result, err := client.GetUsers(&query)
			if err != nil {
				return nil, err
			}

			// []interface{} を []models.User に変換
			interfaceSlice := result.([]interface{})
			return utils.ConvertToUsers(interfaceSlice)
		})
		if err != nil {
			fmt.Printf("Error fetching users: %v\n", err)
			os.Exit(1)
		}

		// 指定された形式で出力
		if err := utils.PrintOutput(users, utils.OutputFormat(output), os.Stdout); err != nil {
			fmt.Printf("Error printing output: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	userCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format (yaml, json)")
}
