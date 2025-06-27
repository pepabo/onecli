package onelogin

import (
	"encoding/json"
)

// EventType represents an OneLogin event type
type EventType struct {
	ID          int32  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type EventTypesResponse struct {
	Status struct {
		Error   bool   `json:"error"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"status"`
	Data []EventType `json:"data"`
}

// GetEventTypes retrieves event types from OneLogin
func (o *Onelogin) GetEventTypes() ([]EventType, error) {
	result, err := o.client.GetEventTypes(nil)
	if err != nil {
		return nil, err
	}

	// Convert the result to EventTypesResponse
	response, err := convertToEventTypesResponse(result.(map[string]any))
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func convertToEventTypesResponse(data map[string]any) (*EventTypesResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var response EventTypesResponse
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
