package onelogin

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/pepabo/onecli/utils"
	"github.com/stretchr/testify/assert"
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
				Email: func() *string { v := "test@example.com"; return &v }(),
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
				Firstname: func() *string { v := "Test"; return &v }(),
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
				Email: func() *string { v := "test@example.com"; return &v }(),
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
				Limit: "",
				Page:  "1",
			}
			if tt.query.Email != nil {
				expectedQuery.Email = tt.query.Email
			}
			if tt.query.Username != nil {
				expectedQuery.Username = tt.query.Username
			}
			if tt.query.Firstname != nil {
				expectedQuery.Firstname = tt.query.Firstname
			}
			if tt.query.Lastname != nil {
				expectedQuery.Lastname = tt.query.Lastname
			}
			if tt.query.UserIDs != nil {
				expectedQuery.UserIDs = tt.query.UserIDs
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

func TestSendInviteLink(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		personalEmail string
		mockError     error
		expectedError error
	}{
		{
			name:  "successful invite to primary email",
			email: "user@example.com",
		},
		{
			name:          "successful invite with personal email",
			email:         "user@example.com",
			personalEmail: "user@gmail.com",
		},
		{
			name:          "error from client",
			email:         "user@example.com",
			mockError:     assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{client: mockClient}

			expectedInvite := models.Invite{
				Email:         tt.email,
				PersonalEmail: tt.personalEmail,
			}
			mockClient.On("SendInviteLink", expectedInvite).Return(nil, tt.mockError)

			err := o.SendInviteLink(tt.email, tt.personalEmail)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
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
		expectedID    int
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
			expectedID: 3,
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
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{
				client: mockClient,
			}

			mockClient.On("CreateUser", tt.inputUser).Return(tt.mockResponse, tt.mockError)

			id, err := o.CreateUser(tt.inputUser)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		user          User
		mockError     error
		expectedError error
	}{
		{
			name:   "successful update",
			userID: 1,
			user:   User{Email: "updated@example.com"},
		},
		{
			name:          "error from client",
			userID:        1,
			user:          User{Email: "updated@example.com"},
			mockError:     assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{client: mockClient}

			mockClient.On("UpdateUser", tt.userID, tt.user).Return(nil, tt.mockError)

			err := o.UpdateUser(tt.userID, tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestSetPassword(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		password      string
		mockError     error
		expectedError error
	}{
		{
			name:     "successful password set",
			userID:   1,
			password: "newpass123",
		},
		{
			name:          "error from client",
			userID:        1,
			password:      "newpass123",
			mockError:     assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(utils.MockClient)
			o := &Onelogin{client: mockClient}

			expectedBody := map[string]string{
				"password":              tt.password,
				"password_confirmation": tt.password,
			}
			mockClient.On("UpdatePasswordInsecure", tt.userID, expectedBody).Return(nil, tt.mockError)

			err := o.SetPassword(tt.userID, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}
