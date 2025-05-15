package repository

import (
	"context"
	"time"
)

// User представляет модель пользователя в базе данных
type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Repository определяет интерфейс для работы с хранилищем данных
type Repository interface {
	// UserRepository определяет методы для работы с пользователями
	UserRepository
}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}
