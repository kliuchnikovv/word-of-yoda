package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/kliuchnikovv/word-of-yoda/domain"

	"github.com/redis/go-redis/v9"
)

// SaveChallenge saves a challenge with TTL
func (r *RedisStore) SaveChallenge(ctx context.Context, challenge *domain.Challenge) error {
	if challenge == nil {
		return fmt.Errorf("challenge cannot be nil")
	}

	logger := r.logger.With(
		slog.String("method", "SaveChallenge"),
		slog.String("challenge_id", challenge.ID),
		"expires_at", challenge.ExpiresAt,
	)

	logger.Debug("Saving challenge")

	// Serialize to JSON
	data, err := json.Marshal(challenge)
	if err != nil {
		logger.Error("Failed to marshal challenge", slog.String("error", err.Error()))
		return fmt.Errorf("failed to marshal challenge: %w", err)
	}

	key := r.config.ChallengePrefix + challenge.ID

	// Calculate TTL
	var ttl time.Duration
	if !challenge.ExpiresAt.IsZero() {
		ttl = time.Until(challenge.ExpiresAt)
		if ttl <= 0 {
			logger.Warn("Challenge already expired")
			return fmt.Errorf("challenge already expired")
		}
	}

	// Save to Redis with TTL
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		logger.Error("Failed to save challenge to Redis", slog.String("error", err.Error()))
		return fmt.Errorf("failed to save challenge: %w", err)
	}

	logger.Info("Challenge saved successfully", "ttl_seconds", ttl.Seconds())
	return nil
}

// GetChallenge retrieves a challenge by ID
func (r *RedisStore) GetChallenge(ctx context.Context, id string) (*domain.Challenge, error) {
	logger := r.logger.With(
		slog.String("method", "GetChallenge"),
		slog.String("challenge_id", id),
	)

	logger.Debug("Getting challenge")

	key := r.config.ChallengePrefix + id

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Debug("Challenge not found")
			return nil, ErrNotFound{Key: key}
		}
		logger.Error("Failed to get challenge from Redis", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}

	var challenge domain.Challenge
	if err := json.Unmarshal([]byte(data), &challenge); err != nil {
		logger.Error("Failed to unmarshal challenge", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to unmarshal challenge: %w", err)
	}

	// Check if challenge has expired
	if !challenge.ExpiresAt.IsZero() && time.Now().UTC().After(challenge.ExpiresAt.UTC()) {
		logger.Debug("Challenge expired, deleting")
		r.client.Del(ctx, key) // delete expired challenge
		return nil, ErrNotFound{Key: key}
	}

	logger.Debug("Challenge retrieved successfully", challenge.Log()...)
	return &challenge, nil
}

// DeleteChallenge deletes a challenge
func (r *RedisStore) DeleteChallenge(
	ctx context.Context,
	id string,
) error {
	logger := r.logger.With(
		slog.String("method", "DeleteChallenge"),
		slog.String("challenge_id", id),
	)

	logger.Debug("Deleting challenge")

	key := r.config.ChallengePrefix + id

	result := r.client.Del(ctx, key)
	if err := result.Err(); err != nil {
		logger.Error("Failed to delete challenge", slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete challenge: %w", err)
	}

	if result.Val() == 0 {
		logger.Debug("Challenge not found for deletion")
		return ErrNotFound{Key: key}
	}

	logger.Info("Challenge deleted successfully")
	return nil
}

// ListChallenges returns a list of challenges
func (r *RedisStore) ListChallenges(
	ctx context.Context,
	pattern string,
	limit int,
) ([]*domain.Challenge, error) {
	logger := r.logger.With(
		slog.String("method", "ListChallenges"),
		slog.String("pattern", pattern),
		"limit", limit,
	)

	logger.Debug("Listing challenges")

	if pattern == "" {
		pattern = "*"
	}

	searchPattern := r.config.ChallengePrefix + pattern

	// Retrieve keys
	keys, err := r.client.Keys(ctx, searchPattern).Result()
	if err != nil {
		logger.Error("Failed to get challenge keys", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get challenge keys: %w", err)
	}

	if limit > 0 && len(keys) > limit {
		keys = keys[:limit]
	}

	if len(keys) == 0 {
		logger.Debug("No challenges found")
		return []*domain.Challenge{}, nil
	}

	// Retrieve data
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))

	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		logger.Error("Failed to execute pipeline", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	var challenges []*domain.Challenge
	now := time.Now()

	for i, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			logger.Warn("Failed to get challenge data", "key", keys[i], "error", err)
			continue
		}

		var challenge domain.Challenge
		if err := json.Unmarshal([]byte(data), &challenge); err != nil {
			logger.Warn("Failed to unmarshal challenge", "key", keys[i], "error", err)
			continue
		}

		// Skip expired challenges
		if !challenge.ExpiresAt.IsZero() && now.After(challenge.ExpiresAt) {
			logger.Debug("Skipping expired challenge", "challenge_id", challenge.ID)
			// delete expired challenge asynchronously
			go r.client.Del(context.Background(), keys[i])
			continue
		}

		challenges = append(challenges, &challenge)
	}

	logger.Debug("Challenges listed successfully", "count", len(challenges))
	return challenges, nil
}
