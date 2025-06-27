package onelogin

import (
	"encoding/json"
	"strconv"
	"time"
)

// Event represents an OneLogin event
type Event struct {
	AccountID            int32      `json:"account_id,omitempty"`
	ActorSystem          string     `json:"actor_system,omitempty"`
	ActorUserID          int32      `json:"actor_user_id,omitempty"`
	ActorUserName        string     `json:"actor_user_name,omitempty"`
	AppID                int32      `json:"app_id,omitempty"`
	AppName              string     `json:"app_name,omitempty"`
	AssumingActingUserID int32      `json:"assuming_acting_user_id,omitempty"`
	BrowserFingerprint   string     `json:"browser_fingerprint,omitempty"`
	ClientID             string     `json:"client_id,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	CustomMessage        string     `json:"custom_message,omitempty"`
	DirectoryID          int32      `json:"directory_id,omitempty"`
	DirectorySyncRunID   int32      `json:"directory_sync_run_id,omitempty"`
	ErrorDescription     string     `json:"error_description,omitempty"`
	EventTypeID          int32      `json:"event_type_id,omitempty"`
	GroupID              int32      `json:"group_id,omitempty"`
	GroupName            string     `json:"group_name,omitempty"`
	ID                   uint64     `json:"id,omitempty"`
	IPAddr               string     `json:"ipaddr,omitempty"`
	Notes                string     `json:"notes,omitempty"`
	OperationName        string     `json:"operation_name,omitempty"`
	OTPDeviceID          int32      `json:"otp_device_id,omitempty"`
	OTPDeviceName        string     `json:"otp_device_name,omitempty"`
	PolicyID             int32      `json:"policy_id,omitempty"`
	PolicyName           string     `json:"policy_name,omitempty"`
	ProxyIP              string     `json:"proxy_ip,omitempty"`
	Resolution           int32      `json:"resolution,omitempty"`
	ResourceTypeID       int32      `json:"resource_type_id,omitempty"`
	RiskCookieID         string     `json:"risk_cookie_id,omitempty"`
	RiskReasons          string     `json:"risk_reasons,omitempty"`
	RiskScore            int32      `json:"risk_score,omitempty"`
	RoleID               int32      `json:"role_id,omitempty"`
	RoleName             string     `json:"role_name,omitempty"`
	Since                string     `json:"since,omitempty"`
	Until                string     `json:"until,omitempty"`
	UserID               int32      `json:"user_id,omitempty"`
	UserName             string     `json:"user_name,omitempty"`
}

type EventsResponse struct {
	Status struct {
		Error   bool   `json:"error"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"status"`
	Pagination struct {
		BeforeCursor *string `json:"before_cursor"`
		AfterCursor  *string `json:"after_cursor"`
		PreviousLink *string `json:"previous_link"`
		NextLink     *string `json:"next_link"`
	} `json:"pagination"`
	Data []Event `json:"data"`
}

// EventsQuery represents query parameters for events
type EventsQuery struct {
	Limit       string  `json:"limit,omitempty"`
	Cursor      string  `json:"after_cursor,omitempty"`
	ClientID    *string `json:"client_id,omitempty"`
	CreatedAt   *string `json:"created_at,omitempty"`
	DirectoryID *string `json:"directory_id,omitempty"`
	EventTypeID *string `json:"event_type_id,omitempty"`
	Resolution  *string `json:"resolution,omitempty"`
	ID          *string `json:"id,omitempty"`
	Since       *string `json:"since,omitempty"`
	Until       *string `json:"until,omitempty"`
	UserID      *string `json:"user_id,omitempty"`
}

// GetKeyValidators returns the validators for the query parameters
func (q EventsQuery) GetKeyValidators() map[string]func(any) bool {
	return map[string]func(any) bool{
		"limit":         validateString,
		"cursor":        validateString,
		"client_id":     validateString,
		"created_at":    validateString,
		"directory_id":  validateString,
		"event_type_id": validateString,
		"resolution":    validateString,
		"id":            validateString,
		"since":         validateString,
		"until":         validateString,
		"user_id":       validateString,
	}
}

// ListEvents retrieves events from OneLogin
func (o *Onelogin) ListEvents(query EventsQuery) ([]Event, error) {
	query.Limit = strconv.Itoa(DefaultPageSize)
	nextCursor := ""
	events := []Event{}

	for {
		if nextCursor != "" {
			query.Cursor = nextCursor
		}

		result, err := o.client.ListEvents(&query)
		if err != nil {
			return nil, err
		}

		// TODO: Inefficient workaround - using JSON marshaling/unmarshaling as a shortcut for complex type conversion
		response, err := convertToEventsResponse(result.(map[string]any))
		if err != nil {
			return nil, err
		}

		events = append(events, response.Data...)

		// Check if AfterCursor is nil before dereferencing
		if response.Pagination.AfterCursor == nil {
			break
		}
		nextCursor = *response.Pagination.AfterCursor

		if nextCursor == "" {
			break
		}
	}

	return events, nil
}

func convertToEventsResponse(data map[string]any) (*EventsResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var response EventsResponse
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
