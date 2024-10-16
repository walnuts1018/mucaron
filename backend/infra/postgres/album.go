package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gorm.io/gorm/clause"
)

func (p *PostgresClient) CreateAlbum(ctx context.Context, a entity.Album) error {
	result := p.DB(ctx).Create(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateAlbum(ctx context.Context, a entity.Album) error {
	result := p.DB(ctx).Save(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteAlbums(ctx context.Context, a []entity.Album) error {
	result := p.DB(ctx).Select(clause.Associations).Delete(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetAlbumByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Album, error) {
	a := make([]entity.Album, 0, len(ids))
	result := p.DB(ctx).Preload(clause.Associations).Find(&a, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get album by ids: %w", result.Error)
	}
	return a, nil
}

func (p *PostgresClient) GetAlbumByID(ctx context.Context, id uuid.UUID) (entity.Album, error) {
	var a entity.Album
	result := p.DB(ctx).Preload(clause.Associations).First(&a, id)
	if result.Error != nil {
		return entity.Album{}, fmt.Errorf("failed to get album by id: %w", result.Error)
	}
	return a, nil
}

func (p *PostgresClient) GetAlbumsByNameAndArtist(ctx context.Context, ownerID uuid.UUID, albumName string, artist entity.Artist) ([]entity.Album, error) {
	albums := make([]entity.Album, 0)
	result := p.DB(ctx).Preload(clause.Associations).Where("owner_id = ? AND name = ? AND id IN (?)", ownerID, albumName, p.DB(ctx).Table("album_musics").Select("album_id").Where("music_id IN (?)", p.DB(ctx).Table("music_artists").Select("music_id").Where("artist_id = ?", artist.ID))).Find(&albums)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get album by name and artist: %w", result.Error)
	}
	return albums, nil
}
