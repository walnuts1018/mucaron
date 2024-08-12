package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
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

var (
	hostAndPort string
	minioClient *MinIO
	testData    = []byte("testdata")
)

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
			Cmd: []string{"server", "/export", "--console-address", ":9001"},
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

func TestMinIO_IOReader(t *testing.T) {
	type args struct {
		ctx          context.Context
		objectName   string
		cacheControl string
	}
	tests := []struct {
		name    string
		args    args
		wantURL string
	}{
		{
			name: "normal",
			args: args{
				ctx:        context.Background(),
				objectName: "test/1",
			},
			wantURL: fmt.Sprintf("http://%s/%s/%s", hostAndPort, bucketName, "test/1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buff bytes.Buffer
			buff.Write(testData)
			minioClient.UploadObject(tt.args.ctx, tt.args.objectName, &buff, int64(buff.Len()))

			got, err := minioClient.GetObjectURL(tt.args.ctx, tt.args.objectName, tt.args.cacheControl)
			if err != nil {
				t.Errorf("MinIO.GetObjectURL() error = %v", err)
				return
			}
			slog.Info("ObjectURL", slog.String(
				"url", got.String(),
			))

			gotWithoutQuery := *got
			gotWithoutQuery.RawQuery = ""

			if !reflect.DeepEqual(gotWithoutQuery.String(), tt.wantURL) {
				t.Errorf("MinIO.GetObjectURL() = %v, want %v", gotWithoutQuery.String(), tt.wantURL)
			}

			resp, err := http.Get(got.String())
			if err != nil {
				t.Errorf("failed to get: %v", err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
			}

			if !reflect.DeepEqual(body, testData) {
				t.Errorf("got = %s, want %s", body, testData)
			}

			if err := minioClient.DeleteObject(tt.args.ctx, tt.args.objectName); err != nil {
				t.Errorf("MinIO.DeleteObject() error = %v", err)
				return
			}
		})
	}
}

func TestMinIO_Directory(t *testing.T) {
	type args struct {
		ctx           context.Context
		objectBaseDir string
		localDir      string
		cacheControl  string
	}
	tests := []struct {
		name    string
		args    args
		wantDir string
		wantURL string
	}{
		{
			name: "normal",
			args: args{
				ctx:           context.Background(),
				objectBaseDir: "test_directory",
				localDir:      "tmp/1",
			},
			wantDir: "test_directory/files",
			wantURL: fmt.Sprintf("http://%s/%s/%s", hostAndPort, bucketName, "test_directory/files/test_0.txt"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "mucaron-test_directory")
			if err != nil {
				t.Errorf("failed to create tmp directory error = %v", err)
				return
			}

			dirPath := path.Join(tmpDir, tt.args.localDir)
			if err := os.MkdirAll(path.Join(dirPath, "files"), 0755); err != nil {
				t.Errorf("failed to create directory error = %v", err)
				return
			}

			for i := 0; i < 3; i++ {
				file, err := os.Create(path.Join(dirPath, fmt.Sprintf("files/test_%d.txt", i)))
				if err != nil {
					t.Errorf("failed to create file: %v", err)
					return
				}
				defer file.Close()

				_, err = file.Write(testData)
				if err != nil {
					t.Errorf("failed to write file: %v", err)
					return
				}
			}

			if err := minioClient.UploadDirectory(tt.args.ctx, tt.args.objectBaseDir, dirPath); err != nil {
				t.Errorf("MinIO.UploadDirectory() error = %v", err)
				return
			}

			got, err := minioClient.GetObjectURL(tt.args.ctx, fmt.Sprintf("%v/test_0.txt", tt.wantDir), tt.args.cacheControl)
			if err != nil {
				t.Errorf("MinIO.GetObjectURL() error = %v", err)
				return
			}
			slog.Info("ObjectURL", slog.String(
				"url", got.String(),
			))

			gotWithoutQuery := *got
			gotWithoutQuery.RawQuery = ""

			if !reflect.DeepEqual(gotWithoutQuery.String(), tt.wantURL) {
				t.Errorf("MinIO.GetObjectURL() = %v, want %v", gotWithoutQuery.String(), tt.wantURL)
				return
			}

			resp, err := http.Get(got.String())
			if err != nil {
				t.Errorf("failed to get: %v", err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
				return
			}

			if !reflect.DeepEqual(body, testData) {
				t.Errorf("got = %s, want %s", body, testData)
				return
			}

			if err := minioClient.DeleteObject(tt.args.ctx, tt.args.objectBaseDir); err != nil {
				t.Errorf("MinIO.DeleteObject() error = %v", err)
				return
			}
		})
	}
}
