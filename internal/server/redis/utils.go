package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Exists проверяет существование ключа
func (r *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to check key existence",
			slog.String("key", key),
			slog.String("error", err.Error()))
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}

	return result > 0, nil
}

// SetTTL устанавливает TTL для ключа
func (r *RedisStore) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
		r.logger.Error("Failed to set TTL",
			slog.String("key", key),
			slog.Int64("ttl_ns", ttl.Nanoseconds()),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to set TTL: %w", err)
	}

	r.logger.Debug("TTL set successfully",
		"key", key,
		"ttl", ttl)
	return nil
}

// GetTTL получает TTL для ключа
func (r *RedisStore) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to get TTL",
			slog.String("key", key),
			slog.String("error", err.Error()),
		)
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return ttl, nil
}

// Ping проверяет соединение с Redis
func (r *RedisStore) Ping(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		r.logger.Error("Redis ping failed", slog.String("error", err.Error()))
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

// Close закрывает соединение с Redis
func (r *RedisStore) Close() error {
	if err := r.client.Close(); err != nil {
		r.logger.Error("Failed to close Redis connection", slog.String("error", err.Error()))
		return fmt.Errorf("failed to close redis connection: %w", err)
	}

	r.logger.Info("Redis connection closed")
	return nil
}
