package challenge

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/kliuchnikovv/word-of-yoda/domain"
	"github.com/kliuchnikovv/word-of-yoda/internal/server/redis"
	"github.com/kliuchnikovv/word-of-yoda/internal/utils"
)

type Challenger struct {
	logger *slog.Logger
	store  redis.Store
}

func NewChallenger(
	logger *slog.Logger,
	store redis.Store) *Challenger {
	return &Challenger{
		logger: logger,
		store:  store,
	}
}

// GenerateChallenge creates a new Challenge with a unique ID, random data, and difficulty
func (c *Challenger) GenerateChallenge(
	ctx context.Context,
	difficulty int,
	ttl time.Duration,
) (*domain.Challenge, error) {
	logger := c.logger.With(
		slog.String("method", "GenerateChallenge"),
		slog.String("difficulty", strconv.Itoa(difficulty)),
		slog.String("ttl", ttl.String()),
	)

	logger.Debug("generating challenge")

	// Generate a unique ID
	id, err := generateRandomHex(8) // 16 hex symbols
	if err != nil {
		logger.Error("failed to generate ID", slog.String("error", err.Error()))
		return nil, err
	}
	// Generate a random part for Data
	randomPart, err := generateRandomHex(16) // 32 hex symbols
	if err != nil {
		logger.Error("failed to generate random part", slog.String("error", err.Error()))
		return nil, err
	}

	// Generate a timestamp to include in the data
	ts := time.Now().UTC().UnixNano()
	data := fmt.Sprintf("%s:%d", randomPart, ts)

	ch := domain.Challenge{
		ID:         id,
		Data:       data,
		Difficulty: difficulty,
		ExpiresAt:  time.Now().UTC().Add(ttl),
	}
	// Save the challenge in the store

	if err := c.store.SaveChallenge(ctx, &ch); err != nil {
		logger.Error("failed to save challenge", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Debug("challenge generated", ch.Log()...)

	return &ch, nil
}

// VerifySolution checks the solution: existence of the puzzle, not expired, and correct nonce
func (c *Challenger) VerifySolution(
	ctx context.Context,
	id string,
	nonce uint64,
) error {
	logger := c.logger.With(
		slog.String("method", "VerifySolution"),
		slog.String("id", id),
		"nonce", nonce,
	)

	ch, err := c.store.GetChallenge(ctx, id)
	if err != nil {
		logger.Error("failed to get challenge", slog.String("error", err.Error()))
		return fmt.Errorf("failed to get challenge: %w", err)
	}

	if time.Now().UTC().After(ch.ExpiresAt.UTC()) {
		logger.Warn("challenge expired", ch.Log()...)
		return fmt.Errorf("challenge expired")
	}

	// Check the solution
	h := sha256.Sum256([]byte(ch.Data + strconv.FormatUint(nonce, 10)))
	if !utils.HasLeadingZeroBits(h[:], ch.Difficulty) {
		logger.Warn("invalid solution", slog.String("hash", hex.EncodeToString(h[:])))
		return fmt.Errorf("invalid solution")
	}

	return nil
}
