FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .
COPY config/ /app/config/

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/config/ /app/config/

ENTRYPOINT ["/app/server"]
