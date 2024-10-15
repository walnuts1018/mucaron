package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain/logger"
	"github.com/walnuts1018/mucaron/backend/tracer"
	"github.com/walnuts1018/mucaron/backend/wire"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger.CreateAndSetLogger(cfg.LogLevel, cfg.LogType)

	ctx := context.Background()
	close, err := tracer.NewTracerProvider(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create tracer provider: %v", err))
	}
	defer close()

	router, err := wire.CreateRouter(cfg)
	if err != nil {
		slog.Error("Failed to create router", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Server is running", slog.String("port", cfg.ServerPort))
	if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		slog.Error("Failed to run server", slog.Any("error", err))
		os.Exit(1)
	}
}
