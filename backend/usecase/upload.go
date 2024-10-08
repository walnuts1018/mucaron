package usecase

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/util/filehash"
	"github.com/walnuts1018/mucaron/backend/util/temp"
)

const timeout = 1 * time.Minute

func (u *Usecase) UploadMusic(ctx context.Context, user entity.User, r io.Reader) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	tmpfile, err := temp.CreateTempFile(r, id.String())
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	file := tmpfile.UseFile()
	defer file.Close()
	raw, err := u.metadataReader.GetMetadata(ctx, file.Name())
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	music := raw.ToEntity(user)

	hash, err := filehash.FileHash(file.Name())
	if err != nil {
		return fmt.Errorf("failed to get file hash: %w", err)
	}
	music.ID = hash

	if err := u.MusicRepository.CreateMusic(music); err != nil {
		return fmt.Errorf("failed to create music: %w", err)
	}

	u.encodeQueue <- queueItem{
		ID:        id,
		Path:      file,
		AudioOnly: true,
	}

	return nil
}
