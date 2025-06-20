//go:build integration
// +build integration

package challenge

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	mocks "github.com/kliuchnikovv/word-of-yoda/internal/server/redis/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite для интеграционных тестов
type IntegrationTestSuite struct {
	suite.Suite
	challenger *Challenger
	mockStore  *mocks.MockStore
	ctx        context.Context
}

func (suite *IntegrationTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.mockStore = mocks.NewMockStore(ctrl)
	suite.challenger = NewChallenger(slog.Default(), suite.mockStore)
	suite.ctx = context.Background()
}

func (suite *IntegrationTestSuite) TestFullChallengeFlow() {
	t := suite.T()

	difficulty := 4
	ttl := 5 * time.Minute

	// Настраиваем мок для сохранения
	suite.mockStore.EXPECT().SaveChallenge(gomock.Any(), gomock.Any()).Return(nil)

	// 1. Генерируем challenge
	challenge, err := suite.challenger.GenerateChallenge(suite.ctx, difficulty, ttl)
	require.NoError(t, err)
	require.NotNil(t, challenge)

	// 2. Находим валидный nonce
	validNonce := calculateValidNonce(challenge.Data, challenge.Difficulty)

	// 3. Настраиваем мок для получения challenge
	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), challenge.ID).Return(challenge, nil)

	// 4. Верифицируем решение
	err = suite.challenger.VerifySolution(suite.ctx, challenge.ID, validNonce)
	assert.NoError(t, err)
}

func (suite *IntegrationTestSuite) TestGenerateAndVerifyMultipleChallenges() {
	t := suite.T()

	difficulties := []int{1, 4, 6, 8}

	for _, difficulty := range difficulties {
		suite.Run(fmt.Sprintf("difficulty_%d", difficulty), func() {
			challenge := createTestChallenge(fmt.Sprintf("test_%d", difficulty), difficulty, 10*time.Minute)
			validNonce := calculateValidNonce(challenge.Data, challenge.Difficulty)

			suite.mockStore.EXPECT().GetChallenge(gomock.Any(), challenge.ID).Return(challenge, nil)

			err := suite.challenger.VerifySolution(suite.ctx, challenge.ID, validNonce)
			assert.NoError(t, err)
		})
	}
}

func (suite *IntegrationTestSuite) TestConcurrentChallengeGeneration() {
	t := suite.T()

	const numChallenges = 10
	results := make(chan error, numChallenges)

	// Настраиваем мок для всех сохранений
	suite.mockStore.EXPECT().SaveChallenge(gomock.Any(), gomock.Any()).Return(nil).Times(numChallenges)

	// Запускаем горутины
	for i := 0; i < numChallenges; i++ {
		go func() {
			_, err := suite.challenger.GenerateChallenge(suite.ctx, 4, 5*time.Minute)
			results <- err
		}()
	}

	// Собираем результаты
	for i := 0; i < numChallenges; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
