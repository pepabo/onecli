package onelogin

import (
	o "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

const (
	// DefaultPageSize は1ページあたりのデフォルトの取得件数です
	DefaultPageSize = 1000
)

// OneloginClient represents the interface for Onelogin client operations
type OneloginClient interface {
	GetUsers(query models.Queryable) (interface{}, error)
	UpdateUser(userID int, user models.User) (interface{}, error)
	CreateUser(user models.User) (interface{}, error)
}

// Onelogin represents the Onelogin client
type Onelogin struct {
	client OneloginClient
}

// New creates a new Onelogin client
func New() (*Onelogin, error) {
	client, err := o.NewOneloginSDK()
	if err != nil {
		return nil, err
	}

	return &Onelogin{client: client}, nil
}
