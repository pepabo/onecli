package cmd

import (
	"fmt"
	"os"
	"strings"

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
	eventQueryEventType   string
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
		query, err := getEventQuery(client)
		if err != nil {
			return err
		}
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

var eventTypesCmd = &cobra.Command{
	Use:          "types",
	Aliases:      []string{"t", "type"},
	Short:        "List all event types",
	Long:         `List all event types in your OneLogin organization`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := initClient()
		if err != nil {
			return err
		}

		eventTypes, err := client.GetEventTypes()
		if err != nil {
			return fmt.Errorf("error getting event types: %v", err)
		}

		if err := utils.PrintOutput(eventTypes, utils.OutputFormat(eventOutput), os.Stdout); err != nil {
			return fmt.Errorf("error printing output: %v", err)
		}
		return nil
	},
}

func getEventQuery(client *onelogin.Onelogin) (onelogin.EventsQuery, error) {
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

	// --typeと--type-idの排他チェック
	if eventQueryEventTypeID != "" && eventQueryEventType != "" {
		return query, fmt.Errorf("--type and --type-id cannot be used together")
	}

	if eventQueryEventTypeID != "" {
		query.EventTypeID = &eventQueryEventTypeID
	} else if eventQueryEventType != "" {
		eventTypes, err := client.GetEventTypes()
		if err != nil {
			return query, fmt.Errorf("error getting event types: %v", err)
		}
		nameToID := onelogin.EventTypeNameIDMap(eventTypes)
		typeNames := strings.Split(eventQueryEventType, ",")
		var typeIDs []string
		var invalidNames []string
		for _, name := range typeNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			if id, exists := nameToID[name]; exists {
				typeIDs = append(typeIDs, fmt.Sprintf("%d", id))
			} else {
				invalidNames = append(invalidNames, name)
			}
		}
		if len(invalidNames) > 0 {
			return query, fmt.Errorf("invalid event type name(s): %s. Use 'onecli event types' to see available event types", strings.Join(invalidNames, ", "))
		}
		if len(typeIDs) > 0 {
			eventTypeIDs := strings.Join(typeIDs, ",")
			query.EventTypeID = &eventTypeIDs
		}
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

	return query, nil
}

func init() {
	eventCmd.AddCommand(eventListCmd)
	eventCmd.AddCommand(eventTypesCmd)

	eventListCmd.Flags().StringVarP(&eventOutput, "output", "o", "yaml", "Output format (yaml, json)")
	eventListCmd.Flags().StringVar(&eventQueryClientID, "client-id", "", "Filter events by client ID")
	eventListCmd.Flags().StringVar(&eventQueryCreatedAt, "created-at", "", "Filter events by created at")
	eventListCmd.Flags().StringVar(&eventQueryDirectoryID, "directory-id", "", "Filter events by directory ID")
	eventListCmd.Flags().StringVar(&eventQueryEventTypeID, "type-id", "", "Filter events by event type ID (comma-separated for multiple values)")
	eventListCmd.Flags().StringVar(&eventQueryEventType, "type", "", "Filter events by event type name (comma-separated for multiple values)")
	eventListCmd.Flags().StringVar(&eventQueryResolution, "resolution", "", "Filter events by resolution")
	eventListCmd.Flags().StringVar(&eventQueryID, "id", "", "Filter events by ID")
	eventListCmd.Flags().StringVar(&eventQuerySince, "since", "", "Filter events from date (YYYY-MM-DD)")
	eventListCmd.Flags().StringVar(&eventQueryUntil, "until", "", "Filter events to date (YYYY-MM-DD)")
	eventListCmd.Flags().StringVar(&eventQueryUserID, "user-id", "", "Filter events by user ID")

	// Make --type and --type-id mutually exclusive
	eventListCmd.MarkFlagsMutuallyExclusive("type", "type-id")

	eventTypesCmd.Flags().StringVarP(&eventOutput, "output", "o", "yaml", "Output format (yaml, json)")
}
