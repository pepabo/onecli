package onelogin

import (
	"testing"

	"github.com/pepabo/onecli/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetEventTypes(t *testing.T) {
	tests := []struct {
		name               string
		mockResponse       any
		expectedEventTypes []EventType
		expectedError      error
	}{
		{
			name: "successful event types retrieval",
			mockResponse: map[string]any{
				"status": map[string]any{
					"error":   false,
					"code":    float64(200),
					"message": "Success",
					"type":    "success",
				},
				"data": []any{
					map[string]any{
						"id":          float64(1),
						"name":        "User Login",
						"description": "User login event",
					},
					map[string]any{
						"id":          float64(2),
						"name":        "User Logout",
						"description": "User logout event",
					},
					map[string]any{
						"id":          float64(3),
						"name":        "Password Reset",
						"description": "Password reset event",
					},
				},
			},
			expectedEventTypes: []EventType{
				{
					ID:          1,
					Name:        "User Login",
					Description: "User login event",
				},
				{
					ID:          2,
					Name:        "User Logout",
					Description: "User logout event",
				},
				{
					ID:          3,
					Name:        "Password Reset",
					Description: "Password reset event",
				},
			},
		},
		{
			name: "successful event types retrieval with empty result",
			mockResponse: map[string]any{
				"status": map[string]any{
					"error":   false,
					"code":    float64(200),
					"message": "Success",
					"type":    "success",
				},
				"data": []any{},
			},
			expectedEventTypes: []EventType{},
		},
		{
			name:          "error from client",
			mockResponse:  nil,
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
				mockClient.On("GetEventTypes", nil).Return(nil, tt.expectedError)
			} else {
				mockClient.On("GetEventTypes", nil).Return(tt.mockResponse, nil)
			}

			eventTypes, err := o.GetEventTypes()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				if eventTypes == nil {
					eventTypes = []EventType{}
				}
				assert.Equal(t, tt.expectedEventTypes, eventTypes)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestConvertToEventTypesResponse(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]any
		expectedResult *EventTypesResponse
		expectedError  bool
	}{
		{
			name: "valid event types response",
			input: map[string]any{
				"status": map[string]any{
					"error":   false,
					"code":    float64(200),
					"message": "Success",
					"type":    "success",
				},
				"data": []any{
					map[string]any{
						"id":          float64(1),
						"name":        "User Login",
						"description": "User login event",
					},
					map[string]any{
						"id":          float64(2),
						"name":        "User Logout",
						"description": "User logout event",
					},
				},
			},
			expectedResult: &EventTypesResponse{
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
				Data: []EventType{
					{
						ID:          1,
						Name:        "User Login",
						Description: "User login event",
					},
					{
						ID:          2,
						Name:        "User Logout",
						Description: "User logout event",
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
				"data": []any{},
			},
			expectedResult: &EventTypesResponse{
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
				Data: []EventType{},
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
			result, err := convertToEventTypesResponse(tt.input)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Status, result.Status)
				assert.Equal(t, tt.expectedResult.Data, result.Data)
			}
		})
	}
}

func TestEventTypeStruct(t *testing.T) {
	// Test EventType struct field mapping
	eventType := EventType{
		ID:          1,
		Name:        "Test Event Type",
		Description: "Test event type description",
	}

	// Test that all fields are properly set
	assert.Equal(t, int32(1), eventType.ID)
	assert.Equal(t, "Test Event Type", eventType.Name)
	assert.Equal(t, "Test event type description", eventType.Description)
}
