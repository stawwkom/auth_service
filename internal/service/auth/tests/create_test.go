package tests

import (
	"context"
	"errors"
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
	mockRepo.CreateMock.Set(func(ctx context.Context, u *model.User) (int64, error) {
		// Проверяем, что пароль захеширован
		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("password123")) != nil {
			return 0, errors.New("password not properly hashed")
		}

		if u.Login == "test_user" && u.Email == "test@example.com" && u.Role == 0 {
			return expectedID, nil
		}

		return 0, errors.New("unexpected user data")
	})

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

	mockRepo.CreateMock.Set(func(ctx context.Context, u *model.User) (int64, error) {
		return 0, assert.AnError
	})

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
