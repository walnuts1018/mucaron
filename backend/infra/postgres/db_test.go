package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain/logger"
)

const (
	user     = "postgres"
	password = "postgres"
	dbname   = "mucaron_test"
)

var p *PostgresClient

func TestMain(m *testing.M) {
	if err := setupTest(m); err != nil {
		slog.Error("failed to setup test", slog.Any("error", err))
		os.Exit(1)
	}
}

func TestPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Suite")
}

func setupTest(m *testing.M) error {
	logger.CreateAndSetLogger(slog.LevelDebug, config.LogTypeText)

	pool, err := dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Client.Ping(); err != nil {
		return fmt.Errorf("failed to connect to Docker: %w", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	initPath := filepath.Join(currentDir, "..", "..", "..", "psql", "init_test")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "16",
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", user),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
			},
			Mounts: []string{
				fmt.Sprintf("%s:/docker-entrypoint-initdb.d", initPath),
			},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}
	defer func() {
		if err := pool.Purge(resource); err != nil {
			slog.Error("failed to purge resources", slog.Any("error", err))
		}
	}()

	host, port, err := net.SplitHostPort(resource.GetHostPort("5432/tcp"))
	if err != nil {
		return fmt.Errorf("failed to split host and port: %w", err)
	}

	if err := pool.Retry(func() error {
		cfg := config.Config{
			LogLevel: slog.LevelDebug,
			PSQLDSN:  fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", host, port, user, password, dbname, "disable", "Asia/Tokyo"),
		}

		var err error
		p, err = NewPostgres(context.Background(), cfg)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	code := m.Run()
	if code != 0 {
		return fmt.Errorf("test failed: %d", code)
	}

	return nil
}
