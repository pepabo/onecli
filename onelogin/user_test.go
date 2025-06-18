package onelogin

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/pepabo/onecli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name          string
		query         UserQuery
		mockResponse  []any
		mockError     error
		expectedUsers []models.User
		expectedError error
	}{
		{
			name: "successful user retrieval with email query",
			query: UserQuery{
				Email: "test@example.com",
			},
			mockResponse: []any{
				map[string]any{
					"id":        1,
					"email":     "test@example.com",
					"username":  "testuser",
					"firstname": "Test",
					"lastname":  "User",
				},
			},
			expectedUsers: []models.User{
				{
					ID:        1,
					Email:     "test@example.com",
					Username:  "testuser",
					Firstname: "Test",
					Lastname:  "User",
				},
			},
		},
		{
			name: "successful user retrieval with multiple users",
			query: UserQuery{
				Firstname: "Test",
			},
			mockResponse: []any{
				map[string]any{
					"id":        1,
					"email":     "test1@example.com",
					"username":  "testuser1",
					"firstname": "Test",
					"lastname":  "User1",
				},
				map[string]any{
					"id":        2,
					"email":     "test2@example.com",
					"username":  "testuser2",
					"firstname": "Test",
					"lastname":  "User2",
				},
			},
			expectedUsers: []models.User{
				{
					ID:        1,
					Email:     "test1@example.com",
					Username:  "testuser1",
					Firstname: "Test",
					Lastname:  "User1",
				},
				{
					ID:        2,
					Email:     "test2@example.com",
					Username:  "testuser2",
					Firstname: "Test",
					Lastname:  "User2",
				},
			},
		},
		{
			name: "error from client",
			query: UserQuery{
				Email: "test@example.com",
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
			expectedQuery := &models.UserQuery{
				Limit: strconv.Itoa(DefaultPageSize),
				Page:  "1",
			}
			if tt.query.Email != "" {
				expectedQuery.Email = &tt.query.Email
			}
			if tt.query.Username != "" {
				expectedQuery.Username = &tt.query.Username
			}
			if tt.query.Firstname != "" {
				expectedQuery.Firstname = &tt.query.Firstname
			}
			if tt.query.Lastname != "" {
				expectedQuery.Lastname = &tt.query.Lastname
			}
			if tt.query.ID != "" {
				expectedQuery.UserIDs = &tt.query.ID
			}

			mockClient.On("GetUsers", expectedQuery).Return(tt.mockResponse, tt.mockError)

			users, err := o.GetUsers(tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUsers, users)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		inputUser     models.User
		mockResponse  any
		mockError     error
		expectedUser  models.User
		expectedError error
	}{
		{
			name: "successful user creation",
			inputUser: models.User{
				Email:     "newuser@example.com",
				Username:  "newuser",
				Firstname: "New",
				Lastname:  "User",
			},
			mockResponse: map[string]any{
				"id":        float64(3),
				"email":     "newuser@example.com",
				"username":  "newuser",
				"firstname": "New",
				"lastname":  "User",
			},
			expectedUser: models.User{
				ID:        3,
				Email:     "newuser@example.com",
				Username:  "newuser",
				Firstname: "New",
				Lastname:  "User",
			},
		},
		{
			name: "error creating user",
			inputUser: models.User{
				Email:     "erroruser@example.com",
				Username:  "erroruser",
				Firstname: "Error",
				Lastname:  "User",
			},
			mockError:     assert.AnError,
			expectedError: fmt.Errorf("error creating user: %v", assert.AnError),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{
				client: mockClient,
			}

			mockClient.On("CreateUser", tt.inputUser).Return(tt.mockResponse, tt.mockError)

			// Add expectation for UpdateUser call in SetUserState
			if tt.expectedError == nil {
				mockClient.On("UpdateUser", 3, mock.AnythingOfType("models.User")).Return(nil, nil)
			}

			createdUser, err := o.CreateUser(tt.inputUser)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, createdUser)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
