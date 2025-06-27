package onelogin

import (
	"net/http"
	"os"
	"log/slog"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/api"
	o "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

type OneloginSDK struct {
	sdk *o.OneloginSDK
}

func NewOneloginSDKWrapper() (*OneloginSDK, error) {
	sdk, err := o.NewOneloginSDK()
	if err != nil {
		return nil, err
	}

	if os.Getenv("ONECLI_DEBUG") != "" {
		sdk.Client.HttpClient = NewHTTPDebuggerClient(sdk.Client.HttpClient)
	}

	return &OneloginSDK{sdk: sdk}, nil
}

func (s *OneloginSDK) GetUsers(query models.Queryable) (any, error) {
	return s.sdk.GetUsers(query)
}

func (s *OneloginSDK) UpdateUser(userID int, user models.User) (any, error) {
	return s.sdk.UpdateUser(userID, user)
}

func (s *OneloginSDK) CreateUser(user models.User) (any, error) {
	return s.sdk.CreateUser(user)
}

func (s *OneloginSDK) GetApps(query models.Queryable) (any, error) {
	return s.sdk.GetApps(query)
}

// Since GetAppUsers does not support pagination, we need to create a wrapper to support pagination.
// This wrapper will become unnecessary once onelogin-go-sdk supports it.
func (s *OneloginSDK) GetAppUsers(appID int, query models.Queryable) (any, error) {
	p, err := utl.BuildAPIPath(o.AppPath, appID, "users")
	if err != nil {
		return nil, err
	}

	return s.get(p, query)
}

func (s *OneloginSDK) ListEvents(query models.Queryable) (any, error) {
	return s.sdk.ListEvents(query)
}

func (s *OneloginSDK) get(path string, query models.Queryable) (any, error) {
	r, err := s.sdk.Client.Get(&path, query)
	if err != nil {
		return nil, err
	}

	return utl.CheckHTTPResponse(r)
}

type HTTPDebuggerClient struct {
	client api.HTTPClient
	logger *slog.Logger
}

func NewHTTPDebuggerClient(client api.HTTPClient) *HTTPDebuggerClient {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return &HTTPDebuggerClient{client: client, logger: logger}
}

func (c *HTTPDebuggerClient) Do(req *http.Request) (*http.Response, error) {
	method := req.Method
	url := req.URL.String()

	c.logger.Debug("Request", "method", method, "url", url)

	return c.client.Do(req)
}
