package usecase

import (
	"context"
	"log/slog"
)

func (u *Usecase) StartEncoder(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case item := <-u.encodeQueue:
			go func(item queueItem) {
				path, err := u.encoder.Encode(item.ID, item.Path, item.AudioOnly)
				if err != nil {
					slog.Error("failed to encode", slog.Any("id", item.ID), slog.Any("error", err))
					return
				}
				if err := u.objectStorage.UploadDirectory(ctx, item.ID.String(), path); err != nil {
					slog.Error("failed to upload directory", slog.Any("id", item.ID), slog.Any("error", err))
					return
				}
			}(item)
		}
	}
}
