FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

RUN pwd
RUN ls -la
RUN ls -la config

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/client ./cmd/client

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/client /app/client
COPY --from=builder /app/config/ /app/config/

ENTRYPOINT ["/app/client"]
