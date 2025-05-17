package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int
	Login        string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdateAt     sql.NullTime
}
