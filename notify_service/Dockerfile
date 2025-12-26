FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY .. .

RUN apk add --no-cache protoc protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go mod download

WORKDIR /app/notify_service
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/notify_service/main .
COPY notify_service/config ./config

EXPOSE 50053

CMD ["./main"]
