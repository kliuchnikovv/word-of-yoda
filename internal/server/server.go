package server

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/kliuchnikovv/word-of-yoda/internal/server/redis"

	"github.com/kliuchnikovv/word-of-yoda/domain"
	"github.com/kliuchnikovv/word-of-yoda/internal/server/challenge"
	"github.com/kliuchnikovv/word-of-yoda/internal/server/quote"
	"github.com/kliuchnikovv/word-of-yoda/internal/utils"
)

const (
	defaultDifficulty   = 5
	defaultTTL          = 3600
	defaultServeTimeout = 10
)

type Server struct {
	logger   *slog.Logger
	listener net.Listener
	group    sync.WaitGroup

	challenger *challenge.Challenger

	ttl          time.Duration
	difficulty   int
	serveTimeout time.Duration

	connections chan net.Conn
	errors      chan error
}

func New(
	logger *slog.Logger,
	store redis.Store,
	listener net.Listener,
	ttlS int64,
	difficulty int,
	serveTimeoutS int64,
) (*Server, error) {
	if difficulty <= 0 {
		difficulty = defaultDifficulty
	}

	if ttlS <= 0 {
		ttlS = defaultTTL
	}

	if serveTimeoutS <= 0 {
		serveTimeoutS = defaultServeTimeout
	}

	if logger == nil {
		logger = slog.Default()
	}

	if store == nil {
		return nil, fmt.Errorf("redis store is nil")
	}

	if listener == nil {
		return nil, fmt.Errorf("listener is nil")
	}

	return &Server{
		logger:      logger,
		listener:    listener,
		challenger:  challenge.NewChallenger(logger, store),
		connections: make(chan net.Conn, 10),
		errors:      make(chan error, 10),

		ttl:          time.Duration(ttlS) * time.Second,
		serveTimeout: time.Duration(serveTimeoutS) * time.Second,
		difficulty:   difficulty,
	}, nil
}

func (server *Server) ListenAndServe(ctx context.Context) error {
	logger := server.logger.With(
		slog.String("method", "ListenAndServe"),
	)

	go server.listen(ctx)

	var shouldWork = true
	for shouldWork {
		select {
		case conn, ok := <-server.connections:
			if !ok {
				logger.Error("connections channel closed")
				return fmt.Errorf("connections channel closed")
			}

			logger.Debug("got connection", "addr", conn.RemoteAddr())

			server.group.Add(1)
			go server.serve(context.Background(), conn)
		case err, ok := <-server.errors:
			if !ok {
				logger.Error("errors channel closed")
				return fmt.Errorf("errors channel closed")
			}

			logger.Error("got error", slog.String("error", err.Error()))
		case <-ctx.Done():
			logger.Info("stopping...")
			shouldWork = false
		}
	}

	server.group.Wait()
	logger.Info("stopped...")

	return nil
}

func (server *Server) listen(ctx context.Context) {
	logger := server.logger.With(
		slog.String("method", "listen"),
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Debug("waiting for connection")

			conn, err := server.listener.Accept()
			if err != nil {
				logger.Error("got error")
				server.errors <- err
			} else {
				logger.Debug("got connection")
				server.connections <- conn
			}
		}
	}
}

func (server *Server) serve(ctx context.Context, conn net.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), server.serveTimeout)

	defer cancel()
	defer server.group.Done()
	defer conn.Close()

	var (
		logger = server.logger.With(slog.String("method", "serve"))
		reader = utils.NewReader[domain.Solution](ctx, bufio.NewReader(conn), 0)
		writer = bufio.NewWriter(conn)
	)

	challenge, err := server.challenger.GenerateChallenge(ctx, server.difficulty, server.ttl)
	if err != nil {
		logger.Error("can't generate challenge", slog.String("error", err.Error()))
		return
	}

	if err := utils.WriteMessage(writer, challenge); err != nil {
		logger.Error("can't write message", slog.String("error", err.Error()))
		return
	}

	msgCh, errCh := reader.Start()
	select {
	case msg := <-msgCh:
		logger.Debug("got message", "message", msg)

		if err := server.challenger.VerifySolution(ctx, msg.ID, msg.Nonce); err != nil {
			logger.Error("can't verify solution", slog.String("error", err.Error()))
			return
		}
	case err := <-errCh:
		logger.Error("can't read message", slog.String("error", err.Error()))
		return
	case <-ctx.Done():
		logger.Error("request timeout exceeded")
		return
	}

	var quote = quote.GetRandomQuote()
	if err := utils.WriteMessage(writer, quote); err != nil {
		logger.Error("can't write message", slog.String("error", err.Error()))
		return
	}

	logger.Info("request handled", "quote_id", quote.ID)
}
