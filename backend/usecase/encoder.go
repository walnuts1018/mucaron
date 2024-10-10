package usecase

import (
	"context"
	"log/slog"

	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (u *Usecase) encode(ctx context.Context, music entity.Music, uploadedFilePath string, audioOnly bool) {
	path, err := u.encoder.Encode(music.ID.String(), uploadedFilePath, audioOnly)
	if err != nil {
		slog.Error("failed to encode", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		return
	}
	if err := u.objectStorage.UploadDirectory(ctx, music.ID.String(), path); err != nil {
		slog.Error("failed to upload directory", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		return
	}
}
