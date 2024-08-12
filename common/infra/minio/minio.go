package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/walnuts1018/mucaron/common/config"
)

const expires = time.Second * 24 * 60 * 60

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

func (m MinIO) GetObjectURL(ctx context.Context, objectName string, cacheControl string) (*url.URL, error) {
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

func (m MinIO) UploadObject(ctx context.Context, objectName string, data io.Reader, size int64) error {
	if _, err := m.client.PutObject(ctx, m.minioBucket, objectName, data, size, minio.PutObjectOptions{}); err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}
	return nil
}

func (m MinIO) DeleteObject(ctx context.Context, objectName string) error {
	if err := m.client.RemoveObject(ctx, m.minioBucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}
