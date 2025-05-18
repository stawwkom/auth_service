package auth

import (
	"context"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Delete(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.authService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
