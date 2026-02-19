package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/app"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/config"
	platformdb "github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/database"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel)
	pool, err := platformdb.NewConnectionPool(ctx, *cfg)
	if err != nil {
		log.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	application := app.New(ctx, cfg, log, pool)

	go func() {
		log.Info("server_started", "port", cfg.HTTP.Port)
		if err := application.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server_error", "error", err)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := application.Server.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful_shutdown_failed", "error", err)
		os.Exit(1)
	}

	log.Info("server_stopped")
}
