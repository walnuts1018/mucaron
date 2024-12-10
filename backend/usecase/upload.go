package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/util/hash"
	"github.com/walnuts1018/mucaron/backend/util/temp"
	"gorm.io/gorm"
)

const (
	uploadedExtension = ".mucaronuploaded"
	encodedExtension  = ".mucaronencoded"
)

func (u *Usecase) UploadMusic(ctx context.Context, user entity.User, r io.Reader, fileName string) (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate uuid: %w", err)
	}

	tmpfile, err := temp.CreateTempFile(r, fmt.Sprintf("%s%s", id.String(), uploadedExtension))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpfile.Close()

	raw, err := u.metadataReader.GetMetadata(ctx, tmpfile.Name())
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	music, album, artist, genre := raw.ToEntity(user, fileName)

	hash, err := hash.ReaderHash(tmpfile)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get file hash: %w", err)
	}

	music.ID = id
	music.FileHash = hash

	if err := u.entityRepository.Transaction(ctx, func(ctx context.Context) error {
		if artist != nil {
			existingArtist, err := u.entityRepository.GetArtistByName(ctx, user.ID, artist.Name)
			if err != nil {
				if !errors.Is(err, domain.ErrNotFound) {
					return fmt.Errorf("failed to get artist by name: %w", err)
				}
			}

			if existingArtist.ID != uuid.Nil {
				if err := u.entityRepository.UpdateArtist(ctx, *artist); err != nil {
					return fmt.Errorf("failed to update artist: %w", err)
				}
				music.Artists = []entity.Artist{*artist}
			} else {
				if err := u.entityRepository.CreateArtist(ctx, *artist); err != nil {
					return fmt.Errorf("failed to create artist: %w", err)
				}
				music.Artists = []entity.Artist{*artist}
			}
		}

		if genre != nil {
			existingGenre, err := u.entityRepository.GetGenreByName(ctx, user.ID, genre.Name)
			if err != nil {
				if !errors.Is(err, domain.ErrNotFound) {
					return fmt.Errorf("failed to get genre by name: %w", err)
				}
			}

			if existingGenre.ID != uuid.Nil {
				if err := u.entityRepository.UpdateGenre(ctx, *genre); err != nil {
					return fmt.Errorf("failed to update genre: %w", err)
				}
				music.Genre = genre
			} else {
				if err := u.entityRepository.CreateGenre(ctx, *genre); err != nil {
					return fmt.Errorf("failed to create genre: %w", err)
				}
				music.Genre = genre
			}

		}

		// ファイルハッシュが一致するMusicが存在しないことを確認
		{
			var existingMusic entity.Music
			if err := u.entityRepository.GetMusicByFileHash(ctx, user.ID, music.FileHash, &existingMusic); err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("failed to get music by filehash: %w", err)
				}
			}

			if existingMusic.ID != uuid.Nil {
				// 既存のMusicが論理削除されている場合は物理削除する
				// 論理削除されていない場合はエラーを返す
				if existingMusic.DeletedAt.Valid {
					if err := u.entityRepository.HardDeleteMusic(ctx, existingMusic); err != nil {
						return fmt.Errorf("failed to delete existing music: %w", err)
					}
				} else {
					return fmt.Errorf("music already exists: %w", domain.ErrAlreadyExists)
				}
			}
		}

		// Musicを新規作成
		if err := u.entityRepository.CreateMusic(ctx, music); err != nil {
			return fmt.Errorf("failed to create music: %w", err)
		}

		if album != nil {
			// // 同名のAlbumが存在しない場合は新規作成
			existingAlbums, err := u.entityRepository.GetAlbumsByNameAndArtist(ctx, user.ID, album.Name, *artist)
			if err != nil {
				if !errors.Is(err, domain.ErrNotFound) {
					return fmt.Errorf("failed to get album by name and artist: %w", err)
				}
			}

			album.Musics = []entity.Music{music}
			if len(existingAlbums) < 1 {
				if err := u.entityRepository.CreateAlbum(ctx, *album); err != nil {
					return fmt.Errorf("failed to create album: %w", err)
				}
			} else {
				if err := u.entityRepository.UpdateAlbum(ctx, *album); err != nil {
					return fmt.Errorf("failed to update album: %w", err)
				}
			}
		}

		return nil
	}); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create music: %w", err)
	}

	return music.ID, nil
}
