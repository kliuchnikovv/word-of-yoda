package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/kliuchnikovv/word-of-yoda/internal/server"
	"github.com/kliuchnikovv/word-of-yoda/internal/server/config"
	"github.com/kliuchnikovv/word-of-yoda/internal/server/redis"
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

	listener, err := net.Listen("tcp", cfg.Server.Address)
	if err != nil {
		logger.Error("failed to setup listener", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer listener.Close()

	store, err := redis.NewRedisStore(logger, &redis.RedisConfig{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})
	if err != nil {
		logger.Error("failed to setup redis store", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer store.Close()

	server, err := server.New(logger, store, listener, cfg.Server.TTLS, cfg.Server.Difficulty, cfg.Server.TimeoutS)
	if err != nil {
		logger.Error("failed to setup server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go server.ListenAndServe(ctx)
	logger.Info("server started")

	sig := <-sigChan
	logger.Info("received signal", "signal", sig.String())

	os.Exit(0)
}
