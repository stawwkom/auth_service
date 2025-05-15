# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -o auth_service ./cmd/auth_server

# Этап запуска
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/auth_service .

# Запускаем сервис
CMD ["./auth_service"] 