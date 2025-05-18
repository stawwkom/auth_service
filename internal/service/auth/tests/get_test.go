package tests

import (
	"context"
	"testing"

	"github.com/stawwkom/auth_service/internal/model"
	"github.com/stawwkom/auth_service/internal/repository/mocks"
	"github.com/stawwkom/auth_service/internal/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAuthRepositoryMock(t)
	service := auth.NewAuthService(mockRepo)

	userID := int64(1)
	expectedUser := &model.UserInfo{
		Login: "test_user",
		Email: "test@example.com",
	}

	mockRepo.GetMock.When(context.Background(), userID).Then(expectedUser, nil)

	// Act
	user, err := service.GetUser(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Login, user.Login)
	assert.Equal(t, expectedUser.Email, user.Email)
	mockRepo.MinimockGetDone()
}

func TestGetUser_Error(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAuthRepositoryMock(t)
	service := auth.NewAuthService(mockRepo)

	userID := int64(1)
	mockRepo.GetMock.When(context.Background(), userID).Then(nil, assert.AnError)

	// Act
	user, err := service.GetUser(context.Background(), userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, assert.AnError, err)
	mockRepo.MinimockGetDone()
}

func TestGetUser_NotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAuthRepositoryMock(t)
	service := auth.NewAuthService(mockRepo)

	userID := int64(999)
	mockRepo.GetMock.When(context.Background(), userID).Then(nil, nil)

	// Act
	user, err := service.GetUser(context.Background(), userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.MinimockGetDone()
}
