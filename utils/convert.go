package utils

import (
	"encoding/json"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// ConvertToUsers はmap[string]interface{}のスライスを[]models.Userに変換します
func ConvertToUsers(data []interface{}) ([]models.User, error) {
	users := make([]models.User, len(data))
	for i, v := range data {
		// map[string]interface{}をJSONに変換
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		// JSONをmodels.Userに変換
		var user models.User
		if err := json.Unmarshal(jsonData, &user); err != nil {
			return nil, err
		}

		users[i] = user
	}
	return users, nil
}
