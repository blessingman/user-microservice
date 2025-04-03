package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/janson/usermicroservice/internal/config"
	"github.com/janson/usermicroservice/internal/handler"
	"github.com/janson/usermicroservice/internal/repository/postgres"
	"github.com/janson/usermicroservice/internal/service"
)

func main() {
	// Загрузка конфигурации из файла config.json
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Настройка логгера для записи в файл
	logFile, err := os.OpenFile(cfg.Logging.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка открытия файла логов: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	logger.Printf("Запуск микросервиса пользователей...")

	// Запуск миграций базы данных для создания таблиц
	runMigrations(cfg, logger)

	// Подключение к базе данных PostgreSQL
	dbpool, err := pgxpool.Connect(context.Background(), cfg.Database.ConnectionString())
	if err != nil {
		logger.Fatalf("Невозможно подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	// Инициализация репозитория, сервиса и обработчика для работы с пользователями
	userRepo := postgres.NewUserRepository(dbpool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, logger)

	// Настройка маршрутизатора и регистрация маршрутов API
	router := mux.NewRouter()
	userHandler.RegisterRoutes(router)

	// Добавление middleware для логирования всех запросов
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
		})
	})

	// Запуск HTTP сервера
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		logger.Printf("Сервер запущен на порту %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Обработка сигналов для корректного завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Сервер принудительно закрыт: %v", err)
	}

	logger.Println("Сервер корректно завершил работу")
}

// runMigrations применяет миграции для создания необходимых таблиц в базе данных
func runMigrations(cfg *config.Config, logger *log.Logger) {
	// Путь к миграциям и строка подключения к базе данных
	migrationsPath := "file://migrations"
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Создание экземпляра для управления миграциями
	m, err := migrate.New(migrationsPath, connString)
	if err != nil {
		logger.Printf("Ошибка инициализации миграций: %v", err)
		return
	}

	// Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Printf("Ошибка применения миграций: %v", err)
		return
	}

	logger.Println("Миграции базы данных успешно применены")
}
