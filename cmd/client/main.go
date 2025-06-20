package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kliuchnikovv/word-of-yoda/internal/client"
	"github.com/kliuchnikovv/word-of-yoda/internal/client/config"
)

var cfg *config.Config

func init() {
	path, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		panic("config path not set up")
	}

	var err error
	if cfg, err = config.New(path); err != nil {
		panic(err)
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     cfg.Logger.Level,
		AddSource: cfg.Logger.AddSource,
	}))

	client, err := client.New(logger, cfg.Server, cfg.Client.TimeoutS)
	if err != nil {
		logger.Error("failed to create client", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go client.SolvePuzzles(ctx, cfg.Client.SleepS)
	logger.Info("client started")

	sig := <-sigChan
	logger.Info("received signal", "signal", sig.String())

	os.Exit(0)
}
