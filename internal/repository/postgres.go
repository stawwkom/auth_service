package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresRepository реализует интерфейс Repository для PostgreSQL
// Добавлен StatementBuilderType для squirrel
// Используется sqlx для выполнения запросов

type PostgresRepository struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

// NewPostgresRepository создает новый экземпляр PostgresRepository
func NewPostgresRepository(connString string) (*PostgresRepository, error) {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

// Close закрывает соединение с базой данных
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

// CreateUser создает нового пользователя
func (r *PostgresRepository) CreateUser(ctx context.Context, user *User) error {
	now := time.Now()
	query, args, err := r.sb.Insert("users").
		Columns("email", "password", "created_at", "updated_at").
		Values(user.Email, user.Password, now, now).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("ошибка построения запроса: %w", err)
	}
	return r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
}

// GetUserByEmail получает пользователя по email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query, args, err := r.sb.Select("id", "email", "password", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ошибка построения запроса: %w", err)
	}
	var user User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}
	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	query, args, err := r.sb.Select("id", "email", "password", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ошибка построения запроса: %w", err)
	}
	var user User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}
	return &user, nil
}

// UpdateUser обновляет данные пользователя
func (r *PostgresRepository) UpdateUser(ctx context.Context, user *User) error {
	query, args, err := r.sb.Update("users").
		Set("email", user.Email).
		Set("password", user.Password).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ошибка построения запроса: %w", err)
	}
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка обновления пользователя: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", user.ID)
	}
	return nil
}

// DeleteUser удаляет пользователя
func (r *PostgresRepository) DeleteUser(ctx context.Context, id int64) error {
	query, args, err := r.sb.Delete("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ошибка построения запроса: %w", err)
	}
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удалённых строк: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", id)
	}
	return nil
}
