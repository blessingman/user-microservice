-- Миграция для создания таблицы пользователей
-- Выполняется при обновлении базы данных

-- Создание таблицы пользователей (если не существует)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                             -- Уникальный идентификатор пользователя, автоинкрементный
    name VARCHAR(100) NOT NULL,                        -- Имя пользователя, обязательное поле
    email VARCHAR(100) NOT NULL UNIQUE,                -- Email пользователя, обязательное и уникальное поле
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата и время создания с часовым поясом, по умолчанию текущее время
);

-- Создание индекса по полю email для быстрого поиска пользователей
CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);