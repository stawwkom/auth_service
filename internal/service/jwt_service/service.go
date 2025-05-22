package jwt_service

import "context"

type AuthService interface {
	Login(username string) (string, error)
	Refresh(refreshToken string) (string, error)
	GenerateAccess(refreshToken string) (string, error)
}

type AccessService interface {
	CheckAccess(ctx context.Context, token string, endpoint string) error
}
