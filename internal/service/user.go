package service

import (
	"context"
	"errors"

	"github.com/janson/usermicroservice/internal/model"
	"github.com/janson/usermicroservice/internal/repository/postgres"
)

// UserService обрабатывает бизнес-логику, связанную с пользователями
// Этот слой служит промежуточным звеном между обработчиками HTTP и репозиторием данных
type UserService struct {
	repo *postgres.UserRepository // Репозиторий для доступа к данным пользователей
}

// NewUserService создает новый сервис пользователей
// repo - репозиторий пользователей для работы с данными
func NewUserService(repo *postgres.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Определение стандартных ошибок для сервиса пользователей
var (
	ErrUserNotFound = errors.New("user not found")     // Пользователь не найден
	ErrInvalidInput = errors.New("invalid input data") // Некорректные входные данные
)

// Create создает нового пользователя
// ctx - контекст операции
// user - данные для создания пользователя
func (s *UserService) Create(ctx context.Context, user model.UserCreate) (*model.User, error) {
	// Валидация входных данных
	if user.Name == "" || user.Email == "" {
		return nil, ErrInvalidInput
	}

	// Делегирование операции создания репозиторию
	return s.repo.Create(ctx, user)
}

// GetByID получает пользователя по ID
// ctx - контекст операции
// id - идентификатор пользователя
func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Если пользователь не найден, возвращаем соответствующую ошибку
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetAll получает всех пользователей
// ctx - контекст операции
func (s *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	// Простая передача вызова к репозиторию
	return s.repo.GetAll(ctx)
}

// Update обновляет информацию о пользователе
// ctx - контекст операции
// id - идентификатор пользователя
// user - данные для обновления
func (s *UserService) Update(ctx context.Context, id int64, user model.UserUpdate) (*model.User, error) {
	// Должно быть обновлено хотя бы одно поле
	if user.Name == "" && user.Email == "" {
		return nil, ErrInvalidInput
	}

	updatedUser, err := s.repo.Update(ctx, id, user)
	if err != nil {
		return nil, err
	}

	// Если пользователь не найден, возвращаем соответствующую ошибку
	if updatedUser == nil {
		return nil, ErrUserNotFound
	}

	return updatedUser, nil
}

// Delete удаляет пользователя по ID
// ctx - контекст операции
// id - идентификатор пользователя
func (s *UserService) Delete(ctx context.Context, id int64) error {
	// Сначала проверяем, существует ли пользователь
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Делегируем операцию удаления репозиторию
	return s.repo.Delete(ctx, id)
}
