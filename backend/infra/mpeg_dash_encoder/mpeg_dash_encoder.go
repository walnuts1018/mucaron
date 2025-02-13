package mpegdashencoder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/walnuts1018/mucaron/backend/config"
)

type Controller struct {
	serverEndpoint url.URL
	adminToken     string

	minioClient             *minio.Client
	minioSourceUploadBucket string
}

func NewController(cfg config.Config) (*Controller, error) {
	endpoint, err := url.Parse(cfg.MpegDashServerEndpoint)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &Controller{
		serverEndpoint:          *endpoint,
		adminToken:              cfg.MpegDashAdminToken,
		minioClient:             minioClient,
		minioSourceUploadBucket: cfg.MinIOSourceUploadBucket,
	}, nil
}

func (c *Controller) GetUserToken(mediaIDs []uuid.UUID) (string, error) {

	var reqBody struct {
		MediaIDs []string `json:"media_ids"`
	}

	for _, id := range mediaIDs {
		reqBody.MediaIDs = append(reqBody.MediaIDs, id.String())
	}

	buff := new(bytes.Buffer)
	if err := json.NewEncoder(buff).Encode(reqBody); err != nil {
		return "", fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.serverEndpoint.JoinPath("/v1/admin/create_user_token").String(), buff)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.adminToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to request: %w", err)
	}

	var respBody struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	return respBody.Token, nil
}

func (c *Controller) UploadMedia(ctx context.Context, mediaID uuid.UUID, filePath string) error {
	_, err := c.minioClient.FPutObject(
		ctx,
		c.minioSourceUploadBucket,
		mediaID.String(),
		filePath,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to upload media: %w", err)
	}

	return nil
}
