package onelogin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/pepabo/onecli/utils"
)

type UserQuery struct {
	Email     string
	Username  string
	Firstname string
	Lastname  string
	ID        string
}

// GetUsers retrieves users from Onelogin
func (o *Onelogin) GetUsers(query UserQuery) ([]models.User, error) {
	q := &models.UserQuery{
		Limit: strconv.Itoa(DefaultPageSize),
		Page:  "1",
	}

	if query.Email != "" {
		q.Email = &query.Email
	}
	if query.Username != "" {
		q.Username = &query.Username
	}
	if query.Firstname != "" {
		q.Firstname = &query.Firstname
	}
	if query.Lastname != "" {
		q.Lastname = &query.Lastname
	}
	if query.ID != "" {
		q.UserIDs = &query.ID
	}

	return utils.Paginate(func(page int) ([]models.User, error) {
		q.Page = strconv.Itoa(page)
		result, err := o.client.GetUsers(q)
		if err != nil {
			return nil, err
		}

		// []interface{} を []models.User に変換
		interfaceSlice := result.([]any)
		return utils.ConvertToSlice[models.User](interfaceSlice)
	}, DefaultPageSize)
}

// UpdateUser updates a user in Onelogin
func (o *Onelogin) UpdateUser(userID int, user models.User) (any, error) {
	return o.client.UpdateUser(userID, user)
}

// SetUserState sets the user state to active and updates the last login time
func (o *Onelogin) SetUserState(userID int) error {
	user := models.User{
		Status:    1, // 1はアクティブを示す
		LastLogin: time.Now(),
	}
	_, err := o.UpdateUser(userID, user)
	if err != nil {
		return fmt.Errorf("error setting user state: %v", err)
	}
	return nil
}

// CreateUser creates a new user in Onelogin
func (o *Onelogin) CreateUser(user models.User) (models.User, error) {
	createdUserInterface, err := o.client.CreateUser(user)
	if err != nil {
		return models.User{}, fmt.Errorf("error creating user: %v", err)
	}

	createdUserMap, ok := createdUserInterface.(map[string]any)
	if !ok {
		return models.User{}, fmt.Errorf("error asserting created user to map[string]interface{}")
	}

	createdUser := models.User{
		ID:        int32(createdUserMap["id"].(float64)),
		Email:     createdUserMap["email"].(string),
		Username:  createdUserMap["username"].(string),
		Firstname: createdUserMap["firstname"].(string),
		Lastname:  createdUserMap["lastname"].(string),
	}

	// ユーザーのステータスをアクティブに設定し、最終ログイン日時を更新
	if err := o.SetUserState(int(createdUser.ID)); err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}
