package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/janson/usermicroservice/internal/model"
)

// Базовый URL для тестового API
var baseURL = "http://localhost:8080"

// Инициализация базового URL из переменных окружения
func init() {
	if envURL := os.Getenv("API_URL"); envURL != "" {
		baseURL = envURL
		fmt.Printf("Используется URL из переменной окружения: %s\n", baseURL)
	}
}

// Структура для хранения созданного пользователя между тестами
var createdUserID int64

// TestMain подготавливает окружение для тестов
func TestMain(m *testing.M) {
	// Ждем, пока API станет доступным
	waitForAPI()

	// Запускаем тесты
	exitCode := m.Run()

	// Выходим с кодом, возвращенным из тестов
	os.Exit(exitCode)
}

// waitForAPI ожидает, пока API не станет доступным
func waitForAPI() {
	maxRetries := 15                 // Увеличиваем количество попыток
	retryInterval := 4 * time.Second // Увеличиваем интервал между попытками

	fmt.Println("Ожидание доступности API...")

	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(fmt.Sprintf("%s/users", baseURL))
		if err == nil {
			resp.Body.Close()
			fmt.Println("API доступен!")
			return
		}

		// Пробуем альтернативный URL с IP-адресом
		altURL := "http://127.0.0.1:8080"
		resp, err = http.Get(fmt.Sprintf("%s/users", altURL))
		if err == nil {
			resp.Body.Close()
			fmt.Println("API доступен через 127.0.0.1!")
			baseURL = altURL // Переключаемся на рабочий URL
			return
		}

		fmt.Printf("Попытка %d/%d: API не доступен, ожидание %s...\n", i+1, maxRetries, retryInterval)
		time.Sleep(retryInterval)
	}
	fmt.Println("API не стал доступным после всех попыток, тесты могут не пройти")
}

// TestCreateUser проверяет создание пользователя
func TestCreateUser(t *testing.T) {
	userData := map[string]string{
		"name":  "Test User",
		"email": fmt.Sprintf("test%d@example.com", time.Now().Unix()), // Уникальный email
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Ошибка маршалинга JSON: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/users", baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusCreated, resp.StatusCode)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	if user.ID == 0 {
		t.Error("ID созданного пользователя равен 0")
	}

	if user.Name != userData["name"] {
		t.Errorf("Ожидалось имя %s, получено %s", userData["name"], user.Name)
	}

	if user.Email != userData["email"] {
		t.Errorf("Ожидался email %s, получен %s", userData["email"], user.Email)
	}

	// Сохраняем ID созданного пользователя для следующих тестов
	createdUserID = user.ID
}

// TestGetUser проверяет получение пользователя по ID
func TestGetUser(t *testing.T) {
	if createdUserID == 0 {
		t.Skip("Пропуск теста: не найден ID пользователя")
	}

	resp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, createdUserID))
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusOK, resp.StatusCode)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	if user.ID != createdUserID {
		t.Errorf("Ожидался ID %d, получен %d", createdUserID, user.ID)
	}
}

// TestGetAllUsers проверяет получение всех пользователей
func TestGetAllUsers(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/users", baseURL))
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusOK, resp.StatusCode)
	}

	var users []model.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	if len(users) == 0 {
		t.Error("Список пользователей пуст")
	}
}

// TestUpdateUser проверяет обновление пользователя
func TestUpdateUser(t *testing.T) {
	if createdUserID == 0 {
		t.Skip("Пропуск теста: не найден ID пользователя")
	}

	userData := map[string]string{
		"name": "Updated Test User",
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Ошибка маршалинга JSON: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/users/%d", baseURL, createdUserID),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusOK, resp.StatusCode)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	if user.Name != userData["name"] {
		t.Errorf("Ожидалось имя %s, получено %s", userData["name"], user.Name)
	}
}

// TestDeleteUser проверяет удаление пользователя
func TestDeleteUser(t *testing.T) {
	if createdUserID == 0 {
		t.Skip("Пропуск теста: не найден ID пользователя")
	}

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/users/%d", baseURL, createdUserID),
		nil,
	)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusNoContent, resp.StatusCode)
	}
}
