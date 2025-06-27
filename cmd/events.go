package cmd

import (
	"fmt"
	"os"

	"github.com/pepabo/onecli/onelogin"
	"github.com/pepabo/onecli/utils"
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:     "events",
	Aliases: []string{"ev"},
	Short:   "Events management commands",
	Long:    `Commands for managing OneLogin events in your organization`,
}

var (
	eventsQueryClientID    string
	eventsQueryCreatedAt   string
	eventsQueryDirectoryID string
	eventsQueryEventTypeID string
	eventsQueryResolution  string
	eventsQueryID          string
	eventsQuerySince       string
	eventsQueryUntil       string
	eventsQueryUserID      string
	eventsOutput           string
)

var eventsListCmd = &cobra.Command{
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

		query := getEventsQuery()
		events, err := client.ListEvents(query)
		if err != nil {
			return fmt.Errorf("error getting events: %v", err)
		}

		if err := utils.PrintOutput(events, utils.OutputFormat(eventsOutput), os.Stdout); err != nil {
			return fmt.Errorf("error printing output: %v", err)
		}
		return nil
	},
}

func getEventsQuery() onelogin.EventsQuery {
	query := onelogin.EventsQuery{}

	if eventsQueryClientID != "" {
		query.ClientID = &eventsQueryClientID
	}

	if eventsQueryCreatedAt != "" {
		query.CreatedAt = &eventsQueryCreatedAt
	}

	if eventsQueryDirectoryID != "" {
		query.DirectoryID = &eventsQueryDirectoryID
	}

	if eventsQueryEventTypeID != "" {
		query.EventTypeID = &eventsQueryEventTypeID
	}

	if eventsQueryResolution != "" {
		query.Resolution = &eventsQueryResolution
	}

	if eventsQueryID != "" {
		query.ID = &eventsQueryID
	}

	if eventsQuerySince != "" {
		query.Since = &eventsQuerySince
	}

	if eventsQueryUntil != "" {
		query.Until = &eventsQueryUntil
	}

	if eventsQueryUserID != "" {
		query.UserID = &eventsQueryUserID
	}

	return query
}

func init() {
	eventsCmd.AddCommand(eventsListCmd)

	eventsListCmd.Flags().StringVarP(&eventsOutput, "output", "o", "yaml", "Output format (yaml, json)")
	eventsListCmd.Flags().StringVar(&eventsQueryClientID, "client-id", "", "Filter events by client ID")
	eventsListCmd.Flags().StringVar(&eventsQueryCreatedAt, "created-at", "", "Filter events by created at")
	eventsListCmd.Flags().StringVar(&eventsQueryDirectoryID, "directory-id", "", "Filter events by directory ID")
	eventsListCmd.Flags().StringVar(&eventsQueryEventTypeID, "event-type-id", "", "Filter events by event type ID")
	eventsListCmd.Flags().StringVar(&eventsQueryResolution, "resolution", "", "Filter events by resolution")
	eventsListCmd.Flags().StringVar(&eventsQueryID, "id", "", "Filter events by ID")
	eventsListCmd.Flags().StringVar(&eventsQuerySince, "since", "", "Filter events from date (YYYY-MM-DD)")
	eventsListCmd.Flags().StringVar(&eventsQueryUntil, "until", "", "Filter events to date (YYYY-MM-DD)")
	eventsListCmd.Flags().StringVar(&eventsQueryUserID, "user-id", "", "Filter events by user ID")
}
