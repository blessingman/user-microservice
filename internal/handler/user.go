package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/janson/usermicroservice/internal/model"
	"github.com/janson/usermicroservice/internal/service"
)

// UserHandler обрабатывает HTTP запросы, связанные с пользователями
// Этот тип отвечает за преобразование HTTP запросов в вызовы сервиса
type UserHandler struct {
	service *service.UserService // Сервис для выполнения бизнес-логики
	logger  *log.Logger          // Логгер для записи информации о запросах
}

// NewUserHandler создает новый обработчик пользователей
// service - сервис пользователей
// logger - логгер для записи событий
func NewUserHandler(service *service.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes регистрирует все маршруты для работы с пользователями
// r - маршрутизатор, в который будут добавлены маршруты
func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", h.GetAllUsers).Methods(http.MethodGet)        // GET /users - получить всех пользователей
	r.HandleFunc("/users/{id}", h.GetUser).Methods(http.MethodGet)       // GET /users/{id} - получить пользователя по ID
	r.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)        // POST /users - создать нового пользователя
	r.HandleFunc("/users/{id}", h.UpdateUser).Methods(http.MethodPut)    // PUT /users/{id} - обновить пользователя
	r.HandleFunc("/users/{id}", h.DeleteUser).Methods(http.MethodDelete) // DELETE /users/{id} - удалить пользователя
}

// GetAllUsers обрабатывает GET /users
// Возвращает список всех пользователей
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("Обработка запроса: %s %s", r.Method, r.URL.Path)

	users, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Printf("Ошибка получения пользователей: %v", err)
		http.Error(w, "Не удалось получить пользователей", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

// GetUser обрабатывает GET /users/{id}
// Возвращает пользователя с указанным ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("Обработка запроса: %s %s", r.Method, r.URL.Path)

	id, err := parseIDFromRequest(r)
	if err != nil {
		h.logger.Printf("Ошибка разбора ID пользователя: %v", err)
		http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		h.logger.Printf("Ошибка получения пользователя: %v", err)
		http.Error(w, "Не удалось получить пользователя", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// CreateUser обрабатывает POST /users
// Создает нового пользователя
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("Обработка запроса: %s %s", r.Method, r.URL.Path)

	var userCreate model.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&userCreate); err != nil {
		http.Error(w, "Некорректное тело запроса", http.StatusBadRequest)
		return
	}

	user, err := h.service.Create(r.Context(), userCreate)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			http.Error(w, "Некорректные входные данные", http.StatusBadRequest)
			return
		}
		h.logger.Printf("Ошибка создания пользователя: %v", err)
		http.Error(w, "Не удалось создать пользователя", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// UpdateUser обрабатывает PUT /users/{id}
// Обновляет информацию о пользователе
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("Обработка запроса: %s %s", r.Method, r.URL.Path)

	id, err := parseIDFromRequest(r)
	if err != nil {
		h.logger.Printf("Ошибка разбора ID пользователя: %v", err)
		http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	var userUpdate model.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		http.Error(w, "Некорректное тело запроса", http.StatusBadRequest)
		return
	}

	user, err := h.service.Update(r.Context(), id, userUpdate)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			http.Error(w, "Некорректные входные данные", http.StatusBadRequest)
			return
		}
		h.logger.Printf("Ошибка обновления пользователя: %v", err)
		http.Error(w, "Не удалось обновить пользователя", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser обрабатывает DELETE /users/{id}
// Удаляет пользователя с указанным ID
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("Обработка запроса: %s %s", r.Method, r.URL.Path)

	id, err := parseIDFromRequest(r)
	if err != nil {
		h.logger.Printf("Ошибка разбора ID пользователя: %v", err)
		http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		h.logger.Printf("Ошибка удаления пользователя: %v", err)
		http.Error(w, "Не удалось удалить пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Вспомогательные функции

// parseIDFromRequest извлекает ID пользователя из параметров запроса
func parseIDFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return 0, errors.New("не указан параметр id")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// respondWithJSON отправляет ответ в формате JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}
