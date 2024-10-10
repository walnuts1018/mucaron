package usecase

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/util/hash"
	"github.com/walnuts1018/mucaron/backend/util/temp"
)

func (u *Usecase) UploadMusic(ctx context.Context, user entity.User, r io.Reader) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	tmpfile, err := temp.CreateTempFile(r, id.String())
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	raw, err := u.metadataReader.GetMetadata(ctx, tmpfile.File().Name())
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	music := raw.ToEntity(user)

	hash, err := hash.ReaderHash(tmpfile.File())
	if err != nil {
		return fmt.Errorf("failed to get file hash: %w", err)
	}

	music.FileHash = hash

	if err := u.MusicRepository.CreateMusic(music); err != nil {
		return fmt.Errorf("failed to create music: %w", err)
	}

	go func() {
		defer tmpfile.Close()

		ctx, cancel := context.WithTimeout(context.Background(), u.cfg.EncodeTimeout)
		defer cancel()

		u.encodeMutex.Lock()
		defer u.encodeMutex.Unlock()

		slog.Info("start encoding", slog.String("music_id", music.ID.String()))
		u.encode(ctx, music, tmpfile.File().Name(), false)
		slog.Info("finish encoding", slog.String("music_id", music.ID.String()))
	}()

	return nil
}
