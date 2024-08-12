package minio

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/walnuts1018/mucaron/common/config"
)

const (
	accessKey  = "mockaccesskey"
	secretKey  = "mocksecretkey"
	bucketName = "mucaron-test"
)

var hostAndPort string
var minioClient *MinIO

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

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "minio/minio",
			Tag:        "latest",
			Env: []string{
				fmt.Sprintf("MINIO_ACCESS_KEY=%s", accessKey),
				fmt.Sprintf("MINIO_SECRET_KEY=%s", secretKey),
			},
		},
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
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

	hostAndPort = resource.GetHostPort("9000/tcp")
	if err := pool.Retry(func() error {
		var err error
		minioClient, err = NewMinIOClient(config.Config{
			MinIOEndpoint:  hostAndPort,
			MinIOAccessKey: accessKey,
			MinIOSecretKey: secretKey,
			MinIOUseSSL:    false,
			MinIOBucket:    bucketName,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("failed to create minio client: %v", err))
			os.Exit(1)
		}

		cancel, err := minioClient.client.HealthCheck(10 * time.Second)
		if err != nil {
			return fmt.Errorf("failed to call healthcheck: %v", err)
		}
		defer cancel()

		online := minioClient.client.IsOnline()
		if online {
			return nil
		} else {
			return fmt.Errorf("minio is offline")
		}

	}); err != nil {
		slog.Error(fmt.Sprintf("failed to connect to minio: %v", err))
		os.Exit(1)
	}

	ctx := context.Background()
	bucketExist, err := minioClient.client.BucketExists(ctx, bucketName)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to check bucket: %v", err))
		os.Exit(1)
	}
	if !bucketExist {
		if err := minioClient.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			slog.Error(fmt.Sprintf("failed to create bucket: %v", err))
			os.Exit(1)
		}
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			slog.Error(fmt.Sprintf("failed to purge resources: %v", err))
			os.Exit(1)
		}

	}()

	m.Run()
}

func TestMinIO_GetObjectURL(t *testing.T) {
	type args struct {
		ctx          context.Context
		objectName   string
		cacheControl string
	}
	tests := []struct {
		name    string
		args    args
		wantURL string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:        context.Background(),
				objectName: "test/1",
			},
			wantURL: fmt.Sprintf("https://%s/%s/%s", hostAndPort, bucketName, "test/1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := minioClient.GetObjectURL(tt.args.ctx, tt.args.objectName, tt.args.cacheControl)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinIO.GetObjectURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotURL := got.String()

			if !reflect.DeepEqual(gotURL, tt.wantURL) {
				t.Errorf("MinIO.GetObjectURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}
