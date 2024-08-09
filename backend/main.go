package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/walnuts1018/mucaron/config"
	"github.com/walnuts1018/mucaron/wire"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load config: %v", err))
		os.Exit(1)
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.RFC3339,
		Level:      cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	router, err := wire.CreateRouter(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create router: %v", err))
		os.Exit(1)
	}

	// サーバー起動
	go func() {
		if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
			slog.Error(fmt.Sprintf("Failed to run router: %v", err))
			os.Exit(1)
		}
	}()
}
