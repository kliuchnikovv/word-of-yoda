package challenge

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	mocks "github.com/kliuchnikovv/word-of-yoda/internal/server/redis/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ChallengerTestSuite группирует тесты
type ChallengerTestSuite struct {
	suite.Suite
	challenger *Challenger
	mockStore  *mocks.MockStore
	logger     *slog.Logger
	ctx        context.Context
}

// SetupTest запускается перед каждым тестом
func (suite *ChallengerTestSuite) SetupTest() {
	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctrl := gomock.NewController(suite.T())

	suite.mockStore = mocks.NewMockStore(ctrl)
	suite.challenger = NewChallenger(suite.logger, suite.mockStore)
	suite.ctx = context.Background()
}

func (suite *ChallengerTestSuite) TestNewChallenger() {
	t := suite.T()

	challenger := NewChallenger(suite.logger, suite.mockStore)

	assert.NotNil(t, challenger)
	assert.Equal(t, suite.logger, challenger.logger)
	assert.Equal(t, suite.mockStore, challenger.store)
}

func (suite *ChallengerTestSuite) TestGenerateChallenge_Success() {
	t := suite.T()

	difficulty := 4
	ttl := 5 * time.Minute

	// Настраиваем мок
	suite.mockStore.EXPECT().SaveChallenge(gomock.Any(), gomock.Any()).Return(nil)

	// Выполняем тест
	challenge, err := suite.challenger.GenerateChallenge(suite.ctx, difficulty, ttl)

	// Проверки
	require.NoError(t, err)
	require.NotNil(t, challenge)

	assert.NotEmpty(t, challenge.ID)
	assert.NotEmpty(t, challenge.Data)
	assert.Equal(t, difficulty, challenge.Difficulty)
	assert.True(t, challenge.ExpiresAt.After(time.Now()))
	assert.True(t, challenge.ExpiresAt.Before(time.Now().Add(ttl+time.Minute)))
}

func (suite *ChallengerTestSuite) TestGenerateChallenge_ZeroDifficulty() {
	t := suite.T()

	difficulty := 0
	ttl := 5 * time.Minute

	suite.mockStore.EXPECT().SaveChallenge(gomock.Any(), gomock.Any()).Return(nil)

	challenge, err := suite.challenger.GenerateChallenge(suite.ctx, difficulty, ttl)

	require.NoError(t, err)
	require.NotNil(t, challenge)
	assert.Equal(t, 0, challenge.Difficulty)
}

func (suite *ChallengerTestSuite) TestGenerateChallenge_HighDifficulty() {
	t := suite.T()

	difficulty := 16
	ttl := 10 * time.Minute

	suite.mockStore.EXPECT().SaveChallenge(gomock.Any(), gomock.Any()).Return(nil)

	challenge, err := suite.challenger.GenerateChallenge(suite.ctx, difficulty, ttl)

	require.NoError(t, err)
	require.NotNil(t, challenge)
	assert.Equal(t, difficulty, challenge.Difficulty)
}

func (suite *ChallengerTestSuite) TestGenerateChallenge_StoreError() {
	t := suite.T()

	difficulty := 4
	ttl := 5 * time.Minute
	expectedError := errors.New("redis connection failed")

	suite.mockStore.EXPECT().SaveChallenge(gomock.Any(), gomock.Any()).Return(expectedError)

	challenge, err := suite.challenger.GenerateChallenge(suite.ctx, difficulty, ttl)

	assert.Error(t, err)
	assert.Nil(t, challenge)
	assert.Contains(t, err.Error(), "redis connection failed")
}

func (suite *ChallengerTestSuite) TestGenerateChallenge_ShortTTL() {
	t := suite.T()

	difficulty := 4
	ttl := 1 * time.Second

	suite.mockStore.EXPECT().SaveChallenge(gomock.Any(), gomock.Any()).Return(nil)

	challenge, err := suite.challenger.GenerateChallenge(suite.ctx, difficulty, ttl)

	require.NoError(t, err)
	require.NotNil(t, challenge)

	// Проверяем, что время истечения установлено правильно
	expectedExpiry := time.Now().Add(ttl)
	assert.WithinDuration(t, expectedExpiry, challenge.ExpiresAt, 1*time.Second)
}

// Запускаем тесты
func TestChallengerTestSuite(t *testing.T) {
	suite.Run(t, new(ChallengerTestSuite))
}
