package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/janson/usermicroservice/internal/model"
)

// UserRepository обрабатывает операции с базой данных, связанные с пользователями
// Этот тип реализует доступ к данным пользователей в PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool // Пул соединений с базой данных PostgreSQL
}

// NewUserRepository создает новый репозиторий пользователей
// db - пул соединений с базой данных
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create добавляет нового пользователя в базу данных
// ctx - контекст для операции с базой данных
// user - данные для создания пользователя
func (r *UserRepository) Create(ctx context.Context, user model.UserCreate) (*model.User, error) {
	// SQL запрос для вставки нового пользователя и получения его данных
	query := `
		INSERT INTO users (name, email, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, email, created_at
	`

	createdAt := time.Now()
	var createdUser model.User

	// Выполнение запроса и сканирование результатов в структуру User
	err := r.db.QueryRow(ctx, query, user.Name, user.Email, createdAt).
		Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

// GetByID получает пользователя по его идентификатору
// ctx - контекст для операции с базой данных
// id - идентификатор пользователя
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	// SQL запрос для получения пользователя по ID
	query := `
		SELECT id, name, email, created_at 
		FROM users 
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пользователь не найден, возвращаем nil без ошибки
		}
		return nil, err
	}

	return &user, nil
}

// GetAll получает всех пользователей из базы данных
// ctx - контекст для операции с базой данных
func (r *UserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	// SQL запрос для получения всех пользователей, отсортированных по ID
	query := `
		SELECT id, name, email, created_at 
		FROM users
		ORDER BY id
	`

	// Выполнение запроса
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Сканирование результатов в слайс пользователей
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Проверка наличия ошибок при итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update обновляет информацию о пользователе
// ctx - контекст для операции с базой данных
// id - идентификатор пользователя для обновления
// user - данные для обновления
func (r *UserRepository) Update(ctx context.Context, id int64, user model.UserUpdate) (*model.User, error) {
	// Сначала получаем текущего пользователя, чтобы убедиться, что он существует
	currentUser, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if currentUser == nil {
		return nil, nil // Пользователь не найден
	}

	// Обновляем непустые поля
	if user.Name != "" {
		currentUser.Name = user.Name
	}
	if user.Email != "" {
		currentUser.Email = user.Email
	}

	// SQL запрос для обновления пользователя
	query := `
		UPDATE users 
		SET name = $1, email = $2
		WHERE id = $3
		RETURNING id, name, email, created_at
	`

	var updatedUser model.User
	err = r.db.QueryRow(ctx, query, currentUser.Name, currentUser.Email, id).
		Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// Delete удаляет пользователя по ID
// ctx - контекст для операции с базой данных
// id - идентификатор пользователя для удаления
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = $1"

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return nil // Пользователь не найден, не считается ошибкой
	}

	return nil
}
