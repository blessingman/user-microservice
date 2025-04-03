# User Microservice

Микросервис для управления пользователями с REST API, написанный на Go и использующий PostgreSQL в качестве базы данных!

## Содержание

- [Описание проекта](#описание-проекта)
- [Требования](#требования)
- [Структура проекта](#структура-проекта)
- [Запуск проекта](#запуск-проекта)
- [API Endpoints](#api-endpoints)
- [Тестирование API](#тестирование-api)
- [База данных](#база-данных)
- [CI/CD](#cicd)

## Описание проекта

User Microservice - это REST API для управления пользователями, предоставляющий базовые операции CRUD (Create, Read, Update, Delete). 
Проект использует:
- Go 1.20
- PostgreSQL для хранения данных
- Docker и Docker Compose для контейнеризации
- Миграции для управления схемой базы данных

## Требования

Для запуска проекта вам потребуется:
- Docker и Docker Compose
- Go 1.20 (только для локальной разработки)
- curl или другой инструмент для тестирования API

## Структура проекта

Проект следует принципам чистой архитектуры:
- `cmd/api/` - точка входа в приложение
- `internal/` - внутренние пакеты приложения:
  - `config/` - конфигурация приложения
  - `handler/` - HTTP обработчики
  - `model/` - модели данных
  - `repository/` - слой доступа к данным
  - `service/` - бизнес-логика
- `migrations/` - SQL миграции для создания и наполнения БД
- `docker-compose.yml` - конфигурация Docker Compose
- `Dockerfile` - инструкции для сборки Docker образа

## Запуск проекта

### Запуск с Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <репозиторий>
cd <директория-проекта>
```

2. Запустите сервисы с помощью Docker Compose:
```bash
docker-compose up -d
```

3. Проверьте, что контейнеры запущены:
```bash
docker-compose ps
```

### Остановка проекта

```bash
docker-compose down
```

### Просмотр логов

```bash
docker-compose logs -f app
```

## API Endpoints

Сервис предоставляет следующие API endpoints:

| Метод | URL | Описание |
|-------|-----|----------|
| GET | /users | Получить список всех пользователей |
| GET | /users/{id} | Получить пользователя по ID |
| POST | /users | Создать нового пользователя |
| PUT | /users/{id} | Обновить данные пользователя |
| DELETE | /users/{id} | Удалить пользователя |

## Тестирование API

Ниже приведены примеры curl-запросов для тестирования API:

### Получение всех пользователей

```bash
curl -X GET http://localhost:8080/users
```

### Получение пользователя по ID

```bash
curl -X GET http://localhost:8080/users/1
```

### Создание нового пользователя

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com"
  }'
```

### Обновление пользователя

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated User",
    "email": "updated@example.com"
  }'
```

### Частичное обновление пользователя (только имя)

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Only Name Updated"
  }'
```

### Удаление пользователя

```bash
curl -X DELETE http://localhost:8080/users/1
```

## База данных

### Структура базы данных

База данных содержит таблицу `users` со следующими полями:
- `id`: SERIAL PRIMARY KEY - уникальный идентификатор пользователя
- `name`: VARCHAR(100) NOT NULL - имя пользователя
- `email`: VARCHAR(100) NOT NULL UNIQUE - электронная почта пользователя
- `created_at`: TIMESTAMP WITH TIME ZONE - дата и время создания

### Начальные данные

При первичном запуске в базу данных добавляются тестовые пользователи:
- John Doe (john@example.com)
- Jane Doe (jane@example.com)
- Test User (test@example.com)

### Доступ к базе данных напрямую

Вы можете подключиться к базе данных PostgreSQL напрямую:

```bash
docker-compose exec postgres psql -U postgres -d userservice
```

Полезные PostgreSQL команды:
- `\dt` - показать все таблицы
- `SELECT * FROM users;` - получить всех пользователей
- `\q` - выйти из psql

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки:
- Автоматическая сборка и тестирование при push в main ветку
- Автоматическая сборка и тестирование при создании Pull Request
- Публикация Docker образа в Docker Hub при успешной сборке в main ветке

## Конфигурация

Конфигурация приложения доступна в файле `config.json`. Вы можете изменить настройки для:
- HTTP-сервера (порт)
- Базы данных (хост, порт, имя пользователя, пароль)
- Логирования (путь к файлу логов, уровень логирования)