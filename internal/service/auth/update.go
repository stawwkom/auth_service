package auth

import (
	"context"
	"github.com/stawwkom/auth_service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Обновление информации пользователя
func (s *service) UpdateUser(ctx context.Context, id int64, user *model.User) error {
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
	}

	return s.authRepo.Update(ctx, id, user)
}
