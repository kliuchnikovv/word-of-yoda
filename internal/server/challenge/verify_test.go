package challenge

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	mocks "github.com/kliuchnikovv/word-of-yoda/internal/server/redis/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// VerifySolutionTestSuite для тестов верификации
type VerifySolutionTestSuite struct {
	suite.Suite
	challenger *Challenger
	mockStore  *mocks.MockStore
	ctx        context.Context
}

func (suite *VerifySolutionTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.mockStore = mocks.NewMockStore(ctrl)
	suite.challenger = NewChallenger(slog.Default(), suite.mockStore)
	suite.ctx = context.Background()
}

func (suite *VerifySolutionTestSuite) TestVerifySolution_ValidSolution() {
	t := suite.T()

	// Создаем тестовый challenge
	challenge := createTestChallenge("test123", 4, 10*time.Minute)

	// Находим валидный nonce
	validNonce := calculateValidNonce(challenge.Data, challenge.Difficulty)

	// Настраиваем мок
	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), "test123").Return(challenge, nil)

	// Выполняем тест
	err := suite.challenger.VerifySolution(suite.ctx, "test123", validNonce)

	// Проверки
	assert.NoError(t, err)
}

func (suite *VerifySolutionTestSuite) TestVerifySolution_InvalidNonce() {
	t := suite.T()

	challenge := createTestChallenge("test123", 4, 10*time.Minute)
	invalidNonce := uint64(999999) // Заведомо неправильный nonce

	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), "test123").Return(challenge, nil)

	err := suite.challenger.VerifySolution(suite.ctx, "test123", invalidNonce)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid solution")
}

func (suite *VerifySolutionTestSuite) TestVerifySolution_ChallengeNotFound() {
	t := suite.T()

	expectedError := errors.New("challenge not found")

	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), "nonexistent").Return(nil, expectedError)

	err := suite.challenger.VerifySolution(suite.ctx, "nonexistent", 12345)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get challenge")
}

func (suite *VerifySolutionTestSuite) TestVerifySolution_ExpiredChallenge() {
	t := suite.T()

	// Создаем истекший challenge
	expiredChallenge := createExpiredChallenge("expired123", 4)

	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), "expired123").Return(expiredChallenge, nil)

	err := suite.challenger.VerifySolution(suite.ctx, "expired123", 12345)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "challenge expired")
}

func (suite *VerifySolutionTestSuite) TestVerifySolution_StoreError() {
	t := suite.T()

	expectedError := errors.New("redis connection failed")

	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), "test123").Return(nil, expectedError)

	err := suite.challenger.VerifySolution(suite.ctx, "test123", 12345)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get challenge")
}

func (suite *VerifySolutionTestSuite) TestVerifySolution_ZeroDifficulty() {
	t := suite.T()

	// Challenge с нулевой сложностью должен приниматься с любым nonce
	challenge := createTestChallenge("easy123", 0, 10*time.Minute)

	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), "easy123").Return(challenge, nil)

	err := suite.challenger.VerifySolution(suite.ctx, "easy123", 0)

	assert.NoError(t, err)
}

func (suite *VerifySolutionTestSuite) TestVerifySolution_HighDifficulty() {
	t := suite.T()

	// Высокая сложность - найдем валидный nonce
	challenge := createTestChallenge("hard123", 8, 10*time.Minute)
	validNonce := calculateValidNonce(challenge.Data, challenge.Difficulty)

	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), "hard123").Return(challenge, nil)

	err := suite.challenger.VerifySolution(suite.ctx, "hard123", validNonce)

	assert.NoError(t, err)
}

func (suite *VerifySolutionTestSuite) TestVerifySolution_ChallengeAlmostExpired() {
	t := suite.T()

	// Challenge истекает через 1 секунду
	challenge := createTestChallenge("almostexpired123", 4, 1*time.Second)
	validNonce := calculateValidNonce(challenge.Data, challenge.Difficulty)

	suite.mockStore.EXPECT().GetChallenge(gomock.Any(), "almostexpired123").Return(challenge, nil)

	err := suite.challenger.VerifySolution(suite.ctx, "almostexpired123", validNonce)

	// Должен быть валидным, если еще не истек
	assert.NoError(t, err)
}

func TestVerifySolutionTestSuite(t *testing.T) {
	suite.Run(t, new(VerifySolutionTestSuite))
}
