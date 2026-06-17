package onelogin

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	o "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/api"
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

// UpdateUser updates a user. It strips zero-value time.Time fields from the
// payload because models.User declares them as non-pointer time.Time, so
// json's omitempty cannot drop them and they would otherwise be sent as
// "0001-01-01T00:00:00Z" and overwrite OneLogin's server-managed timestamps.
func (s *OneloginSDK) UpdateUser(userID int, user models.User) (any, error) {
	body, err := userPayload(user)
	if err != nil {
		return nil, err
	}

	p, err := utl.BuildAPIPath(o.UserPathV2, userID)
	if err != nil {
		return nil, err
	}

	r, err := s.sdk.Client.Put(&p, body)
	if err != nil {
		return nil, err
	}

	return utl.CheckHTTPResponse(r)
}

// CreateUser creates a user. See UpdateUser for why zero-value time.Time
// fields are stripped from the payload.
func (s *OneloginSDK) CreateUser(user models.User) (any, error) {
	body, err := userPayload(user)
	if err != nil {
		return nil, err
	}

	p, err := utl.BuildAPIPath(o.UserPathV2)
	if err != nil {
		return nil, err
	}

	r, err := s.sdk.Client.Post(&p, body)
	if err != nil {
		return nil, err
	}

	return utl.CheckHTTPResponse(r)
}

// userPayload marshals a user and drops any field whose value is the zero
// time.Time ("0001-01-01T00:00:00Z"), which models.User emits for unset
// timestamps despite omitempty.
func userPayload(user models.User) (map[string]any, error) {
	b, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	const zeroTime = "0001-01-01T00:00:00Z"
	for k, v := range m {
		if s, ok := v.(string); ok && s == zeroTime {
			delete(m, k)
		}
	}

	return m, nil
}

func (s *OneloginSDK) UpdatePasswordInsecure(userID int, requestBody any) (any, error) {
	return s.sdk.UpdatePasswordInsecure(userID, requestBody)
}

func (s *OneloginSDK) SendInviteLink(invite models.Invite) (any, error) {
	return s.sdk.SendInviteLink(invite)
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

func (s *OneloginSDK) GetEventTypes(query models.Queryable) (any, error) {
	return s.sdk.GetEventTypes(query)
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
