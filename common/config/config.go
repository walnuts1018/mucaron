package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	ServerURL  string `env:"SERVER_URL" envDefault:"localhost"`

	LogLevelString string     `env:"LOG_LEVEL" envDefault:"info"`
	LogLevel       slog.Level // Parse from LogLevelString

	// --------------------- PostgreSQL ---------------------
	PSQLDSN      string `env:"PSQL_DSN" envDefault:""` // If PSQL_DSN is set, other PSQL_* variables will be ignored
	PSQLHost     string `env:"PSQL_HOST" envDefault:"localhost"`
	PSQLPort     string `env:"PSQL_PORT" envDefault:"5432"`
	PSQLDatabase string `env:"PSQL_DATABASE" envDefault:"ac_hacking"`
	PSQLUser     string `env:"PSQL_USER" envDefault:"postgres"`
	PSQLPassword string `env:"PSQL_PASSWORD" envDefault:"postgres"`
	PSQLSSLMode  string `env:"PSQL_SSL_MODE" envDefault:"disable"`
	PSQLTimeZone string `env:"PSQL_TIMEZONE" envDefault:"Asia/Tokyo"`
	// -------------------------------------------------------

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
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	cfg.parseLogLevel()
	cfg.makePSQLDSN()

	return cfg, nil
}

func (cfg *Config) parseLogLevel() {
	cfg.LogLevel = func() slog.Level {
		switch strings.ToLower(cfg.LogLevelString) {
		case "debug":
			return slog.LevelDebug
		case "info":
			return slog.LevelInfo
		case "warn":
			return slog.LevelWarn
		case "error":
			return slog.LevelError
		default:
			slog.Warn("Invalid log level, use default level: info")
			return slog.LevelInfo
		}
	}()
	
}

func (cfg *Config) makePSQLDSN() {
	if cfg.PSQLDSN != "" {
		return
	}

	cfg.PSQLDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", cfg.PSQLHost, cfg.PSQLPort, cfg.PSQLUser, cfg.PSQLPassword, cfg.PSQLDatabase, cfg.PSQLSSLMode, cfg.PSQLTimeZone)
}
