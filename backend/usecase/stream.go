package usecase

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (u *Usecase) GetStreamM3U8(ctx context.Context, musicID uuid.UUID, streamID string) (string, error) {
	// streamIDの正規化
	streamID = strings.TrimSuffix(streamID, ".m3u8")

	base, err := u.objectStorage.GetObject(ctx, filepath.Join(musicID.String(), fmt.Sprintf("%s.m3u8", streamID)))
	if err != nil {
		return "", fmt.Errorf("failed to get m3u8 file: %w", err)
	}
	defer base.Close()

	// m3u8ファイルの内容を取得
	content, err := io.ReadAll(base)
	if err != nil {
		return "", fmt.Errorf("failed to read m3u8 file: %w", err)
	}

	var newContents string
	for _, line := range strings.Split(string(content), "\n") {
		if line == "" {
			newContents += "\n"
			continue
		}

		if strings.HasPrefix(line, "#") {
			newContents += line + "\n"
			continue
		}

		url, err := url.Parse(line)
		if err != nil {
			return "", fmt.Errorf("failed to parse URL: %w", err)
		}
		slog.Debug("url", slog.String("url", url.String()))

		presignedURL, err := u.objectStorage.GetObjectURL(ctx, url.Path, "")
		if err != nil {
			return "", fmt.Errorf("failed to get presigned URL: %w", err)
		}
		newContents += presignedURL.String() + "\n"
	}

	return newContents, nil
}
