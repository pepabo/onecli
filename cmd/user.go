package cmd

import (
	"fmt"
	"os"

	"github.com/pepabo/onecli/onelogin"
	"github.com/pepabo/onecli/utils"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  `Commands for managing OneLogin users in your organization`,
}

var (
	userQueryEmail     string
	userQueryUsername  string
	userQueryFirstname string
	userQueryLastname  string
	userQueryUserID    string
	output             string
)

// initClient initializes the OneLogin client
func initClient() (*onelogin.Onelogin, error) {
	client, err := onelogin.New()
	if err != nil {
		return nil, fmt.Errorf("error initializing OneLogin client: %v", err)
	}
	return client, nil
}

var listCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"l", "ls"},
	Short:        "List all users",
	Long:         `List all users in your OneLogin organization`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := initClient()
		if err != nil {
			return err
		}

		query := getUserQuery()
		users, err := client.GetUsers(query)
		if err != nil {
			return fmt.Errorf("error getting users: %v", err)
		}

		if err := utils.PrintOutput(users, utils.OutputFormat(output), os.Stdout); err != nil {
			return fmt.Errorf("error printing output: %v", err)
		}
		return nil
	},
}

var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify user information",
	Long:  `Modify user information in your OneLogin organization`,
}

var modifyEmailCmd = &cobra.Command{
	Use:          "email <new-email>",
	Aliases:      []string{"m", "mod"},
	Short:        "Modify user's email address",
	Long:         `Modify a user's email address in your OneLogin organization`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := getUserQuery()
		if isQueryParamsEmpty(query) {
			return fmt.Errorf("at least one query parameter (email, username, firstname, lastname, or user-id) must be specified")
		}

		newEmail := args[0]

		client, err := initClient()
		if err != nil {
			return err
		}

		users, err := client.GetUsers(query)
		if err != nil {
			return fmt.Errorf("error getting users: %v", err)
		}

		if len(users) == 0 {
			return fmt.Errorf("no users found matching the query")
		}

		if len(users) > 1 {
			return fmt.Errorf("multiple users found matching the query. Please be more specific")
		}

		user := users[0]
		user.Email = newEmail

		err = client.UpdateUser(int(user.ID), user)
		if err != nil {
			return fmt.Errorf("error updating user: %v", err)
		}

		fmt.Printf("Successfully updated user %s with new email: %s\n", user.Username, newEmail)
		return nil
	},
}

var addCmd = &cobra.Command{
	Use:   "add <first-name> <last-name> <email>",
	Short: "Add a new user",
	Long:  `Add a new user to your OneLogin organization`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		firstName := args[0]
		lastName := args[1]
		email := args[2]

		client, err := initClient()
		if err != nil {
			return fmt.Errorf("error initializing OneLogin client: %v", err)
		}

		newUser := onelogin.User{
			Firstname: firstName,
			Lastname:  lastName,
			Email:     email,
		}

		err = client.CreateUser(newUser)
		if err != nil {
			return fmt.Errorf("error creating user: %v", err)
		}

		fmt.Printf("Successfully added user: %s %s with email: %s\n", newUser.Firstname, newUser.Lastname, newUser.Email)
		return nil
	},
}

func getUserQuery() onelogin.UserQuery {
	return onelogin.UserQuery{
		Email:     &userQueryEmail,
		Username:  &userQueryUsername,
		Firstname: &userQueryFirstname,
		Lastname:  &userQueryLastname,
		UserIDs:   &userQueryUserID,
	}
}

// isQueryParamsEmpty checks if all query parameters are empty
func isQueryParamsEmpty(params onelogin.UserQuery) bool {
	return params.Email == nil && params.Username == nil && params.Firstname == nil && params.Lastname == nil && params.UserIDs == nil
}

func init() {
	userCmd.AddCommand(listCmd)
	userCmd.AddCommand(modifyCmd)
	modifyCmd.AddCommand(modifyEmailCmd)
	userCmd.AddCommand(addCmd)

	listCmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format (yaml, json)")
	listCmd.Flags().StringVar(&userQueryEmail, "email", "", "Filter users by email")
	listCmd.Flags().StringVar(&userQueryUsername, "username", "", "Filter users by username")
	listCmd.Flags().StringVar(&userQueryFirstname, "firstname", "", "Filter users by first name")
	listCmd.Flags().StringVar(&userQueryLastname, "lastname", "", "Filter users by last name")
	listCmd.Flags().StringVar(&userQueryUserID, "user-id", "", "Filter users by user ID")

	modifyEmailCmd.Flags().StringVar(&userQueryEmail, "email", "", "Query by email")
	modifyEmailCmd.Flags().StringVar(&userQueryUsername, "username", "", "Query by username")
	modifyEmailCmd.Flags().StringVar(&userQueryFirstname, "firstname", "", "Query by first name")
	modifyEmailCmd.Flags().StringVar(&userQueryLastname, "lastname", "", "Query by last name")
	modifyEmailCmd.Flags().StringVar(&userQueryUserID, "user-id", "", "Query by user ID")
}
