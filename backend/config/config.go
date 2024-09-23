package config

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ServerPort string     `env:"SERVER_PORT" envDefault:"8080"`
	ServerURL  string     `env:"SERVER_URL" envDefault:"localhost"`
	LogLevel   slog.Level `env:"LOG_LEVEL"`
	PSQLDSN    string     `env:"PSQL_DSN" envDefault:"invalid_value"` // If PSQL_DSN is set, other PSQL_* variables will be ignored

	// ------------------------ MinIO ------------------------
	MinIOEndpoint      string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	MinIOAccessKey     string `env:"MINIO_ACCESS_KEY,required"`
	MinIOSecretKey     string `env:"MINIO_SECRET_KEY,required"`
	MinIOUseSSL        bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MinIOBucket        string `env:"MINIO_BUCKET" envDefault:"mucaron"`
	MinIOPublicBaseURL string `env:"MINIO_PUBLIC_BASE_URL" envDefault:"http://localhost:9000"`
	// -------------------------------------------------------
}

func Load() (Config, error) {
	var cfg Config
	var parseErr error
	if err := env.ParseWithOptions(&cfg, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(slog.Level(0)):    returnAny(ParseLogLevel),
			reflect.TypeOf(time.Duration(0)): returnAny(time.ParseDuration),
		},
		OnSet: func(tag string, value any, isDefault bool) {
			if !isDefault {
				return
			}
			dsn, err := parsePSQLSettings()
			if err != nil {
				parseErr = err
			}
			cfg.PSQLDSN = dsn
		},
	}); err != nil {
		return Config{}, err
	}
	if parseErr != nil {
		return Config{}, parseErr
	}
	return cfg, nil
}

func returnAny[T any](f func(v string) (t T, err error)) func(v string) (any, error) {
	return func(v string) (any, error) {
		t, err := f(v)
		return any(t), err
	}
}

func ParseLogLevel(v string) (slog.Level, error) {
	switch strings.ToLower(v) {
	case "":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		slog.Warn("Invalid log level, use default level: info")
		return slog.LevelInfo, nil
	}
}

type PSQLSettings struct {
	PSQLHost     string `env:"PSQL_HOST" envDefault:"localhost"`
	PSQLPort     string `env:"PSQL_PORT" envDefault:"5432"`
	PSQLDatabase string `env:"PSQL_DATABASE" envDefault:"mucaron"`
	PSQLUser     string `env:"PSQL_USER" envDefault:"postgres"`
	PSQLPassword string `env:"PSQL_PASSWORD" envDefault:"postgres"`
	PSQLSSLMode  string `env:"PSQL_SSL_MODE" envDefault:"disable"`
	PSQLTimeZone string `env:"PSQL_TIMEZONE" envDefault:"Asia/Tokyo"`
}

func parsePSQLSettings() (string, error) {
	var s PSQLSettings
	if err := env.Parse(&s); err != nil {
		return "", err
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", s.PSQLHost, s.PSQLPort, s.PSQLUser, s.PSQLPassword, s.PSQLDatabase, s.PSQLSSLMode, s.PSQLTimeZone),
		nil
}
