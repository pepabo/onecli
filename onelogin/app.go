package onelogin

import (
	"strconv"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/pepabo/onecli/utils"
)

type AppQuery struct {
	Name string
}

// AppDetails represents an app with its associated details
type AppDetails struct {
	models.App `json:",inline"`
	Users      []models.User `json:"users,omitempty"`
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

// GetAppsDetails retrieves apps with user details from Onelogin
func (o *Onelogin) GetAppsDetails(query AppQuery) ([]AppDetails, error) {
	// GetAppsを内部的に呼び出してアプリ情報を取得
	apps, err := o.GetApps(query)
	if err != nil {
		return nil, err
	}

	// 各アプリに対してユーザー情報を取得
	appsWithDetails := make([]AppDetails, len(apps))
	for i, app := range apps {
		appDetails := AppDetails{
			App: app,
		}

		// アプリのIDが存在する場合のみユーザー情報を取得
		if app.ID != nil {
			users, err := o.GetAppUsers(int(*app.ID))
			if err != nil {
				// ユーザー情報の取得に失敗した場合は空のスライスを設定
				appDetails.Users = []models.User{}
			} else {
				appDetails.Users = users
			}
		}

		appsWithDetails[i] = appDetails
	}

	return appsWithDetails, nil
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
