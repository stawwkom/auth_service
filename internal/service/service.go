package service

import (
	"context"
	"github.com/stawwkom/auth_service/internal/model"
)

// Интерфейс описывает, что умеет сервис
type AuthService interface {
	Register(ctx context.Context, user *model.User) (int64, error)
	Login(ctx context.Context, login, password string) (*model.User, error)
	GetUser(ctx context.Context, id int64) (*model.UserInfo, error)
	UpdateUser(ctx context.Context, id int64, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error
}
