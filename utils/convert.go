package utils

import (
	"encoding/json"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

func ConvertToSlice[T any](data []any) ([]T, error) {
	result := make([]T, len(data))
	for i, v := range data {
		// map[string]interface{}をJSONに変換
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		// JSONを指定された型に変換
		var item T
		if err := json.Unmarshal(jsonData, &item); err != nil {
			return nil, err
		}

		result[i] = item
	}
	return result, nil
}

func ConvertToUsers(data []any) ([]models.User, error) {
	return ConvertToSlice[models.User](data)
}

func ConvertToApps(data []any) ([]models.App, error) {
	return ConvertToSlice[models.App](data)
}
