package challenge

import (
	"crypto/sha256"
	"strconv"
	"time"

	"github.com/kliuchnikovv/word-of-yoda/domain"
	"github.com/kliuchnikovv/word-of-yoda/internal/utils"
)

// calculateValidNonce находит nonce, который решает challenge
func calculateValidNonce(data string, difficulty int) uint64 {
	var nonce uint64 = 0

	for nonce < 1000000 { // Защита от бесконечного цикла
		testData := data + strconv.FormatUint(nonce, 10)
		hash := sha256.Sum256([]byte(testData))

		if utils.HasLeadingZeroBits(hash[:], difficulty) {
			return nonce
		}

		nonce++
	}

	return nonce
}

// createTestChallenge создает тестовый challenge
func createTestChallenge(id string, difficulty int, ttl time.Duration) *domain.Challenge {
	return &domain.Challenge{
		ID:         id,
		Data:       "testdata:1640995200000000000",
		Difficulty: difficulty,
		ExpiresAt:  time.Now().UTC().Add(ttl),
	}
}

// createExpiredChallenge создает истекший challenge
func createExpiredChallenge(id string, difficulty int) *domain.Challenge {
	return &domain.Challenge{
		ID:         id,
		Data:       "expireddata:1640995200000000000",
		Difficulty: difficulty,
		ExpiresAt:  time.Now().UTC().Add(-1 * time.Hour), // Истек час назад
	}
}
