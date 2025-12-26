FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем всё из корня проекта
COPY .. .

# Устанавливаем зависимости
RUN go mod download

# Собираем save_service
WORKDIR /app/save_service
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/save_service/main .
COPY save_service/config ./config

EXPOSE 50051

CMD ["./main"]
