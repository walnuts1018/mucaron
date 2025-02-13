package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/walnuts1018/mucaron/backend/util/random"
)

var ErrInvalidSessionSecretLength = errors.New("session secret must be 16, 24, or 32 bytes")

type Config struct {
	ServerPort     string     `env:"SERVER_PORT" envDefault:"8080"`
	ServerEndpoint string     `env:"SERVER_ENDPOINT" envDefault:"http://localhost:80"`
	LogLevel       slog.Level `env:"LOG_LEVEL"`
	LogType        LogType    `env:"LOG_TYPE" envDefault:"json"`

	PSQLDSN string `env:"PSQL_DSN" envDefault:""` // If PSQL_DSN is set, other PSQL_* variables will be ignored

	// ------------------------ MinIO ------------------------
	MinIOEndpoint           string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	MinIOAccessKey          string `env:"MINIO_ACCESS_KEY,required"`
	MinIOSecretKey          string `env:"MINIO_SECRET_KEY,required"`
	MinIOUseSSL             bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MinIOSourceUploadBucket string `env:"MINIO_SOURCE_UPLOAD_BUCKET" envDefault:"mpeg-dash-encoder-source-upload"`

	// ------------------------ Redis ------------------------
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`

	// ------------------------ MpegDash ------------------------
	MpegDashServerEndpoint string `env:"MPEG_DASH_SERVER_ENDPOINT" envDefault:"http://localhost:8081"`
	MpegDashAdminToken     string `env:"MPEG_DASH_ADMIN_TOKEN,required"`

	SessionSecret  string `env:"SESSION_SECRET" envDefault:""`
	SessionOptions SessionOptions
}

type SessionOptions struct {
	SameSite http.SameSite `env:"SESSION_SAME_SITE" envDefault:"strict"`
	HttpOnly bool          `env:"SESSION_HTTP_ONLY" envDefault:"true"`
	Secure   bool          `env:"SESSION_SECURE" envDefault:"true"`
}

func Load() (Config, error) {
	var cfg Config
	var parseErr error
	if err := env.ParseWithOptions(&cfg, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(slog.Level(0)):    returnAny(ParseLogLevel),
			reflect.TypeOf(time.Duration(0)): returnAny(time.ParseDuration),
			reflect.TypeOf(http.SameSite(0)): returnAny(ParseSameSite),
			reflect.TypeOf(LogType("")):      returnAny(ParseLogType),
		},
		OnSet: func(tag string, value any, isDefault bool) {
			switch tag {
			case "PSQL_DSN":
				if isDefault {
					dsn, err := parsePSQLSettings()
					if err != nil {
						parseErr = err
						return
					}
					cfg.PSQLDSN = dsn
				} else {
					dsn, ok := value.(string)
					if !ok {
						parseErr = errors.New("PSQL_DSN must be string")
						return
					}
					cfg.PSQLDSN = dsn
				}
			case "SESSION_SECRET":
				str, ok := value.(string)
				if !ok {
					parseErr = errors.New("SESSION_SECRET must be string")
					return
				}
				v, err := parseSessionSecret(str)
				if err != nil {
					parseErr = err
					return
				}
				cfg.SessionSecret = v
			}
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

func parseSessionSecret(v string) (string, error) {
	if v == "" {
		str, err := random.String(32, random.Alphanumeric)
		if err != nil {
			return "", err
		}
		return str, nil
	} else {
		allowedLen := []int{16, 24, 32}
		for _, l := range allowedLen {
			if len(v) == l {
				return v, nil
			}
		}
		return "", ErrInvalidSessionSecretLength
	}
}

type LogType string

const (
	LogTypeJSON LogType = "json"
	LogTypeText LogType = "text"
)

func ParseLogType(v string) (LogType, error) {
	switch strings.ToLower(v) {
	case "json":
		return LogTypeJSON, nil
	case "text":
		return LogTypeText, nil
	default:
		slog.Warn("Invalid log type, use default type: json")
		return LogTypeJSON, nil
	}
}

func ParseSameSite(v string) (http.SameSite, error) {
	switch strings.ToLower(v) {
	case "strict":
		return http.SameSiteStrictMode, nil
	case "lax":
		return http.SameSiteLaxMode, nil
	case "none":
		return http.SameSiteNoneMode, nil
	default:
		slog.Warn("Invalid SameSite, use default mode strict")
		return http.SameSiteStrictMode, nil
	}
}
