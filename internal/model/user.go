package model

import (
	"time"
)

// User для внутреннего использования
type User struct {
	ID        int
	Login     string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// UserInfo для публичного использования
type UserInfo struct {
	Login string
	Email string
}
