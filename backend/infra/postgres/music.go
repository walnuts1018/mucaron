package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
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

func (p *PostgresClient) HardDeleteMusic(ctx context.Context, music entity.Music) error {
	result := p.DB(ctx).Exec("DELETE FROM album_musics WHERE music_id = ?", music.ID).Exec("DELETE FROM raw_music_metadata_tags WHERE raw_music_metadata_id = ?", music.RawMetaData.ID).Unscoped().Select(clause.Associations).Delete(&music)
	if result.Error != nil {
		return fmt.Errorf("failed to hard delete: %w", result.Error)
	}
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
	result := p.DB(ctx).Preload("Artists").Preload("Genre").Where("owner_id = ?", userID).Find(&m)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music by user id: %w", result.Error)
	}
	return m, nil
}

func (p *PostgresClient) GetMusicByFileHash(ctx context.Context, userID uuid.UUID, fileHash string, m *entity.Music) error {
	result := p.DB(ctx).Unscoped().Preload(clause.Associations).Where("owner_id = ? AND file_hash = ?", userID, fileHash).First(m)
	if result.Error != nil {
		return fmt.Errorf("failed to get music by filehash: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetMusicIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	result := p.DB(ctx).Model(&entity.Music{}).Where("owner_id = ?", userID).Pluck("id", &ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music ids by user id: %w", result.Error)
	}
	return ids, nil
}

func (p *PostgresClient) UpdateMusicStatuses(ctx context.Context, musicIDs []uuid.UUID, status entity.MusicStatus) error {
	result := p.DB(ctx).Model(&entity.Music{}).Where("id IN ?", musicIDs).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update music status: %w", result.Error)
	}
	return nil
}
