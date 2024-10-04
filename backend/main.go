package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain/logger"
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

	logger := slog.New(logger.NewTraceHandler(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     cfg.LogLevel,
			AddSource: cfg.LogLevel == slog.LevelDebug,
		}),
	))
	slog.SetDefault(logger)

	router, err := wire.CreateRouter(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create router: %v", err))
		os.Exit(1)
	}

	go func() {
		if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
			slog.Error(fmt.Sprintf("Failed to run router: %v", err))
			os.Exit(1)
		}
	}()
}
