package onelogin

import (
	"strconv"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/pepabo/onecli/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetApps(t *testing.T) {
	tests := []struct {
		name          string
		query         AppQuery
		mockResponse  []interface{}
		mockError     error
		expectedApps  []models.App
		expectedError error
	}{
		{
			name: "successful app retrieval with name query",
			query: AppQuery{
				Name: "Test App",
			},
			mockResponse: []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "Test App",
				},
			},
			expectedApps: []models.App{
				{
					ID:   func() *int32 { v := int32(1); return &v }(),
					Name: func() *string { v := "Test App"; return &v }(),
				},
			},
		},
		{
			name:  "successful app retrieval with multiple apps",
			query: AppQuery{},
			mockResponse: []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "Test App 1",
				},
				map[string]interface{}{
					"id":   float64(2),
					"name": "Test App 2",
				},
			},
			expectedApps: []models.App{
				{
					ID:   func() *int32 { v := int32(1); return &v }(),
					Name: func() *string { v := "Test App 1"; return &v }(),
				},
				{
					ID:   func() *int32 { v := int32(2); return &v }(),
					Name: func() *string { v := "Test App 2"; return &v }(),
				},
			},
		},
		{
			name: "successful app retrieval with empty result",
			query: AppQuery{
				Name: "Non-existent App",
			},
			mockResponse: []interface{}{},
			expectedApps: []models.App{},
		},
		{
			name: "error from client",
			query: AppQuery{
				Name: "Test App",
			},
			mockError:     assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{
				client: mockClient,
			}

			// Set up mock expectations
			expectedQuery := &models.AppQuery{
				Limit: strconv.Itoa(DefaultPageSize),
				Page:  "1",
			}
			if tt.query.Name != "" {
				expectedQuery.Name = &tt.query.Name
			}

			mockClient.On("GetApps", expectedQuery).Return(tt.mockResponse, tt.mockError)

			apps, err := o.GetApps(tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				if apps == nil {
					apps = []models.App{}
				}
				assert.Equal(t, tt.expectedApps, apps)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGetAppsDetails(t *testing.T) {
	tests := []struct {
		name              string
		query             AppQuery
		mockResponse      []interface{}
		mockUsersResponse []interface{}
		mockUsersError    error
		expectedApps      []AppDetails
		expectedError     error
	}{
		{
			name: "successful app retrieval with user details",
			query: AppQuery{
				Name: "Test App",
			},
			mockResponse: []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "Test App",
				},
			},
			mockUsersResponse: []interface{}{
				map[string]interface{}{
					"id":        float64(1),
					"email":     "user1@example.com",
					"username":  "user1",
					"firstname": "User",
					"lastname":  "One",
				},
				map[string]interface{}{
					"id":        float64(2),
					"email":     "user2@example.com",
					"username":  "user2",
					"firstname": "User",
					"lastname":  "Two",
				},
			},
			expectedApps: []AppDetails{
				{
					App: models.App{
						ID:   func() *int32 { v := int32(1); return &v }(),
						Name: func() *string { v := "Test App"; return &v }(),
					},
					Users: []models.User{
						{
							ID:        1,
							Email:     "user1@example.com",
							Username:  "user1",
							Firstname: "User",
							Lastname:  "One",
						},
						{
							ID:        2,
							Email:     "user2@example.com",
							Username:  "user2",
							Firstname: "User",
							Lastname:  "Two",
						},
					},
				},
			},
		},
		{
			name: "successful app retrieval with empty users",
			query: AppQuery{
				Name: "Test App",
			},
			mockResponse: []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "Test App",
				},
			},
			mockUsersResponse: []interface{}{},
			expectedApps: []AppDetails{
				{
					App: models.App{
						ID:   func() *int32 { v := int32(1); return &v }(),
						Name: func() *string { v := "Test App"; return &v }(),
					},
					Users: []models.User{},
				},
			},
		},
		{
			name: "app retrieval with user fetch error",
			query: AppQuery{
				Name: "Test App",
			},
			mockResponse: []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "Test App",
				},
			},
			mockUsersError: assert.AnError,
			expectedApps: []AppDetails{
				{
					App: models.App{
						ID:   func() *int32 { v := int32(1); return &v }(),
						Name: func() *string { v := "Test App"; return &v }(),
					},
					Users: []models.User{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{
				client: mockClient,
			}

			// Set up mock expectations
			expectedQuery := &models.AppQuery{
				Limit: strconv.Itoa(DefaultPageSize),
				Page:  "1",
			}
			if tt.query.Name != "" {
				expectedQuery.Name = &tt.query.Name
			}

			mockClient.On("GetApps", expectedQuery).Return(tt.mockResponse, nil)

			// Set up mock expectations for GetAppUsers if app has ID
			if len(tt.mockResponse) > 0 {
				if appData, ok := tt.mockResponse[0].(map[string]interface{}); ok {
					if appID, ok := appData["id"].(float64); ok {
						if tt.mockUsersError != nil {
							mockClient.On("GetAppUsers", int(appID)).Return(tt.mockUsersResponse, tt.mockUsersError)
						} else {
							mockClient.On("GetAppUsers", int(appID)).Return(tt.mockUsersResponse, nil)
						}
					}
				}
			}

			apps, err := o.GetAppsDetails(tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				if apps == nil {
					apps = []AppDetails{}
				}
				assert.Equal(t, tt.expectedApps, apps)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGetAppsWithPagination(t *testing.T) {
	tests := []struct {
		name          string
		query         AppQuery
		mockResponses [][]interface{}
		expectedApps  []models.App
		expectedError error
	}{
		{
			name:  "successful app retrieval with pagination",
			query: AppQuery{},
			mockResponses: [][]interface{}{
				func() []interface{} {
					result := make([]interface{}, DefaultPageSize)
					for i := 0; i < DefaultPageSize; i++ {
						result[i] = map[string]interface{}{
							"id":   float64(i + 1),
							"name": "Test App " + strconv.Itoa(i+1),
						}
					}
					return result
				}(),
				{
					map[string]interface{}{
						"id":   float64(1001),
						"name": "Test App 1001",
					},
					map[string]interface{}{
						"id":   float64(1002),
						"name": "Test App 1002",
					},
					map[string]interface{}{
						"id":   float64(1003),
						"name": "Test App 1003",
					},
				},
			},
			expectedApps: func() []models.App {
				result := make([]models.App, DefaultPageSize)
				for i := 0; i < DefaultPageSize; i++ {
					id := int32(i + 1)
					name := "Test App " + strconv.Itoa(i+1)
					result[i] = models.App{
						ID:   &id,
						Name: &name,
					}
				}
				id1 := int32(1001)
				name1 := "Test App 1001"
				id2 := int32(1002)
				name2 := "Test App 1002"
				id3 := int32(1003)
				name3 := "Test App 1003"
				result = append(result, models.App{ID: &id1, Name: &name1})
				result = append(result, models.App{ID: &id2, Name: &name2})
				result = append(result, models.App{ID: &id3, Name: &name3})
				return result
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{
				client: mockClient,
			}

			expectedQuery1 := &models.AppQuery{
				Limit: strconv.Itoa(DefaultPageSize),
				Page:  "1",
			}
			if tt.query.Name != "" {
				expectedQuery1.Name = &tt.query.Name
			}
			mockClient.On("GetApps", expectedQuery1).Return(tt.mockResponses[0], nil)

			expectedQuery2 := &models.AppQuery{
				Limit: strconv.Itoa(DefaultPageSize),
				Page:  "2",
			}
			if tt.query.Name != "" {
				expectedQuery2.Name = &tt.query.Name
			}
			mockClient.On("GetApps", expectedQuery2).Return(tt.mockResponses[1], nil)

			apps, err := o.GetApps(tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				if apps == nil {
					apps = []models.App{}
				}
				assert.Equal(t, tt.expectedApps, apps)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGetAppUsers(t *testing.T) {
	tests := []struct {
		name          string
		appID         int
		mockResponse  []interface{}
		mockError     error
		expectedUsers []models.User
		expectedError error
	}{
		{
			name:  "successful app users retrieval",
			appID: 123,
			mockResponse: []interface{}{
				map[string]interface{}{
					"id":        float64(1),
					"email":     "user1@example.com",
					"username":  "user1",
					"firstname": "User",
					"lastname":  "One",
				},
				map[string]interface{}{
					"id":        float64(2),
					"email":     "user2@example.com",
					"username":  "user2",
					"firstname": "User",
					"lastname":  "Two",
				},
			},
			expectedUsers: []models.User{
				{
					ID:        1,
					Email:     "user1@example.com",
					Username:  "user1",
					Firstname: "User",
					Lastname:  "One",
				},
				{
					ID:        2,
					Email:     "user2@example.com",
					Username:  "user2",
					Firstname: "User",
					Lastname:  "Two",
				},
			},
		},
		{
			name:          "successful app users retrieval with empty result",
			appID:         456,
			mockResponse:  []interface{}{},
			expectedUsers: []models.User{},
		},
		{
			name:          "error from client",
			appID:         789,
			mockError:     assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{
				client: mockClient,
			}

			mockClient.On("GetAppUsers", tt.appID).Return(tt.mockResponse, tt.mockError)

			users, err := o.GetAppUsers(tt.appID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				if users == nil {
					users = []models.User{}
				}
				assert.Equal(t, tt.expectedUsers, users)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
