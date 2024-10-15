package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *PostgresClient) CreateMusic(ctx context.Context, m entity.Music) error {
	result := p.DB(ctx).Create(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateMusic(ctx context.Context, m entity.Music) error {
	result := p.DB(ctx).Save(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateMusicStatus(ctx context.Context, musicID uuid.UUID, status entity.MusicStatus) error {
	result := p.DB(ctx).Model(&entity.Music{}).Where("id = ?", musicID).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update music status: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteMusics(ctx context.Context, musicIDs []uuid.UUID) error {
	var deleted []entity.Music
	result := p.DB(ctx).Select(clause.Associations).Select("RawMetaData").Delete(&deleted, musicIDs)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	slog.Debug("deleted musics", slog.Any("deleted", deleted))
	return nil
}

func (p *PostgresClient) GetMusicByID(ctx context.Context, id uuid.UUID) (entity.Music, error) {
	var m entity.Music
	result := p.DB(ctx).Preload("Artists").First(&m, id)
	if result.Error != nil {
		return entity.Music{}, fmt.Errorf("failed to get music by id: %w", result.Error)
	}
	return m, nil
}

func (p *PostgresClient) GetMusicByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Music, error) {
	m := make([]entity.Music, 0, len(ids))
	result := p.DB(ctx).Preload("Artists").Find(&m, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music by ids: %w", result.Error)
	}
	return m, nil
}

func (p *PostgresClient) GetMusicsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Music, error) {
	m := make([]entity.Music, 0)
	result := p.DB(ctx).Preload("Artists").Where("owner_id = ?", userID).Find(&m)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music by user id: %w", result.Error)
	}
	return m, nil
}

func (p *PostgresClient) GetMusicIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	result := p.DB(ctx).Model(&entity.Music{}).Where("owner_id = ?", userID).Pluck("id", &ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music ids by user id: %w", result.Error)
	}
	return ids, nil
}

func (p *PostgresClient) CreateMusicWithDependencies(ctx context.Context, userID uuid.UUID, m entity.Music, album *entity.Album, artist *entity.Artist, genre *entity.Genre) error {
	return p.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if artist != nil {
			// 同名のArtistが存在しない場合は新規作成
			if err := tx.Where("owner_id = ? AND name = ?", userID, artist.Name).Attrs(*artist).FirstOrCreate(artist).Error; err != nil {
				return fmt.Errorf("failed to create artist: %w", err)
			}
			m.Artists = []entity.Artist{*artist}
		}

		if genre != nil {
			// 同名のGenreが存在しない場合は新規作成
			if err := tx.Where("owner_id = ? AND name = ?", userID, genre.Name).Attrs(*genre).FirstOrCreate(genre).Error; err != nil {
				return fmt.Errorf("failed to create genre: %w", err)
			}
			m.GenreID = &genre.ID
		}

		var existingMusic entity.Music
		if err := tx.Unscoped().Preload(clause.Associations).Where("owner_id = ? AND file_hash = ?", userID, m.FileHash).First(&existingMusic).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to get music by filehash: %w", err)
			}
		}

		if existingMusic.ID != uuid.Nil {
			if existingMusic.DeletedAt.Valid {
				result := tx.Exec("DELETE FROM album_musics WHERE music_id = ?", existingMusic.ID).Exec("DELETE FROM raw_music_metadata_tags WHERE raw_music_metadata_id = ?", existingMusic.RawMetaData.ID).Unscoped().Select(clause.Associations).Delete(&existingMusic)
				if result.Error != nil {
					return fmt.Errorf("failed to delete existing music: %w", result.Error)
				}
			} else {
				return fmt.Errorf("music already exists: %w", domain.ErrAlreadyExists)
			}
		}

		// Musicを新規作成
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("failed to create music: %w", err)
		}

		if album != nil {
			album.Musics = []entity.Music{m}
			// 同名のAlbumが存在しない場合は新規作成
			if err := tx.Where("owner_id = ? AND name = ?", userID, album.Name).Attrs(*album).FirstOrCreate(album).Error; err != nil {
				return fmt.Errorf("failed to create album: %w", err)
			}
		}

		return nil
	})
}

func (p *PostgresClient) UpdateMusicStatuses(ctx context.Context, musicIDs []uuid.UUID, status entity.MusicStatus) error {
	result := p.DB(ctx).Model(&entity.Music{}).Where("id IN ?", musicIDs).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update music status: %w", result.Error)
	}
	return nil
}
