package auth

import (
	"context"
	converter "github.com/stawwkom/auth_service/internal/converter"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Get(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	userInfo, err := s.authService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	if userInfo == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return converter.ToProtoUserInfo(userInfo), nil
}
