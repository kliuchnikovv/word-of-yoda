package client

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/kliuchnikovv/word-of-yoda/domain"
	"github.com/kliuchnikovv/word-of-yoda/internal/client/solver"
	"github.com/kliuchnikovv/word-of-yoda/internal/utils"
)

const defaultTimeout int64 = 3

type Client struct {
	logger *slog.Logger

	serverAddress string
	timeout       time.Duration
}

func New(logger *slog.Logger, address string, timeoutS int64) (*Client, error) {
	if len(address) == 0 {
		return nil, fmt.Errorf("address is empty")
	}

	if timeoutS <= 0 {
		timeoutS = defaultTimeout
	} else if logger == nil {
		logger = slog.Default()
	}

	return &Client{
		logger:        logger,
		serverAddress: address,
		timeout:       time.Duration(timeoutS) * time.Second,
	}, nil
}

func (client *Client) SolvePuzzles(ctx context.Context, sleep int64) error {
	logger := client.logger.With(
		slog.String("method", "SolvePuzzles"),
	)

	var ticker = time.NewTicker(time.Duration(sleep) * time.Second)

	for {
		select {
		case <-ctx.Done():
			logger.Info("stop solving")
			return nil
		case <-ticker.C:
			logger.Info("ready to get quote")

			quote, err := client.GetQuote(ctx)
			if err != nil {
				logger.Error("failed to get quote", slog.String("error", err.Error()))
			} else {
				logger.Info("got quote", quote.Log()...)
			}
		}
	}
}

func (client *Client) GetQuote(ctx context.Context) (*domain.Quote, error) {
	logger := client.logger.With(
		slog.String("method", "GetQuote"),
	)

	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", client.serverAddress)
	if err != nil {
		logger.Error("failed to dial", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	logger.Debug("connected")

	var (
		reader = bufio.NewReader(conn)
		writer = bufio.NewWriter(conn)
	)

	challenge, err := utils.ReadMessage[domain.Challenge](reader)
	if err != nil {
		logger.Error("failed to read challenge", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to read challenge: %w", err)
	}

	logger.Debug("got challenge", challenge.Log()...)

	solution, err := solver.Solve(ctx, logger, *challenge)
	if err != nil {
		logger.Error("failed to solve challenge", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to solve challenge: %w", err)
	}

	logger.Debug("solved challenge", solution.Log()...)

	if err := utils.WriteMessage(writer, solution); err != nil {
		logger.Error("failed to write solution", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to write solution: %w", err)
	}

	logger.Debug("send solution", solution.Log()...)

	quote, err := utils.ReadMessage[domain.Quote](reader)
	if err != nil {
		logger.Error("failed to read quote", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to read quote: %w", err)
	}

	logger.Debug("read quote", quote.Log()...)

	return quote, nil
}
