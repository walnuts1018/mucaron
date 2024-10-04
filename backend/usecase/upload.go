package usecase

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/util/filehash"
	"github.com/walnuts1018/mucaron/backend/util/result"
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

	raw, err := u.metadataReader.GetMetadata(ctx, tmpfile.Name())
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	music, artist, album, genre := raw.ToEntity()

	hash, err := filehash.FileHash(tmpfile.Name())
	if err != nil {
		return fmt.Errorf("failed to get file hash: %w", err)
	}
	music.Hash = hash

	ch := make(chan result.Result[string])
	go func(ch chan<- result.Result[string]) {
		path, err := u.encoder.Encode(id, r, false)
		ch <- result.Result[string]{
			Result: path,
			Error:  err,
		}
	}(ch)

	if err := u.MusicRepository.CreateMusic(m); err != nil {
		return fmt.Errorf("failed to create music: %w", err)
	}
}
