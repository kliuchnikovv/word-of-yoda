package challenge

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	mocks "github.com/kliuchnikovv/word-of-yoda/internal/server/redis/mocks"
)

func BenchmarkGenerateChallenge(b *testing.B) {
	ctrl := gomock.NewController(b)

	mockStore := mocks.NewMockStore(ctrl)
	challenger := NewChallenger(slog.Default(), mockStore)
	ctx := context.Background()

	// Настраиваем мок для всех вызовов
	mockStore.EXPECT().SaveChallenge(gomock.Any(), gomock.Any()).Return(nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := challenger.GenerateChallenge(ctx, 4, 5*time.Minute)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkVerifySolution(b *testing.B) {
	ctrl := gomock.NewController(b)

	mockStore := mocks.NewMockStore(ctrl)

	challenger := NewChallenger(slog.Default(), mockStore)
	ctx := context.Background()

	// Подготавливаем данные
	challenge := createTestChallenge("benchmark", 4, 10*time.Minute)
	validNonce := calculateValidNonce(challenge.Data, challenge.Difficulty)

	// Настраиваем мок
	mockStore.EXPECT().GetChallenge(gomock.Any(), "benchmark").Return(challenge, nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := challenger.VerifySolution(ctx, "benchmark", validNonce)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
