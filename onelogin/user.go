package onelogin

import (
	"fmt"
	"strconv"

	"github.com/pepabo/onecli/utils"
)

// GetUsers retrieves users from Onelogin
func (o *Onelogin) GetUsers(query UserQuery) ([]User, error) {
	return utils.Paginate(func(page int) ([]User, error) {
		query.Page = strconv.Itoa(page)
		result, err := o.client.GetUsers(&query)
		if err != nil {
			return nil, err
		}
		return utils.ConvertToUsers(result.([]any))
	}, DefaultPageSize)
}

// UpdateUser updates a user in Onelogin
func (o *Onelogin) UpdateUser(userID int, user User) error {
	_, err := o.client.UpdateUser(userID, user)
	return err
}

// CreateUser creates a new user in Onelogin and returns the created user's ID
func (o *Onelogin) CreateUser(user User) (int, error) {
	result, err := o.client.CreateUser(user)
	if err != nil {
		return 0, err
	}
	resultMap, ok := result.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("unexpected response type from create user: %T", result)
	}
	idFloat, ok := resultMap["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("missing or invalid id in create user response")
	}
	return int(idFloat), nil
}

// SetPassword sets a password for a user
func (o *Onelogin) SetPassword(userID int, password string) error {
	body := map[string]string{
		"password":              password,
		"password_confirmation": password,
	}
	_, err := o.client.UpdatePasswordInsecure(userID, body)
	return err
}

// SendInviteLink sends a password setup/reset invite link to a user
func (o *Onelogin) SendInviteLink(email, personalEmail string) error {
	invite := Invite{
		Email:         email,
		PersonalEmail: personalEmail,
	}
	_, err := o.client.SendInviteLink(invite)
	return err
}
