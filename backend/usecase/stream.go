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
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

const maxAgeSecond = 7 * 24 * 60 * 60 // 1 week

func (u *Usecase) GetPrimaryStreamM3U8URL(ctx context.Context, user entity.User, musicID uuid.UUID) (*url.URL, error) {
	m, err := u.entityRepository.GetMusicByID(musicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get music by id: %w", err)
	}
	if m.OwnerID != user.ID {
		slog.Warn("access denied", slog.String("music_id", musicID.String()), slog.String("access_user_id", user.ID.String()), slog.String("owner_id", m.OwnerID.String()))
		return nil, domain.ErrAccessDenied
	}
	return u.objectStorage.GetObjectURL(ctx, filepath.Join(musicID.String(), "primary.m3u8"), "")
}

func (u *Usecase) GetStreamM3U8(ctx context.Context, user entity.User, musicID uuid.UUID, streamID string) (string, error) {
	// streamIDの正規化
	streamID = strings.TrimSuffix(streamID, ".m3u8")

	m, err := u.entityRepository.GetMusicByID(musicID)
	if err != nil {
		return "", fmt.Errorf("failed to get music by id: %w", err)
	}

	if m.OwnerID != user.ID {
		slog.Warn("access denied", slog.String("music_id", musicID.String()), slog.String("access_user_id", user.ID.String()), slog.String("owner_id", m.OwnerID.String()))
		return "", domain.ErrAccessDenied
	}

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

		presignedURL, err := u.objectStorage.GetObjectURL(ctx, strings.TrimPrefix(url.Path, "/"), fmt.Sprintf("max-age=%d", maxAgeSecond))
		if err != nil {
			return "", fmt.Errorf("failed to get presigned URL: %w", err)
		}
		newContents += presignedURL.String() + "\n"
	}

	return newContents, nil
}
