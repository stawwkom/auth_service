package auth

import (
	"context"
	converter "github.com/stawwkom/auth_service/internal/converter"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
)

// Create обрабатывает gRPC запрос на создание пользователя получает запрос и возвращает ID.
func (s *Server) Create(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	// Конвертация запроса в бизнес-модель
	user := converter.ToModelUser(req)

	// Вызов сервисного слоя
	id, err := s.authService.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	// Возврат ответа (ID)
	return &desc.CreateUserResponse{Id: id}, nil
}
