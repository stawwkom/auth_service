package auth

import (
	"context"
	"github.com/stawwkom/auth_service/internal/model"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Update(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	var login, email string

	if req.Name != nil {
		login = req.Name.Value
	}
	if req.Email != nil {
		email = req.Email.Value
	}

	err := s.authService.UpdateUser(ctx, req.Id, &model.User{
		Login: login,
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
