package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdateAt     sql.NullTime
}
