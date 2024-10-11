package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/walnuts1018/mucaron/backend/config"
	"gorm.io/gorm/logger"
)

type Logger struct {
	cfg config.Config
}

func NewLogger(cfg config.Config) *Logger {
	return &Logger{
		cfg: cfg,
	}
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	slog.Warn("Log Level Change is not supported")
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	slog.Warn("Trace is not supported")
	logger.Default.Trace(ctx, begin, fc, err)
}
