package model

import (
	"time"
)

// User представляет собой сущность пользователя в системе
// Эта структура используется для передачи данных о пользователе между слоями приложения
type User struct {
	ID        int64     `json:"id"`         // Уникальный идентификатор пользователя
	Name      string    `json:"name"`       // Имя пользователя
	Email     string    `json:"email"`      // Электронная почта (уникальна для каждого пользователя)
	CreatedAt time.Time `json:"created_at"` // Дата и время создания пользователя
}

// UserCreate используется для создания нового пользователя
// Содержит только поля, необходимые для создания пользователя
type UserCreate struct {
	Name  string `json:"name"`  // Имя нового пользователя
	Email string `json:"email"` // Электронная почта нового пользователя
}

// UserUpdate используется для обновления существующего пользователя
// Поля помечены как omitempty, чтобы можно было обновлять только часть полей
type UserUpdate struct {
	Name  string `json:"name,omitempty"`  // Новое имя пользователя (опционально)
	Email string `json:"email,omitempty"` // Новая электронная почта (опционально)
}
