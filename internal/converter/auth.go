package auth

import (
	"github.com/stawwkom/auth_service/internal/model"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoUser(u *model.User) *desc.GetUserResponse {
	resp := &desc.GetUserResponse{
		Id:        u.ID,
		Name:      u.Login,
		Email:     u.Email,
		Role:      desc.Role(u.Role),
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
	if u.UpdatedAt != nil {
		resp.UpdatedAt = timestamppb.New(*u.UpdatedAt)
	}
	return resp
}

func ToModelUser(req *desc.CreateUserRequest) *model.User {
	return &model.User{
		Login:    req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     int(req.Role), // Приведение enum к int, если поле Role есть в модели
	}
}

func ToProtoUserInfo(u *model.UserInfo) *desc.GetUserResponse {
	return &desc.GetUserResponse{
		Id:    u.ID,
		Name:  u.Login,
		Email: u.Email,
		// Остальные поля (id, role, timestamps) оставить пустыми или не заполнять
	}
}
