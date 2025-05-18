package auth

import (
	"context"
	auth "github.com/stawwkom/auth_service/internal/converter"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
)

func (s *Server) Get(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	user, err := s.authService.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return auth.ToProtoUserInfo(user), nil
}
