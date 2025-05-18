package repository

import (
	"context"
	"github.com/stawwkom/auth_service/internal/model"
)

type AuthRepository interface {
	Create(ctx context.Context, info *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.UserInfo, error)
	GetByLogin(ctx context.Context, login string) (*model.UserInfo, error)
	Update(ctx context.Context, id int64, info *model.User) error
	Delete(ctx context.Context, id int64) error
}
