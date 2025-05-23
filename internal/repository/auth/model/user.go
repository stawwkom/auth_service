package model

import (
	"database/sql"
	"time"
)

// User для внутреннего использования определяет структуру БД
type User struct {
	ID        int64
	Login     string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
