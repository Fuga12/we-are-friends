package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fuga12/we-are-friends/backend-go/internal/config"
	"github.com/Fuga12/we-are-friends/backend-go/internal/httpapi"
	"github.com/Fuga12/we-are-friends/backend-go/internal/media"
	"github.com/Fuga12/we-are-friends/backend-go/internal/migrate"
	"github.com/Fuga12/we-are-friends/backend-go/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect PostgreSQL", "error", err)
		os.Exit(1)
	}
	database := store.New(pool)
	defer database.Close()
	if err := database.Ping(ctx); err != nil {
		slog.Error("ping PostgreSQL", "error", err)
		os.Exit(1)
	}
	if err := migrate.Up(ctx, pool); err != nil {
		slog.Error("apply database schema", "error", err)
		os.Exit(1)
	}
	if err := database.BootstrapAdmin(ctx, cfg.AdminUser); err != nil {
		slog.Error("bootstrap administrator", "error", err)
		os.Exit(1)
	}
	mediaClient, err := media.New(cfg)
	if err != nil {
		slog.Error("configure MinIO", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.NewRouter(cfg, database, mediaClient),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("API server started", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("API server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
}
