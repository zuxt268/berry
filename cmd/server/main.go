package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/di"
	"github.com/zuxt268/berry/internal/interface/httpserver"
)

func main() {

	// DI コンテナ初期化
	container, err := di.NewContainer()
	if err != nil {
		slog.Error("failed to initialize DI container", "error", err)
		os.Exit(1)
	}

	// Server に依存を注入
	srv := httpserver.NewServer(
		":"+config.Env.AppPort,
		container,
	)

	// Create error channel to capture server errors
	errChan := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		slog.Info("starting server", "port", config.Env.AppPort)
		if err := srv.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-errChan:
		slog.Error("server error", "error", err)
	case sig := <-sigChan:
		slog.Info("received signal", "signal", sig)
	}

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("error during server shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server shutdown complete")
}
