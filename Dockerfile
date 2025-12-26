# Используем многостадийную сборку
FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для protobuf
RUN apk add --no-cache protoc protobuf-dev

# Устанавливаем Go плагины для protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .
