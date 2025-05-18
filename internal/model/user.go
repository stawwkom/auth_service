package model

import (
	"time"
)

// User для внутреннего использования
type User struct {
	ID        int64 // int64 для совместимости с protobuf
	Login     string
	Email     string
	Password  string
	Role      int // Role (0 = USER, 1 = ADMIN)
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// UserInfo для публичного использования
type UserInfo struct {
	Login string
	Email string
}

type UpdateUserInput struct {
	ID    int64
	Name  *string
	Email *string
}
