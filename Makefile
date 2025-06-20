.PHONY: test test-unit test-integration test-storage test-ci test-all

# Unit tests (fast, without external dependencies)
test-unit:
	go test -short ./...

# Integration tests with build tags
test-integration:
	go test -tags=integration ./...

# Storage tests with build tags
test-storage:
	docker-compose -f docker-compose.test.yml up -d redis-test
	sleep 10  # Ждем запуска Redis
	CONTAINER_TESTS=true REDIS_CONTAINER_ADDR=localhost:6380 go test -tags=storage ./internal/server/redis/...
	docker-compose -f docker-compose.test.yml down

# CI тесты (с Redis в фоне)
test-ci:
	docker-compose -f docker-compose.test.yml up -d
	sleep 20
	REDIS_TESTS=true CONTAINER_TESTS=true go test -v ./...
	docker-compose -f docker-compose.test.yml down

# Все тесты
test-all: test-integration test-storage

# Очистка
clean-test:
	docker-compose -f docker-compose.test.yml down
	docker volume prune -f

# Запуск тестов с coverage
test-coverage:
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Запуск тестов с race detector
test-race:
	go test -race ./...

# Запуск benchmark тестов
benchmark:
	go test -bench=. ./...
