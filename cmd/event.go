package cmd

import (
	"fmt"
	"os"

	"github.com/pepabo/onecli/onelogin"
	"github.com/pepabo/onecli/utils"
	"github.com/spf13/cobra"
)

var eventCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"ev"},
	Short:   "Event management commands",
	Long:    `Commands for managing OneLogin events in your organization`,
}

var (
	eventQueryClientID    string
	eventQueryCreatedAt   string
	eventQueryDirectoryID string
	eventQueryEventTypeID string
	eventQueryResolution  string
	eventQueryID          string
	eventQuerySince       string
	eventQueryUntil       string
	eventQueryUserID      string
	eventOutput           string
)

var eventListCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"l", "ls"},
	Short:        "List all events",
	Long:         `List all events in your OneLogin organization`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := initClient()
		if err != nil {
			return err
		}

		query := getEventQuery()
		events, err := client.ListEvents(query)
		if err != nil {
			return fmt.Errorf("error getting events: %v", err)
		}

		if err := utils.PrintOutput(events, utils.OutputFormat(eventOutput), os.Stdout); err != nil {
			return fmt.Errorf("error printing output: %v", err)
		}
		return nil
	},
}

func getEventQuery() onelogin.EventsQuery {
	query := onelogin.EventsQuery{}

	if eventQueryClientID != "" {
		query.ClientID = &eventQueryClientID
	}

	if eventQueryCreatedAt != "" {
		query.CreatedAt = &eventQueryCreatedAt
	}

	if eventQueryDirectoryID != "" {
		query.DirectoryID = &eventQueryDirectoryID
	}

	if eventQueryEventTypeID != "" {
		query.EventTypeID = &eventQueryEventTypeID
	}

	if eventQueryResolution != "" {
		query.Resolution = &eventQueryResolution
	}

	if eventQueryID != "" {
		query.ID = &eventQueryID
	}

	if eventQuerySince != "" {
		query.Since = &eventQuerySince
	}

	if eventQueryUntil != "" {
		query.Until = &eventQueryUntil
	}

	if eventQueryUserID != "" {
		query.UserID = &eventQueryUserID
	}

	return query
}

func init() {
	eventCmd.AddCommand(eventListCmd)

	eventListCmd.Flags().StringVarP(&eventOutput, "output", "o", "yaml", "Output format (yaml, json)")
	eventListCmd.Flags().StringVar(&eventQueryClientID, "client-id", "", "Filter events by client ID")
	eventListCmd.Flags().StringVar(&eventQueryCreatedAt, "created-at", "", "Filter events by created at")
	eventListCmd.Flags().StringVar(&eventQueryDirectoryID, "directory-id", "", "Filter events by directory ID")
	eventListCmd.Flags().StringVar(&eventQueryEventTypeID, "event-type-id", "", "Filter events by event type ID")
	eventListCmd.Flags().StringVar(&eventQueryResolution, "resolution", "", "Filter events by resolution")
	eventListCmd.Flags().StringVar(&eventQueryID, "id", "", "Filter events by ID")
	eventListCmd.Flags().StringVar(&eventQuerySince, "since", "", "Filter events from date (YYYY-MM-DD)")
	eventListCmd.Flags().StringVar(&eventQueryUntil, "until", "", "Filter events to date (YYYY-MM-DD)")
	eventListCmd.Flags().StringVar(&eventQueryUserID, "user-id", "", "Filter events by user ID")
}
