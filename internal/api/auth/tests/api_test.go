package tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stawwkom/auth_service/internal/api/auth"
	"github.com/stawwkom/auth_service/internal/converter"
	"github.com/stawwkom/auth_service/internal/model"
	"github.com/stawwkom/auth_service/internal/service/mocks"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	// Arrange
	mockService := mocks.NewAuthServiceMock(t)
	server := auth.NewServer(mockService)

	req := &desc.CreateUserRequest{
		Name:            "Test User",
		Email:           "test@example.com",
		Password:        "password123",
		PasswordConfirm: "password123",
		Role:            desc.Role_USER,
	}

	expectedID := int64(1)
	mockService.RegisterMock.When(context.Background(), mock.MatchedBy(func(user *model.User) bool {
		return user.Login == req.Name &&
			user.Email == req.Email &&
			user.Password == req.Password &&
			user.Role == int(req.Role)
	})).Then(expectedID, nil)

	// Act
	resp, err := server.Create(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedID, resp.Id)
	mockService.MinimockRegisterDone()
}

func TestCreate_Error(t *testing.T) {
	// Arrange
	mockService := mocks.NewAuthServiceMock(t)
	server := auth.NewServer(mockService)

	req := &desc.CreateUserRequest{
		Name:            "Test User",
		Email:           "test@example.com",
		Password:        "password123",
		PasswordConfirm: "password123",
		Role:            desc.Role_USER,
	}

	mockService.RegisterMock.When(context.Background(), mock.Anything).Then(int64(0), assert.AnError)

	// Act
	resp, err := server.Create(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, assert.AnError, err)
	mockService.MinimockRegisterDone()
}

func TestGet(t *testing.T) {
	// Arrange
	mockService := mocks.NewAuthServiceMock(t)
	server := auth.NewServer(mockService)

	userID := int64(1)
	req := &desc.GetUserRequest{
		Id: userID,
	}

	expectedUser := &model.UserInfo{
		Login: "test_user",
		Email: "test@example.com",
	}

	mockService.GetUserMock.When(context.Background(), userID).Then(expectedUser, nil)

	// Act
	resp, err := server.Get(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser.Login, resp.Name)
	assert.Equal(t, expectedUser.Email, resp.Email)
	mockService.MinimockGetUserDone()
}

func TestGet_Error(t *testing.T) {
	// Arrange
	mockService := mocks.NewAuthServiceMock(t)
	server := auth.NewServer(mockService)

	userID := int64(1)
	req := &desc.GetUserRequest{
		Id: userID,
	}

	mockService.GetUserMock.When(context.Background(), userID).Then(nil, assert.AnError)

	// Act
	resp, err := server.Get(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, assert.AnError, err)
	mockService.MinimockGetUserDone()
}

func TestGet_UserNotFound(t *testing.T) {
	// Arrange
	mockService := mocks.NewAuthServiceMock(t)
	server := auth.NewServer(mockService)

	userID := int64(999)
	req := &desc.GetUserRequest{
		Id: userID,
	}

	mockService.GetUserMock.When(context.Background(), userID).Then(nil, nil)

	// Act
	resp, err := server.Get(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.MinimockGetUserDone()
}
