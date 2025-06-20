package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kliuchnikovv/word-of-yoda/domain"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=redis.go -destination=mocks/mock_redis.go -package=redis_mocks

// Store defines the interface for working with the store
type Store interface {
	// Challenge operations
	SaveChallenge(ctx context.Context, challenge *domain.Challenge) error
	GetChallenge(ctx context.Context, id string) (*domain.Challenge, error)
	DeleteChallenge(ctx context.Context, id string) error
	ListChallenges(ctx context.Context, pattern string, limit int) ([]*domain.Challenge, error)

	// Utility operations
	Exists(ctx context.Context, key string) (bool, error)
	SetTTL(ctx context.Context, key string, ttl time.Duration) error
	GetTTL(ctx context.Context, key string) (time.Duration, error)

	// Health check
	Ping(ctx context.Context) error

	// Close connection
	Close() error
}

// ErrNotFound is returned when an object is not found
type ErrNotFound struct {
	Key string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("key not found: %s", e.Key)
}

// RedisStore represents a Redis store
type RedisStore struct {
	client *redis.Client
	logger *slog.Logger
	config *RedisConfig
}

// RedisConfig represents the configuration for Redis
type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`

	// Prefixes для разных типов данных
	MessagePrefix   string `json:"message_prefix"`
	ChallengePrefix string `json:"challenge_prefix"`

	// Timeouts
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`

	// Pool settings
	PoolSize     int           `json:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns"`
	MaxIdleTime  time.Duration `json:"max_idle_time"`
}

// DefaultRedisConfig returns the default configuration
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,

		MessagePrefix:   "msg:",
		ChallengePrefix: "chl:",

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		PoolSize:     10,
		MinIdleConns: 5,
		MaxIdleTime:  5 * time.Minute,
	}
}

// NewRedisStore creates a new Redis store
func NewRedisStore(
	logger *slog.Logger,
	config *RedisConfig,
) (*RedisStore, error) {
	if config == nil {
		config = DefaultRedisConfig()
	}

	if logger == nil {
		logger = slog.Default()
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,

		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,

		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	store := &RedisStore{
		client: rdb,
		logger: logger.With(slog.String("component", "redis_store")),
		config: config,
	}

	// Check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if status := store.client.Ping(ctx); status.Err() != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", status.Err())
	}

	store.logger.Info("Redis store initialized",
		"addr", config.Addr,
		"db", config.DB,
	)

	return store, nil
}
