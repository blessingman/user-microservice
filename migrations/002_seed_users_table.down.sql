-- Миграция для отката заполнения таблицы тестовыми данными
-- Выполняется при откате базы данных

-- Удаление тестовых пользователей, добавленных в миграции seed
DELETE FROM users
WHERE email IN ('john@example.com', 'jane@example.com', 'test@example.com');
