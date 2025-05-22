package auth_handler

import (
	"context"
	"github.com/stawwkom/auth_service/internal/service/jwt_service"
	descAuth "github.com/stawwkom/auth_service/pkg/auth_login"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService struct {
	descAuth.UnimplementedAuthV1Server
	authService jwt_service.AuthService
}

func NewAuthHandler(auth jwt_service.AuthService) descAuth.AuthV1Server {
	return &authService{
		authService: auth,
	}
}

func (a *authService) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	token, err := a.authService.Login(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login: %v", err)
	}
	return &descAuth.LoginResponse{RefreshToken: token}, nil
}

func (a *authService) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	token, err := a.authService.Refresh(req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to refresh token: %v", err)
	}
	return &descAuth.GetRefreshTokenResponse{RefreshToken: token}, nil
}

func (a *authService) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
	token, err := a.authService.GenerateAccess(req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to generate access token: %v", err)
	}
	return &descAuth.GetAccessTokenResponse{AccessToken: token}, nil
}
