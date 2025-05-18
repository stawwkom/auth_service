package tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stawwkom/auth_service/internal/model"
	"github.com/stawwkom/auth_service/internal/repository/mocks"
	serv "github.com/stawwkom/auth_service/internal/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAuthRepositoryMock(t)
	service := serv.NewAuthService(mockRepo)

	user := &model.User{
		Login:    "test_user",
		Email:    "test@example.com",
		Password: "password123",
		Role:     0,
	}

	expectedID := int64(1)
	mockRepo.CreateMock.When(context.Background(), mock.MatchedBy(func(u *model.User) bool {
		// Проверяем, что пароль был захеширован
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
		return err == nil && u.Login == user.Login && u.Email == user.Email && u.Role == user.Role
	})).Then(expectedID, nil)

	// Act
	id, err := service.Register(context.Background(), user)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedID, id)
	mockRepo.MinimockCreateDone()
}

func TestRegister_Error(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAuthRepositoryMock(t)
	service := serv.NewAuthService(mockRepo)

	user := &model.User{
		Login:    "test_user",
		Email:    "test@example.com",
		Password: "password123",
		Role:     0,
	}

	mockRepo.CreateMock.When(context.Background(), mock.Anything).Then(int64(0), assert.AnError)

	// Act
	id, err := service.Register(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, assert.AnError, err)
	mockRepo.MinimockCreateDone()
}

func TestRegister_HashError(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewAuthRepositoryMock(t)
	service := serv.NewAuthService(mockRepo)

	// Создаем пользователя с очень длинным паролем, который вызовет ошибку при хешировании
	user := &model.User{
		Login:    "test_user",
		Email:    "test@example.com",
		Password: string(make([]byte, 73)), // bcrypt имеет ограничение в 72 байта
		Role:     0,
	}

	// Act
	id, err := service.Register(context.Background(), user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	mockRepo.MinimockCreateDone()
}
