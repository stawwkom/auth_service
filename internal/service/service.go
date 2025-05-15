package service

import (
	"context"
)

// User представляет модель пользователя в бизнес-логике
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// Service определяет интерфейс для бизнес-логики
type Service interface {
	// UserService определяет методы для работы с пользователями
	UserService
}

// UserService определяет методы для работы с пользователями
type UserService interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email, password string) (*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}
