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

// GetEventTypes retrieves event types from OneLogin (with caching)
func (o *Onelogin) GetEventTypes() ([]EventType, error) {
	o.eventTypesCacheOnce.Do(func() {
		result, err := o.client.GetEventTypes(nil)
		if err != nil {
			o.eventTypesCacheErr = err
			return
		}

		// Convert the result to EventTypesResponse
		response, err := convertToEventTypesResponse(result.(map[string]any))
		if err != nil {
			o.eventTypesCacheErr = err
			return
		}

		o.eventTypesCache = response.Data
	})
	return o.eventTypesCache, o.eventTypesCacheErr
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

// EventTypeIDNameMap returns a map[id]name
func EventTypeIDNameMap(eventTypes []EventType) map[int32]string {
	m := make(map[int32]string)
	for _, et := range eventTypes {
		m[et.ID] = et.Name
	}
	return m
}

// EventTypeNameIDMap returns a map[name]id
func EventTypeNameIDMap(eventTypes []EventType) map[string]int32 {
	m := make(map[string]int32)
	for _, et := range eventTypes {
		m[et.Name] = et.ID
	}
	return m
}
