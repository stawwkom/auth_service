// Определяем, что этот файл принадлежит пакету "user"
package user

// Импортируем необходимые библиотеки и пакеты
import (
	"context"                            // Для передачи контекста, например, отмены или тайм-аутов
	sq "github.com/Masterminds/squirrel" // Удобный SQL-конструктор запросов
	"github.com/jackc/pgx/v5/pgxpool"    // Пул подключений к PostgreSQL (более эффективен, чем стандартный драйвер)

	// Внутренние пакеты нашего приложения
	"github.com/stawwkom/auth_service/internal/model"                           // Доменная модель User
	"github.com/stawwkom/auth_service/internal/repository"                      // Интерфейс, который должен реализовать этот репозиторий
	"github.com/stawwkom/auth_service/internal/repository/auth/converter"       // Конвертация между доменной и репо моделью
	modelRepo "github.com/stawwkom/auth_service/internal/repository/auth/model" // Модель пользователя для базы (с sql.NullTime и т.п.)
)

// Константы с названиями таблицы и колонок в БД, чтобы не писать их вручную много раз
const (
	tableName       = "users"      // Имя таблицы в базе данных
	idColumn        = "id"         // Уникальный идентификатор пользователя
	loginColumn     = "login"      // Логин пользователя
	emailColumn     = "email"      // Электронная почта
	passwordColumn  = "password"   // Хэшированный пароль
	createdAtColumn = "created_at" // Время создания записи
	updatedAtColumn = "updated_at" // Время последнего обновления
)

// Структура репозитория, содержит пул подключений к БД
type repo struct {
	db *pgxpool.Pool // Подключение к PostgreSQL через пул
}

// Фабричная функция для создания нового репозитория
func NewRepository(db *pgxpool.Pool) repository.AuthRepository {
	// Возвращаем реализацию интерфейса AuthRepository
	return &repo{db: db}
}

// Метод для создания нового пользователя в базе
func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	// Строим SQL-запрос на вставку пользователя
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar). // Используем $1, $2 и т.д. (PostgreSQL стиль)
		Columns(loginColumn, emailColumn, passwordColumn).
		Values(user.Login, user.Email, user.Password).
		Suffix("RETURNING id") // Возвращаем ID созданного пользователя

	// Получаем SQL-запрос и аргументы
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err // Ошибка при построении SQL
	}

	var id int64
	// Выполняем запрос и сканируем ID
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err // Ошибка при выполнении запроса
	}

	return id, nil // Возвращаем ID созданного пользователя
}

// Получение краткой информации о пользователе по ID
func (r *repo) Get(ctx context.Context, id int64) (*model.UserInfo, error) {
	// Строим SQL-запрос на выборку по ID
	builder := sq.Select(idColumn, loginColumn, emailColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}). // WHERE id = $1
		Limit(1)                    // Ограничиваем до 1 записи

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelRepo.User // Временная структура для сканирования
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err // Ошибка при получении данных
	}

	// Конвертируем репо-модель в доменную (без пароля)
	return converter.ToAuthInfoFromRepo(&user), nil
}

// GetByLogin возвращает полную информацию по логину (используется при логине)
func (r *repo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	builder := sq.Select(idColumn, loginColumn, emailColumn, passwordColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{loginColumn: login}). // WHERE login = $1
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelRepo.User
	// Здесь сканируем все поля, включая пароль
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Конвертируем в доменную модель
	return converter.ToAuthFromRepo(&user), nil
}

// Обновление данных пользователя по ID
func (r *repo) Update(ctx context.Context, id int64, user *model.User) error {
	// Строим SQL-запрос на обновление
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(loginColumn, user.Login).
		Set(emailColumn, user.Email).
		Set(passwordColumn, user.Password).
		Set(updatedAtColumn, sq.Expr("NOW()")). // Устанавливаем текущее время обновления
		Where(sq.Eq{idColumn: id})              // WHERE id = $1

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	// Выполняем запрос
	_, err = r.db.Exec(ctx, query, args...)
	return err
}

// Удаление пользователя по ID
func (r *repo) Delete(ctx context.Context, id int64) error {
	// Строим SQL-запрос на удаление
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}) // WHERE id = $1

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	// Выполняем удаление
	_, err = r.db.Exec(ctx, query, args...)
	return err
}
