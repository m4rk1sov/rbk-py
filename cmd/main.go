package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/m4rk1sov/rbk-py/internal/config"
	
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file")
	}
	cfg := config.Load()
	log := setupLogger(cfg.Env)
	
	srv := httpserver.New(cfg)
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start the server")
			panic(err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	
	return log
}
