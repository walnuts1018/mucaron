package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gorm.io/gorm/clause"
)

func (p *PostgresClient) CreateArtist(ctx context.Context, a entity.Artist) error {
	result := p.DB(ctx).Create(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateArtist(ctx context.Context, a entity.Artist) error {
	result := p.DB(ctx).Save(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteArtist(ctx context.Context, a entity.Artist) error {
	result := p.DB(ctx).Select(clause.Associations).Delete(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetArtistByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Artist, error) {
	a := make([]entity.Artist, 0, len(ids))
	result := p.DB(ctx).Find(&a, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get artist by ids: %w", result.Error)
	}
	return a, nil
}

func (p *PostgresClient) GetArtistByID(ctx context.Context, id uuid.UUID) (entity.Artist, error) {
	var a entity.Artist
	result := p.DB(ctx).First(&a, id)
	if result.Error != nil {
		return entity.Artist{}, fmt.Errorf("failed to get artist by id: %w", result.Error)
	}
	return a, nil
}

func (p *PostgresClient) GetArtistByName(ctx context.Context, ownerID uuid.UUID, name string) (entity.Artist, error) {
	var a entity.Artist
	result := p.DB(ctx).Where("owner_id = ? AND name = ?", ownerID, name).First(&a)
	if result.Error != nil {
		return entity.Artist{}, fmt.Errorf("failed to get artist by name: %w", result.Error)
	}
	return a, nil
}
