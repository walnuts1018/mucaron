package postgres

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/walnuts1018/mucaron/backend/config"
)

const (
	user     = "postgres"
	password = "postgres"
)

var p *PostgresClient

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create pool: %v", err))
		os.Exit(1)
	}

	if err := pool.Client.Ping(); err != nil {
		slog.Error(fmt.Sprintf("failed to connect to Docker: %v", err))
		os.Exit(1)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get current directory: %v", err))
		os.Exit(1)
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
		slog.Error(fmt.Sprintf("failed to create pool: %v", err))
		os.Exit(1)
	}

	host, port, err := net.SplitHostPort(resource.GetHostPort("5432/tcp"))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to split host and port: %v", err))
		os.Exit(1)
	}

	if err := pool.Retry(func() error {
		cfg := config.Config{
			PSQLDSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", host, port, user, password, "mucaron_test", "disable", "Asia/Tokyo"),
		}

		var err error
		p, err = NewPostgres(cfg)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		slog.Error(fmt.Sprintf("failed to connect to database: %v", err))
		os.Exit(1)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			slog.Error(fmt.Sprintf("failed to purge resources: %v", err))
			os.Exit(1)
		}

	}()

	m.Run()
}
