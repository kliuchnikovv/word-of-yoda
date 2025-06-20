package solver

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/kliuchnikovv/word-of-yoda/domain"
	"github.com/kliuchnikovv/word-of-yoda/internal/utils"
)

func Solve(
	ctx context.Context,
	logger *slog.Logger,
	challenge domain.Challenge,
) (*domain.Solution, error) {
	logger = logger.With(slog.String("method", "Solve"))

	var nonce uint64 = 0

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("solving was cancelled")
		default:
			dataToHash := challenge.Data + strconv.FormatUint(nonce, 10)

			hash := sha256.Sum256([]byte(dataToHash))

			if utils.HasLeadingZeroBits(hash[:], challenge.Difficulty) {
				logger.Debug("solution found", slog.String("hash", hex.EncodeToString(hash[:])))
				return &domain.Solution{
					ID:    challenge.ID,
					Nonce: nonce,
				}, nil
			}

			nonce++
		}
	}
}
