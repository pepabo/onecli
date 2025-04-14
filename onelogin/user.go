package onelogin

import (
	"strconv"

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
	q := models.UserQuery{}

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
		if q.Limit == "" {
			q.Limit = strconv.Itoa(DefaultPageSize)
		}
		q.Page = strconv.Itoa(page)
		result, err := o.client.GetUsers(&q)
		if err != nil {
			return nil, err
		}

		// []interface{} を []models.User に変換
		interfaceSlice := result.([]interface{})
		return utils.ConvertToUsers(interfaceSlice)
	}, DefaultPageSize)
}
