package solver

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/kliuchnikovv/word-of-yoda/domain"
	"github.com/kliuchnikovv/word-of-yoda/internal/utils"
)

func TestSolve_Success(t *testing.T) {
	// Create a test logger that discards output
	logger := slog.Default()

	// Set up a simple challenge with low difficulty for quick testing
	challenge := domain.Challenge{
		ID:         "test-challenge-1",
		Data:       "test-data",
		Difficulty: 8, // Low difficulty for fast tests
	}

	// Solve the challenge
	solution, err := Solve(context.Background(), logger, challenge)

	// Check for errors
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify the solution is not nil
	if solution == nil {
		t.Fatal("Expected solution, got nil")
	}

	// Verify the solution ID matches the challenge ID
	if solution.ID != challenge.ID {
		t.Errorf("Expected solution ID %s, got %s", challenge.ID, solution.ID)
	}

	// Verify the solution is valid
	dataToHash := challenge.Data + strconv.FormatUint(solution.Nonce, 10)
	hash := sha256.Sum256([]byte(dataToHash))
	if !utils.HasLeadingZeroBits(hash[:], challenge.Difficulty) {
		t.Errorf("Solution verification failed: hash does not have %d leading zero bits", challenge.Difficulty)
	}
}

func TestSolve_ContextCancellation(t *testing.T) {
	// Create a test logger that discards output
	logger := slog.Default()

	// Set up a difficult challenge that would take time to solve
	challenge := domain.Challenge{
		ID:         "test-challenge-2",
		Data:       "test-data",
		Difficulty: 24, // High difficulty to ensure we don't solve before cancellation
	}

	// Create a context with cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Try to solve the challenge
	solution, err := Solve(ctx, logger, challenge)

	// Check that we got the expected error
	if err == nil {
		t.Fatal("Expected error due to context cancellation, got nil")
	}

	// Verify that the solution is nil
	if solution != nil {
		t.Errorf("Expected nil solution, got %+v", solution)
	}
}

func TestSolve_DifferentDifficulties(t *testing.T) {
	// Create a test logger that discards output
	logger := slog.Default()

	testCases := []struct {
		name       string
		difficulty int
		timeout    time.Duration
	}{
		{"Very Easy", 4, 1 * time.Second},
		{"Easy", 8, 2 * time.Second},
		{"Medium", 12, 5 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			challenge := domain.Challenge{
				ID:         "difficulty-test",
				Data:       "test-data",
				Difficulty: tc.difficulty,
			}

			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			solution, err := Solve(ctx, logger, challenge)

			if err != nil {
				t.Fatalf("Failed to solve challenge with difficulty %d: %v", tc.difficulty, err)
			}

			if solution == nil {
				t.Fatalf("Expected solution for difficulty %d, got nil", tc.difficulty)
			}

			// Verify the solution
			dataToHash := challenge.Data + strconv.FormatUint(solution.Nonce, 10)
			hash := sha256.Sum256([]byte(dataToHash))
			if !utils.HasLeadingZeroBits(hash[:], challenge.Difficulty) {
				t.Errorf("Solution verification failed for difficulty %d", tc.difficulty)
			}
		})
	}
}

// Test the edge case of difficulty 0
func TestSolve_ZeroDifficulty(t *testing.T) {
	logger := slog.Default()

	challenge := domain.Challenge{
		ID:         "zero-difficulty",
		Data:       "test-data",
		Difficulty: 0,
	}

	solution, err := Solve(context.Background(), logger, challenge)

	if err != nil {
		t.Fatalf("Expected no error for zero difficulty, got: %v", err)
	}

	if solution == nil {
		t.Fatal("Expected solution for zero difficulty, got nil")
	}

	// For zero difficulty, any hash should pass, so nonce should be 0
	if solution.Nonce != 0 {
		t.Errorf("Expected nonce 0 for zero difficulty, got %d", solution.Nonce)
	}
}
