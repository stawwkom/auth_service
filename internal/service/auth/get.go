package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/stawwkom/auth_service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
)

// GetUser Получение полной информации по ID
func (s *service) GetUser(ctx context.Context, id int64) (*model.UserInfo, error) {
	user, err := s.authRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return &model.UserInfo{
		Login: user.Login,
		Email: user.Email,
	}, nil
}

// Login Аутентификация пользователя
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
