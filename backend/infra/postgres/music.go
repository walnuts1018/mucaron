package postgres

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *PostgresClient) CreateMusic(m entity.Music) error {
	result := p.db.Create(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateMusic(m entity.Music) error {
	result := p.db.Save(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateMusicStatus(musicID uuid.UUID, status entity.MusicStatus) error {
	result := p.db.Model(&entity.Music{}).Where("id = ?", musicID).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update music status: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteMusics(musicIDs []uuid.UUID) error {
	var deleted []entity.Music
	result := p.db.Select(clause.Associations).Select("RawMetaData").Delete(&deleted, musicIDs)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	slog.Debug("deleted musics", slog.Any("deleted", deleted))
	return nil
}

func (p *PostgresClient) GetMusicByID(id uuid.UUID) (entity.Music, error) {
	var m entity.Music
	result := p.db.Preload(clause.Associations).First(&m, id)
	if result.Error != nil {
		return entity.Music{}, fmt.Errorf("failed to get music by id: %w", result.Error)
	}
	return m, nil
}

func (p *PostgresClient) GetMusicByIDs(ids []uuid.UUID) ([]entity.Music, error) {
	m := make([]entity.Music, 0, len(ids))
	result := p.db.Find(&m, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music by ids: %w", result.Error)
	}
	return m, nil
}

func (p *PostgresClient) CreateMusicWithDependencies(m entity.Music, album *entity.Album, artist *entity.Artist, genre *entity.Genre) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if artist != nil {
			// 同名のArtistが存在しない場合は新規作成
			if err := tx.Where("owner_id = ? AND name = ?", artist.OwnerID, artist.Name).Attrs(*artist).FirstOrCreate(artist).Error; err != nil {
				return fmt.Errorf("failed to create artist: %w", err)
			}
			m.Artists = []entity.Artist{*artist}
		}

		if genre != nil {
			// 同名のGenreが存在しない場合は新規作成
			if err := tx.Where("owner_id = ? AND name = ?", genre.OwnerID, genre.Name).Attrs(*genre).FirstOrCreate(genre).Error; err != nil {
				return fmt.Errorf("failed to create genre: %w", err)
			}
			m.GenreID = &genre.ID
		}

		var existingMusic entity.Music
		if err := tx.Unscoped().Preload(clause.Associations).Where("owner_id = ? AND file_hash = ?", m.OwnerID, m.FileHash).First(&existingMusic).Error; err != nil {
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
			if err := tx.Where("owner_id = ? AND name = ?", album.OwnerID, album.Name).Attrs(*album).FirstOrCreate(album).Error; err != nil {
				return fmt.Errorf("failed to create album: %w", err)
			}
		}

		return nil
	})
}

func (p *PostgresClient) UpdateMusicStatuses(musicIDs []uuid.UUID, status entity.MusicStatus) error {
	result := p.db.Model(&entity.Music{}).Where("id IN ?", musicIDs).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update music status: %w", result.Error)
	}
	return nil
}
