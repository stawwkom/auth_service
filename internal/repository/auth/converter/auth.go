package converter

import (
	"database/sql"
	"github.com/stawwkom/auth_service/internal/model"
	modelRepo "github.com/stawwkom/auth_service/internal/repository/auth/model"
	"time"
)

// Конвертируем Repo в Domain

func ToAuthFromRepo(auth *modelRepo.User) *model.User {
	var updatedAt *time.Time
	if auth.UpdatedAt.Valid {
		updatedAt = &auth.UpdatedAt.Time
	}

	return &model.User{
		ID:        auth.ID,
		Login:     auth.Login,
		Email:     auth.Email,
		Password:  auth.Password,
		CreatedAt: auth.CreatedAt,
		UpdatedAt: updatedAt,
	}
}

func ToRepoFromAuth(auth *model.User) *modelRepo.User {
	var updatedAt sql.NullTime
	if auth.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: *auth.UpdatedAt, Valid: true}
	}

	return &modelRepo.User{
		ID:        auth.ID,
		Login:     auth.Login,
		Email:     auth.Email,
		Password:  auth.Password,
		CreatedAt: auth.CreatedAt,
		UpdatedAt: updatedAt,
	}
}

func ToAuthInfoFromRepo(user *modelRepo.User) *model.UserInfo {

	return &model.UserInfo{
		ID:    user.ID,
		Login: user.Login,
		Email: user.Email,
	}
}
