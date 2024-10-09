package minio

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/walnuts1018/mucaron/backend/config"
)

const expires = 2 * 24 * time.Hour

type MinIO struct {
	minioBucket    string
	minioPublicUrl url.URL
	client         *minio.Client
}

func NewMinIOClient(cfg config.Config) (*MinIO, error) {
	minioClient, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	url, err := url.Parse(cfg.MinIOPublicBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse minio public url: %w", err)
	}

	return &MinIO{
		minioBucket:    cfg.MinIOBucket,
		minioPublicUrl: *url,
		client:         minioClient,
	}, nil
}

func (m *MinIO) GetObjectURL(ctx context.Context, objectName string, cacheControl string) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment")
	if cacheControl != "" {
		reqParams.Set("response-cache-control", cacheControl)
	}

	url, err := m.client.PresignedGetObject(ctx, m.minioBucket, objectName, expires, reqParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get presigned URL: %w", err)
	}
	return url, nil
}

func (m *MinIO) UploadObject(ctx context.Context, objectName string, data io.Reader, size int64) error {
	if _, err := m.client.PutObject(ctx, m.minioBucket, objectName, data, size, minio.PutObjectOptions{}); err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}
	return nil
}

func (m *MinIO) UploadDirectory(ctx context.Context, objectBaseDir string, localDir string) error {
	if err := filepath.WalkDir(localDir, func(localFilePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}
		localRelativeFilePath, err := filepath.Rel(localDir, localFilePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		objectPath := path.Join(objectBaseDir, strings.ReplaceAll(localRelativeFilePath, "\\", "/"))

		if _, err := m.client.FPutObject(ctx, m.minioBucket, objectPath, localFilePath, minio.PutObjectOptions{}); err != nil {
			return fmt.Errorf("failed to put object: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk dir: %w", err)
	}

	return nil
}

func (m *MinIO) DeleteObject(ctx context.Context, objectName string) error {
	if err := m.client.RemoveObject(ctx, m.minioBucket, objectName, minio.RemoveObjectOptions{
		ForceDelete: true, // recursive
	}); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}
