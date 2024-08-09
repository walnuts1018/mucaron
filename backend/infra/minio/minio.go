package minio

import (
	"fmt"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/walnuts1018/mucaron/config"
)

type client struct {
	minioBucket    string
	minioPublicUrl url.URL
	mc             *minio.Client
}

func NewMinIOClient(cfg config.Config) (domain.IconRepository, error) {
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	url, err := url.Parse(cfg.MinioPublicURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse minio public url: %w", err)
	}

	return &client{
		minioBucket:    cfg.MinioBucket,
		minioPublicUrl: *url,
		mc:             minioClient,
	}, nil
}
