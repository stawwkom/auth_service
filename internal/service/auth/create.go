package auth

import (
	"context"
	"github.com/stawwkom/auth_service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Реализация метода Register
func (s *service) Register(ctx context.Context, user *model.User) (int64, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashed)
	return s.authRepo.Create(ctx, user)
}
