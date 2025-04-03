# Многоэтапная сборка для Go приложения

# Этап 1: Сборка приложения
FROM golang:1.20 AS builder

WORKDIR /app

# Копирование всех файлов проекта
COPY . .

# Сборка приложения с флагом CGO_ENABLED=0 напрямую в команде
RUN CGO_ENABLED=0 go build -o userservice ./cmd/api

# Этап 2: Создание минимального образа для запуска
FROM alpine:3.16

# Установка необходимых сертификатов и создание директории для логов
RUN apk --no-cache add ca-certificates && \
    mkdir -p /var/log/userservice

WORKDIR /app

# Копирование только необходимых файлов из этапа сборки
COPY --from=builder /app/userservice .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config.json .

# Открытие порта 8080 для HTTP сервера
EXPOSE 8080

# Команда для запуска приложения
CMD ["./userservice"]
