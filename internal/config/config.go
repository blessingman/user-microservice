package config

import (
	"encoding/json"
	"os"
)

// Config содержит все настройки для сервиса
// Используется для загрузки конфигурации из JSON файла
type Config struct {
	Server   ServerConfig   `json:"server"`   // Настройки HTTP сервера
	Database DatabaseConfig `json:"database"` // Настройки базы данных
	Logging  LoggingConfig  `json:"logging"`  // Настройки логирования
}

// ServerConfig содержит настройки HTTP сервера
type ServerConfig struct {
	Port string `json:"port"` // Порт, на котором будет работать сервер
}

// DatabaseConfig содержит настройки подключения к базе данных
type DatabaseConfig struct {
	Host     string `json:"host"`     // Хост базы данных
	Port     string `json:"port"`     // Порт базы данных
	User     string `json:"user"`     // Имя пользователя для подключения к БД
	Password string `json:"password"` // Пароль для подключения к БД
	DBName   string `json:"dbname"`   // Имя базы данных
	SSLMode  string `json:"sslmode"`  // Режим SSL (обычно disable для локальной разработки)
}

// LoggingConfig содержит настройки логирования
type LoggingConfig struct {
	FilePath string `json:"file_path"` // Путь к файлу логов
	Level    string `json:"level"`     // Уровень логирования (info, debug, error и т.д.)
}

// ConnectionString возвращает строку подключения к PostgreSQL
// Используется для подключения к базе данных через pgxpool
func (c DatabaseConfig) ConnectionString() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

// Load читает конфигурацию из файла и возвращает структуру Config
// path - путь к JSON файлу конфигурации
func Load(path string) (*Config, error) {
	// Открытие файла конфигурации
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Декодирование JSON в структуру Config
	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
