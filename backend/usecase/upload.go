package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/util/filehash"
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

	hash, err := filehash.FileHash(tmpfile.File().Name())
	if err != nil {
		return fmt.Errorf("failed to get file hash: %w", err)
	}

	music.ID = hash

	if err := u.MusicRepository.CreateMusic(music); err != nil {
		return fmt.Errorf("failed to create music: %w", err)
	}

	go func() {
		defer tmpfile.Close()
		
		u.encodeMutex.Lock()
		defer u.encodeMutex.Unlock()

		u.encode(ctx, music, tmpfile.File().Name(), false)
	}()

	return nil
}
