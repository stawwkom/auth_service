package auth

import (
	"github.com/stawwkom/auth_service/internal/repository"
	serv "github.com/stawwkom/auth_service/internal/service"
)

// Структура, реализующая AuthService
type service struct {
	authRepo repository.AuthRepository
}

// Конструктор Auth-сервиса
func NewAuthService(authRepo repository.AuthRepository) serv.AuthService {
	return &service{authRepo: authRepo}
}
