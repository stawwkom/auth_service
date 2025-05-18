package auth

import (
	"context"
)

// Удаление пользователя
func (s *service) DeleteUser(ctx context.Context, id int64) error {
	return s.authRepo.Delete(ctx, id)
}
