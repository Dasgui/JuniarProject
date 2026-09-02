# Базовый образ для сборки
FROM golang:1.26-alpine AS builder


# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/.

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

# Копируем бинарный файл из builder
COPY --from=builder /app/app .

# Открываем порт
EXPOSE 8081

# Запускаем приложение
CMD ["./app"]