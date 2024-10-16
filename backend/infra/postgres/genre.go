package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gorm.io/gorm/clause"
)

func (p *PostgresClient) CreateGenre(ctx context.Context, g entity.Genre) error {
	result := p.DB(ctx).Create(&g)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateGenre(ctx context.Context, g entity.Genre) error {
	result := p.DB(ctx).Save(&g)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteGenre(ctx context.Context, g entity.Genre) error {
	result := p.DB(ctx).Select(clause.Associations).Delete(&g)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetGenreByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Genre, error) {
	g := make([]entity.Genre, 0, len(ids))
	result := p.DB(ctx).Find(&g, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get genre by ids: %w", result.Error)
	}
	return g, nil
}

func (p *PostgresClient) GetGenreByID(ctx context.Context, id uuid.UUID) (entity.Genre, error) {
	var g entity.Genre
	result := p.DB(ctx).First(&g, id)
	if result.Error != nil {
		return entity.Genre{}, fmt.Errorf("failed to get genre by id: %w", result.Error)
	}
	return g, nil
}

func (p *PostgresClient) GetGenreByName(ctx context.Context, ownerID uuid.UUID, name string) (entity.Genre, error) {
	var g entity.Genre
	result := p.DB(ctx).Where("owner_id = ? AND name = ?", ownerID, name).First(&g)
	if result.Error != nil {
		return entity.Genre{}, fmt.Errorf("failed to get genre by name: %w", result.Error)
	}
	return g, nil
}
