package onelogin

import (
	"strconv"
	"testing"
	"time"

	"github.com/pepabo/onecli/utils"
	"github.com/stretchr/testify/assert"
)

func TestEventsQuery_GetKeyValidators(t *testing.T) {
	query := EventsQuery{}
	validators := query.GetKeyValidators()

	// Test that all expected validators are present
	expectedKeys := []string{"event_type_id", "user_id", "since", "until", "limit", "cursor"}
	for _, key := range expectedKeys {
		assert.Contains(t, validators, key, "Validator for key %s should be present", key)
	}

	// Test that validators are functions
	for key, validator := range validators {
		assert.NotNil(t, validator, "Validator for key %s should not be nil", key)
	}
}

func TestListEvents(t *testing.T) {
	tests := []struct {
		name           string
		query          EventsQuery
		mockResponse   []any
		expectedEvents []Event
		expectedError  error
	}{
		{
			name: "successful events retrieval",
			query: EventsQuery{
				EventTypeID: func() *string { v := "1"; return &v }(),
				UserID:      func() *string { v := "123"; return &v }(),
			},
			mockResponse: []any{
				map[string]any{
					"status": map[string]any{
						"error":   false,
						"code":    float64(200),
						"message": "Success",
						"type":    "success",
					},
					"pagination": map[string]any{
						"before_cursor": nil,
						"after_cursor":  nil,
						"previous_link": nil,
						"next_link":     nil,
					},
					"data": []any{
						map[string]any{
							"id":            float64(1),
							"account_id":    float64(12345),
							"event_type_id": float64(1),
							"user_id":       float64(123),
							"user_name":     "testuser",
							"created_at":    "2023-01-01T00:00:00Z",
							"ipaddr":        "192.168.1.1",
						},
						map[string]any{
							"id":            float64(2),
							"account_id":    float64(12345),
							"event_type_id": float64(2),
							"user_id":       float64(456),
							"user_name":     "testuser2",
							"created_at":    "2023-01-02T00:00:00Z",
							"ipaddr":        "192.168.1.2",
						},
					},
				},
			},
			expectedEvents: []Event{
				{
					ID:          1,
					AccountID:   12345,
					EventTypeID: 1,
					UserID:      123,
					UserName:    "testuser",
					IPAddr:      "192.168.1.1",
				},
				{
					ID:          2,
					AccountID:   12345,
					EventTypeID: 2,
					UserID:      456,
					UserName:    "testuser2",
					IPAddr:      "192.168.1.2",
				},
			},
		},
		{
			name: "successful events retrieval with pagination",
			query: EventsQuery{
				EventTypeID: func() *string { v := "1"; return &v }(),
			},
			mockResponse: []any{
				// First page
				map[string]any{
					"status": map[string]any{
						"error":   false,
						"code":    float64(200),
						"message": "Success",
						"type":    "success",
					},
					"pagination": map[string]any{
						"before_cursor": nil,
						"after_cursor":  func() *string { v := "cursor123"; return &v }(),
						"previous_link": nil,
						"next_link":     func() *string { v := "https://api.onelogin.com/events?after_cursor=cursor123"; return &v }(),
					},
					"data": []any{
						map[string]any{
							"id":            float64(1),
							"account_id":    float64(12345),
							"event_type_id": float64(1),
							"user_id":       float64(123),
							"user_name":     "testuser",
							"created_at":    "2023-01-01T00:00:00Z",
						},
					},
				},
				// Second page
				map[string]any{
					"status": map[string]any{
						"error":   false,
						"code":    float64(200),
						"message": "Success",
						"type":    "success",
					},
					"pagination": map[string]any{
						"before_cursor": func() *string { v := "cursor123"; return &v }(),
						"after_cursor":  nil,
						"previous_link": func() *string { v := "https://api.onelogin.com/events?before_cursor=cursor123"; return &v }(),
						"next_link":     nil,
					},
					"data": []any{
						map[string]any{
							"id":            float64(2),
							"account_id":    float64(12345),
							"event_type_id": float64(1),
							"user_id":       float64(456),
							"user_name":     "testuser2",
							"created_at":    "2023-01-02T00:00:00Z",
						},
					},
				},
			},
			expectedEvents: []Event{
				{
					ID:          1,
					AccountID:   12345,
					EventTypeID: 1,
					UserID:      123,
					UserName:    "testuser",
				},
				{
					ID:          2,
					AccountID:   12345,
					EventTypeID: 1,
					UserID:      456,
					UserName:    "testuser2",
				},
			},
		},
		{
			name: "successful events retrieval with empty result",
			query: EventsQuery{
				EventTypeID: func() *string { v := "999"; return &v }(),
			},
			mockResponse: []any{
				map[string]any{
					"status": map[string]any{
						"error":   false,
						"code":    float64(200),
						"message": "Success",
						"type":    "success",
					},
					"pagination": map[string]any{
						"before_cursor": nil,
						"after_cursor":  nil,
						"previous_link": nil,
						"next_link":     nil,
					},
					"data": []any{},
				},
			},
			expectedEvents: []Event{},
		},
		{
			name: "error from client",
			query: EventsQuery{
				EventTypeID: func() *string { v := "1"; return &v }(),
			},
			mockResponse:  []any{},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{
				client: mockClient,
			}

			// Set up mock expectations
			if tt.expectedError != nil {
				// Create expected query with Limit set to DefaultPageSize
				expectedQuery := tt.query
				expectedQuery.Limit = strconv.Itoa(DefaultPageSize)
				mockClient.On("ListEvents", &expectedQuery).Return(nil, tt.expectedError)
			} else {
				// Set up pagination expectations
				for i, response := range tt.mockResponse {
					if i == 0 {
						// Create expected query with Limit set to DefaultPageSize
						expectedQuery := tt.query
						expectedQuery.Limit = strconv.Itoa(DefaultPageSize)
						mockClient.On("ListEvents", &expectedQuery).Return(response, nil)
					} else {
						// For subsequent pages, we need to create a new query with cursor
						paginatedQuery := tt.query
						paginatedQuery.Limit = strconv.Itoa(DefaultPageSize)
						if i > 0 && i-1 < len(tt.mockResponse) {
							// Get the after_cursor from the previous response
							if prevResponse, ok := tt.mockResponse[i-1].(map[string]any); ok {
								if pagination, ok := prevResponse["pagination"].(map[string]any); ok {
									if afterCursor, ok := pagination["after_cursor"].(*string); ok && afterCursor != nil {
										paginatedQuery.Cursor = *afterCursor
									}
								}
							}
						}
						mockClient.On("ListEvents", &paginatedQuery).Return(response, nil)
					}
				}
			}

			events, err := o.ListEvents(tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				if events == nil {
					events = []Event{}
				}
				// Compare all fields except CreatedAt
				for i := range tt.expectedEvents {
					actual := events[i]
					expected := tt.expectedEvents[i]
					assert.Equal(t, expected.AccountID, actual.AccountID)
					assert.Equal(t, expected.ActorSystem, actual.ActorSystem)
					assert.Equal(t, expected.ActorUserID, actual.ActorUserID)
					assert.Equal(t, expected.ActorUserName, actual.ActorUserName)
					assert.Equal(t, expected.AppID, actual.AppID)
					assert.Equal(t, expected.AppName, actual.AppName)
					assert.Equal(t, expected.AssumingActingUserID, actual.AssumingActingUserID)
					assert.Equal(t, expected.BrowserFingerprint, actual.BrowserFingerprint)
					assert.Equal(t, expected.ClientID, actual.ClientID)
					assert.Equal(t, expected.CustomMessage, actual.CustomMessage)
					assert.Equal(t, expected.DirectoryID, actual.DirectoryID)
					assert.Equal(t, expected.DirectorySyncRunID, actual.DirectorySyncRunID)
					assert.Equal(t, expected.ErrorDescription, actual.ErrorDescription)
					assert.Equal(t, expected.EventTypeID, actual.EventTypeID)
					assert.Equal(t, expected.EventTypeIDs, actual.EventTypeIDs)
					assert.Equal(t, expected.GroupID, actual.GroupID)
					assert.Equal(t, expected.GroupName, actual.GroupName)
					assert.Equal(t, expected.ID, actual.ID)
					assert.Equal(t, expected.IPAddr, actual.IPAddr)
					assert.Equal(t, expected.Notes, actual.Notes)
					assert.Equal(t, expected.OperationName, actual.OperationName)
					assert.Equal(t, expected.OTPDeviceID, actual.OTPDeviceID)
					assert.Equal(t, expected.OTPDeviceName, actual.OTPDeviceName)
					assert.Equal(t, expected.PolicyID, actual.PolicyID)
					assert.Equal(t, expected.PolicyName, actual.PolicyName)
					assert.Equal(t, expected.ProxyIP, actual.ProxyIP)
					assert.Equal(t, expected.Resolution, actual.Resolution)
					assert.Equal(t, expected.ResourceTypeID, actual.ResourceTypeID)
					assert.Equal(t, expected.RiskCookieID, actual.RiskCookieID)
					assert.Equal(t, expected.RiskReasons, actual.RiskReasons)
					assert.Equal(t, expected.RiskScore, actual.RiskScore)
					assert.Equal(t, expected.RoleID, actual.RoleID)
					assert.Equal(t, expected.RoleName, actual.RoleName)
					assert.Equal(t, expected.Since, actual.Since)
					assert.Equal(t, expected.Until, actual.Until)
					assert.Equal(t, expected.UserID, actual.UserID)
					assert.Equal(t, expected.UserName, actual.UserName)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestListEventsWithComplexQuery(t *testing.T) {
	tests := []struct {
		name           string
		query          EventsQuery
		mockResponse   any
		expectedEvents []Event
		expectedError  error
	}{
		{
			name: "events with all query parameters",
			query: EventsQuery{
				EventTypeID: func() *string { v := "1"; return &v }(),
				UserID:      func() *string { v := "123"; return &v }(),
				Since:       func() *string { v := "2023-01-01T00:00:00Z"; return &v }(),
				Until:       func() *string { v := "2023-01-31T23:59:59Z"; return &v }(),
			},
			mockResponse: map[string]any{
				"status": map[string]any{
					"error":   false,
					"code":    float64(200),
					"message": "Success",
					"type":    "success",
				},
				"pagination": map[string]any{
					"before_cursor": nil,
					"after_cursor":  nil,
					"previous_link": nil,
					"next_link":     nil,
				},
				"data": []any{
					map[string]any{
						"id":            float64(1),
						"account_id":    float64(12345),
						"event_type_id": float64(1),
						"user_id":       float64(123),
						"app_id":        float64(456),
						"user_name":     "testuser",
						"app_name":      "Test App",
						"created_at":    "2023-01-15T12:00:00Z",
						"ipaddr":        "192.168.1.1",
						"risk_score":    float64(50),
					},
				},
			},
			expectedEvents: []Event{
				{
					ID:          1,
					AccountID:   12345,
					EventTypeID: 1,
					UserID:      123,
					AppID:       456,
					UserName:    "testuser",
					AppName:     "Test App",
					IPAddr:      "192.168.1.1",
					RiskScore:   50,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{
				client: mockClient,
			}

			// Set up mock expectations
			if tt.expectedError != nil {
				// Create expected query with Limit set to DefaultPageSize
				expectedQuery := tt.query
				expectedQuery.Limit = strconv.Itoa(DefaultPageSize)
				mockClient.On("ListEvents", &expectedQuery).Return(nil, tt.expectedError)
			} else {
				// Create expected query with Limit set to DefaultPageSize
				expectedQuery := tt.query
				expectedQuery.Limit = strconv.Itoa(DefaultPageSize)
				mockClient.On("ListEvents", &expectedQuery).Return(tt.mockResponse, nil)
			}

			events, err := o.ListEvents(tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				if events == nil {
					events = []Event{}
				}
				for i := range tt.expectedEvents {
					actual := events[i]
					expected := tt.expectedEvents[i]
					assert.Equal(t, expected.AccountID, actual.AccountID)
					assert.Equal(t, expected.AppID, actual.AppID)
					assert.Equal(t, expected.AppName, actual.AppName)
					assert.Equal(t, expected.EventTypeID, actual.EventTypeID)
					assert.Equal(t, expected.UserID, actual.UserID)
					assert.Equal(t, expected.UserName, actual.UserName)
					assert.Equal(t, expected.RiskScore, actual.RiskScore)
					assert.Equal(t, expected.ID, actual.ID)
					assert.Equal(t, expected.IPAddr, actual.IPAddr)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestConvertToEventsResponse(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]any
		expectedResult *EventsResponse
		expectedError  bool
	}{
		{
			name: "valid events response",
			input: map[string]any{
				"status": map[string]any{
					"error":   false,
					"code":    float64(200),
					"message": "Success",
					"type":    "success",
				},
				"pagination": map[string]any{
					"before_cursor": nil,
					"after_cursor":  func() *string { v := "cursor123"; return &v }(),
					"previous_link": nil,
					"next_link":     func() *string { v := "https://api.onelogin.com/events?after_cursor=cursor123"; return &v }(),
				},
				"data": []any{
					map[string]any{
						"id":            float64(1),
						"account_id":    float64(12345),
						"event_type_id": float64(1),
						"user_id":       float64(123),
						"user_name":     "testuser",
						"created_at":    "2023-01-01T00:00:00Z",
						"ipaddr":        "192.168.1.1",
					},
				},
			},
			expectedResult: &EventsResponse{
				Status: struct {
					Error   bool   `json:"error"`
					Code    int    `json:"code"`
					Message string `json:"message"`
					Type    string `json:"type"`
				}{
					Error:   false,
					Code:    200,
					Message: "Success",
					Type:    "success",
				},
				Pagination: struct {
					BeforeCursor *string `json:"before_cursor"`
					AfterCursor  *string `json:"after_cursor"`
					PreviousLink *string `json:"previous_link"`
					NextLink     *string `json:"next_link"`
				}{
					BeforeCursor: nil,
					AfterCursor:  func() *string { v := "cursor123"; return &v }(),
					PreviousLink: nil,
					NextLink:     func() *string { v := "https://api.onelogin.com/events?after_cursor=cursor123"; return &v }(),
				},
				Data: []Event{
					{
						ID:          1,
						AccountID:   12345,
						EventTypeID: 1,
						UserID:      123,
						UserName:    "testuser",
						IPAddr:      "192.168.1.1",
					},
				},
			},
			expectedError: false,
		},
		{
			name: "empty data response",
			input: map[string]any{
				"status": map[string]any{
					"error":   false,
					"code":    float64(200),
					"message": "Success",
					"type":    "success",
				},
				"pagination": map[string]any{
					"before_cursor": nil,
					"after_cursor":  nil,
					"previous_link": nil,
					"next_link":     nil,
				},
				"data": []any{},
			},
			expectedResult: &EventsResponse{
				Status: struct {
					Error   bool   `json:"error"`
					Code    int    `json:"code"`
					Message string `json:"message"`
					Type    string `json:"type"`
				}{
					Error:   false,
					Code:    200,
					Message: "Success",
					Type:    "success",
				},
				Pagination: struct {
					BeforeCursor *string `json:"before_cursor"`
					AfterCursor  *string `json:"after_cursor"`
					PreviousLink *string `json:"previous_link"`
					NextLink     *string `json:"next_link"`
				}{
					BeforeCursor: nil,
					AfterCursor:  nil,
					PreviousLink: nil,
					NextLink:     nil,
				},
				Data: []Event{},
			},
			expectedError: false,
		},
		{
			name: "invalid json data",
			input: map[string]any{
				"invalid": make(chan int), // This will cause JSON marshaling to fail
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertToEventsResponse(tt.input)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				for i := range tt.expectedResult.Data {
					actual := result.Data[i]
					expected := tt.expectedResult.Data[i]
					assert.Equal(t, expected.AccountID, actual.AccountID)
					assert.Equal(t, expected.AppID, actual.AppID)
					assert.Equal(t, expected.AppName, actual.AppName)
					assert.Equal(t, expected.EventTypeID, actual.EventTypeID)
					assert.Equal(t, expected.UserID, actual.UserID)
					assert.Equal(t, expected.UserName, actual.UserName)
					assert.Equal(t, expected.RiskScore, actual.RiskScore)
					assert.Equal(t, expected.ID, actual.ID)
					assert.Equal(t, expected.IPAddr, actual.IPAddr)
				}
			}
		})
	}
}

func TestEventStruct(t *testing.T) {
	// Test Event struct field mapping
	event := Event{
		AccountID:            12345,
		ActorSystem:          "system1",
		ActorUserID:          123,
		ActorUserName:        "actoruser",
		AppID:                456,
		AppName:              "Test App",
		AssumingActingUserID: 789,
		BrowserFingerprint:   "fingerprint123",
		ClientID:             "client123",
		CreatedAt:            func() *time.Time { t := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
		CustomMessage:        "Custom message",
		DirectoryID:          101,
		DirectorySyncRunID:   202,
		ErrorDescription:     "Error description",
		EventTypeID:          1,
		EventTypeIDs:         "1,2,3",
		GroupID:              303,
		GroupName:            "Test Group",
		ID:                   1,
		IPAddr:               "192.168.1.1",
		Notes:                "Test notes",
		OperationName:        "test_operation",
		OTPDeviceID:          404,
		OTPDeviceName:        "OTP Device",
		PolicyID:             505,
		PolicyName:           "Test Policy",
		ProxyIP:              "10.0.0.1",
		Resolution:           1,
		ResourceTypeID:       606,
		RiskCookieID:         "cookie123",
		RiskReasons:          "risk reasons",
		RiskScore:            50,
		RoleID:               707,
		RoleName:             "Test Role",
		Since:                "2023-01-01T00:00:00Z",
		Until:                "2023-01-31T23:59:59Z",
		UserID:               123,
		UserName:             "testuser",
	}

	// Test that all fields are properly set
	assert.Equal(t, int32(12345), event.AccountID)
	assert.Equal(t, "system1", event.ActorSystem)
	assert.Equal(t, int32(123), event.ActorUserID)
	assert.Equal(t, "actoruser", event.ActorUserName)
	assert.Equal(t, int32(456), event.AppID)
	assert.Equal(t, "Test App", event.AppName)
	assert.Equal(t, int32(789), event.AssumingActingUserID)
	assert.Equal(t, "fingerprint123", event.BrowserFingerprint)
	assert.Equal(t, "client123", event.ClientID)
	assert.NotNil(t, event.CreatedAt)
	assert.Equal(t, "Custom message", event.CustomMessage)
	assert.Equal(t, int32(101), event.DirectoryID)
	assert.Equal(t, int32(202), event.DirectorySyncRunID)
	assert.Equal(t, "Error description", event.ErrorDescription)
	assert.Equal(t, int32(1), event.EventTypeID)
	assert.Equal(t, "1,2,3", event.EventTypeIDs)
	assert.Equal(t, int32(303), event.GroupID)
	assert.Equal(t, "Test Group", event.GroupName)
	assert.Equal(t, uint64(1), event.ID)
	assert.Equal(t, "192.168.1.1", event.IPAddr)
	assert.Equal(t, "Test notes", event.Notes)
	assert.Equal(t, "test_operation", event.OperationName)
	assert.Equal(t, int32(404), event.OTPDeviceID)
	assert.Equal(t, "OTP Device", event.OTPDeviceName)
	assert.Equal(t, int32(505), event.PolicyID)
	assert.Equal(t, "Test Policy", event.PolicyName)
	assert.Equal(t, "10.0.0.1", event.ProxyIP)
	assert.Equal(t, int32(1), event.Resolution)
	assert.Equal(t, int32(606), event.ResourceTypeID)
	assert.Equal(t, "cookie123", event.RiskCookieID)
	assert.Equal(t, "risk reasons", event.RiskReasons)
	assert.Equal(t, int32(50), event.RiskScore)
	assert.Equal(t, int32(707), event.RoleID)
	assert.Equal(t, "Test Role", event.RoleName)
	assert.Equal(t, "2023-01-01T00:00:00Z", event.Since)
	assert.Equal(t, "2023-01-31T23:59:59Z", event.Until)
	assert.Equal(t, int32(123), event.UserID)
	assert.Equal(t, "testuser", event.UserName)
}
