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
	eventsQueryEventType string
	eventsQueryUserID    string
	eventsQueryAppID     string
	eventsQueryFrom      string
	eventsQueryTo        string
	eventsOutput         string
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

	if eventsQueryUserID != "" {
		query.UserID = &eventsQueryUserID
	}

	return query
}

func init() {
	eventsCmd.AddCommand(eventsListCmd)

	eventsListCmd.Flags().StringVarP(&eventsOutput, "output", "o", "yaml", "Output format (yaml, json)")
	eventsListCmd.Flags().StringVar(&eventsQueryEventType, "event-type", "", "Filter events by event type")
	eventsListCmd.Flags().StringVar(&eventsQueryUserID, "user-id", "", "Filter events by user ID")
	eventsListCmd.Flags().StringVar(&eventsQueryAppID, "app-id", "", "Filter events by app ID")
	eventsListCmd.Flags().StringVar(&eventsQueryFrom, "from", "", "Filter events from date (YYYY-MM-DD)")
	eventsListCmd.Flags().StringVar(&eventsQueryTo, "to", "", "Filter events to date (YYYY-MM-DD)")
}
