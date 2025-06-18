package onelogin

import (
	"strconv"

	"github.com/pepabo/onecli/utils"
)

// AppDetails represents an app with its associated details
type AppDetails struct {
	App   `json:",inline"`
	Users []User `json:"users,omitempty"`
}

// GetApps retrieves apps from Onelogin
func (o *Onelogin) GetApps(query AppQuery) ([]App, error) {
	query.Limit = strconv.Itoa(DefaultPageSize)

	return utils.Paginate(func(page int) ([]App, error) {
		query.Page = strconv.Itoa(page)
		result, err := o.client.GetApps(&query)
		if err != nil {
			return nil, err
		}
		return utils.ConvertToApps(result.([]any))
	}, DefaultPageSize)
}

// GetAppsDetails retrieves apps with user details from Onelogin
func (o *Onelogin) GetAppsDetails(query AppQuery) ([]AppDetails, error) {
	apps, err := o.GetApps(query)
	if err != nil {
		return nil, err
	}

	appsWithDetails := make([]AppDetails, len(apps))
	for i, app := range apps {
		appDetails := AppDetails{
			App: app,
		}

		if app.ID != nil {
			users, err := o.GetAppUsers(int(*app.ID))
			if err != nil {
				appDetails.Users = []User{}
			} else {
				if users == nil {
					appDetails.Users = []User{}
				} else {
					appDetails.Users = users
				}
			}
		}

		appsWithDetails[i] = appDetails
	}

	return appsWithDetails, nil
}

// GetAppUsers retrieves users for a specific app from Onelogin
func (o *Onelogin) GetAppUsers(appID int) ([]User, error) {
	query := UserQuery{
		Limit: strconv.Itoa(DefaultPageSize),
	}
	return utils.Paginate(func(page int) ([]User, error) {
		query.Page = strconv.Itoa(page)
		result, err := o.client.GetAppUsers(appID, &query)
		if err != nil {
			return nil, err
		}
		return utils.ConvertToUsers(result.([]any))
	}, DefaultPageSize)
}
