//go:build storage
// +build storage

package redis

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/kliuchnikovv/word-of-yoda/domain"

	"github.com/google/uuid"
)

func setupTestStore(t *testing.T) *RedisStore {
	config := &RedisConfig{
		Addr:            "localhost:6379",
		DB:              1, // Используем отдельную БД для тестов
		MessagePrefix:   "test:msg:",
		ChallengePrefix: "test:chl:",
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
	}

	store, err := NewRedisStore(config, slog.Default())
	if err != nil {
		t.Fatalf("Failed to create Redis store: %v", err)
	}

	return store
}

func TestRedisStore_Challenge(t *testing.T) {
	store := setupTestStore(t)
	defer store.Close()

	ctx := context.Background()

	// Создаем тестовый challenge
	challenge := &domain.Challenge{
		ID:         uuid.New().String(),
		Data:       "test data",
		Difficulty: 4,
		ExpiresAt:  time.Now().Add(10 * time.Second),
	}

	// Сохраняем
	err := store.SaveChallenge(ctx, challenge)
	if err != nil {
		t.Fatalf("Failed to save challenge: %v", err)
	}

	// Получаем
	retrieved, err := store.GetChallenge(ctx, challenge.ID)
	if err != nil {
		t.Fatalf("Failed to get challenge: %v", err)
	}

	// Проверяем
	if retrieved.ID != challenge.ID || retrieved.Data != challenge.Data {
		t.Errorf("Retrieved challenge doesn't match. Got: %+v, Want: %+v", retrieved, challenge)
	}

	// Проверяем TTL
	key := store.config.ChallengePrefix + challenge.ID
	ttl, err := store.GetTTL(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get TTL: %v", err)
	}

	if ttl <= 0 || ttl > 10*time.Second {
		t.Errorf("Invalid TTL: %v", ttl)
	}

	// Удаляем
	err = store.DeleteChallenge(ctx, challenge.ID)
	if err != nil {
		t.Fatalf("Failed to delete challenge: %v", err)
	}
}
