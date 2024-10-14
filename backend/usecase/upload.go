package usecase

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"regexp"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/util/hash"
	"github.com/walnuts1018/mucaron/backend/util/temp"
)

const (
	uploadedExtension = ".mucaronuploaded"
	encodedExtension  = ".mucaronencoded"
)

func (u *Usecase) UploadMusic(ctx context.Context, user entity.User, r io.Reader, fileName string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	tmpfile, err := temp.CreateTempFile(r, fmt.Sprintf("%s%s", id.String(), uploadedExtension))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpfile.Close()

	raw, err := u.metadataReader.GetMetadata(ctx, tmpfile.Name())
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	music, album, artist, genre := raw.ToEntity(user, fileName)

	hash, err := hash.ReaderHash(tmpfile)
	if err != nil {
		return fmt.Errorf("failed to get file hash: %w", err)
	}

	music.ID = id
	music.FileHash = hash

	if err := u.entityRepository.CreateMusicWithDependencies(music, album, artist, genre); err != nil {
		return fmt.Errorf("failed to create music: %w", err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), u.cfg.EncodeTimeout)
		defer cancel()

		u.encode(ctx, tmpfile.Name(), music)
	}()

	return nil
}

var re = regexp.MustCompile(`(stream_[\d]+.m3u8)`)

func replaceM3U8URL(content string, serverEndpoint, musicID string) (string, error) {
	newURL, err := url.JoinPath(serverEndpoint, "/api/v1/music/", musicID, "/stream/$1")
	if err != nil {
		return "", fmt.Errorf("failed to join url: %w", err)
	}
	replaced := re.ReplaceAllString(content, newURL)
	return replaced, nil
}
