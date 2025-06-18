package onelogin

import (
	"strconv"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/pepabo/onecli/utils"
)

type AppQuery struct {
	Name string
}

// GetApps retrieves apps from Onelogin
func (o *Onelogin) GetApps(query AppQuery) ([]models.App, error) {
	q := &models.AppQuery{
		Limit: strconv.Itoa(DefaultPageSize),
		Page:  "1",
	}

	if query.Name != "" {
		q.Name = &query.Name
	}

	return utils.Paginate(func(page int) ([]models.App, error) {
		q.Page = strconv.Itoa(page)
		result, err := o.client.GetApps(q)
		if err != nil {
			return nil, err
		}

		// []interface{} を []models.App に変換
		interfaceSlice := result.([]interface{})
		return utils.ConvertToSlice[models.App](interfaceSlice)
	}, DefaultPageSize)
}

// GetAppUsers retrieves users for a specific app from Onelogin
func (o *Onelogin) GetAppUsers(appID int) ([]models.User, error) {
	result, err := o.client.GetAppUsers(appID)
	if err != nil {
		return nil, err
	}

	// []interface{} を []models.User に変換
	interfaceSlice := result.([]interface{})
	return utils.ConvertToSlice[models.User](interfaceSlice)
}
