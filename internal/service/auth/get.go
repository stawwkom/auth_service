package auth

import (
	"context"
	"errors"

	"github.com/stawwkom/auth_service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
)

// Получение полной информации по ID
func (s *service) GetUser(ctx context.Context, id int64) (*model.UserInfo, error) {
	return s.authRepo.Get(ctx, id)
}

// Аутентификация пользователя
func (s *service) Login(ctx context.Context, login string, password string) (*model.User, error) {
	user, err := s.authRepo.GetByLogin(ctx, login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
