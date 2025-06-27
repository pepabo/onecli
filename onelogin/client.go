package onelogin

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

const (
	DefaultPageSize = 1000
)

type (
	User      = models.User
	UserQuery = models.UserQuery
)

type (
	App      = models.App
	AppQuery = models.AppQuery
)

type Client interface {
	GetUsers(query models.Queryable) (any, error)
	UpdateUser(userID int, user models.User) (any, error)
	CreateUser(user models.User) (any, error)
	GetApps(query models.Queryable) (any, error)
	GetAppUsers(appID int, query models.Queryable) (any, error)
	ListEvents(query models.Queryable) (any, error)
}

type Onelogin struct {
	client Client
}

// New creates a new Onelogin client
func New() (*Onelogin, error) {
	client, err := NewOneloginSDKWrapper()
	if err != nil {
		return nil, err
	}

	return &Onelogin{client: client}, nil
}
