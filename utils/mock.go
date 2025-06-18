package utils

import (
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/mock"
)

// MockClient represents a mock implementation of the OneloginClient interface
type MockClient struct {
	mock.Mock
}

// GetUsers mocks the GetUsers method
func (m *MockClient) GetUsers(query models.Queryable) (any, error) {
	args := m.Called(query)
	return args.Get(0), args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *MockClient) UpdateUser(userID int, user models.User) (any, error) {
	args := m.Called(userID, user)
	return args.Get(0), args.Error(1)
}

// CreateUser mocks the CreateUser method
func (m *MockClient) CreateUser(user models.User) (any, error) {
	args := m.Called(user)
	return args.Get(0), args.Error(1)
}

// GetApps mocks the GetApps method
func (m *MockClient) GetApps(query models.Queryable) (any, error) {
	args := m.Called(query)
	return args.Get(0), args.Error(1)
}

// GetAppUsers mocks the GetAppUsers method
func (m *MockClient) GetAppUsers(appID int) (any, error) {
	args := m.Called(appID)
	return args.Get(0), args.Error(1)
}
