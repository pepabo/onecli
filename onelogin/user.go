package onelogin

import (
	"fmt"
	"strconv"
	"time"

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
	if err != nil {
		return err
	}
	return nil
}

// SetUserState sets the user state to active and updates the last login time
func (o *Onelogin) SetUserState(userID int) error {
	user := User{
		Status:    1,
		LastLogin: time.Now(),
	}
	err := o.UpdateUser(userID, user)
	if err != nil {
		return fmt.Errorf("error setting user state: %v", err)
	}

	return nil
}

// CreateUser creates a new user in Onelogin
func (o *Onelogin) CreateUser(user User) error {
	result, err := o.client.CreateUser(user)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	if err := o.SetUserState(int(result.(User).ID)); err != nil {
		return err
	}

	return nil
}
