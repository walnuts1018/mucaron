package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gorm.io/gorm/clause"
)

func (p *PostgresClient) CreateArtist(a entity.Artist) error {
	result := p.db.Create(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateArtist(a entity.Artist) error {
	result := p.db.Save(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteArtist(a entity.Artist) error {
	result := p.db.Select(clause.Associations).Delete(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetArtistByIDs(ids []uuid.UUID) ([]entity.Artist, error) {
	a := make([]entity.Artist, 0, len(ids))
	result := p.db.Find(&a, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get artist by ids: %w", result.Error)
	}
	return a, nil
}

func (p *PostgresClient) GetArtistByID(id uuid.UUID) (entity.Artist, error) {
	var a entity.Artist
	result := p.db.First(&a, id)
	if result.Error != nil {
		return entity.Artist{}, fmt.Errorf("failed to get artist by id: %w", result.Error)
	}
	return a, nil
}

func (p *PostgresClient) GetArtistByName(ownerID uuid.UUID, name string) (entity.Artist, error) {
	var a entity.Artist
	result := p.db.Where("owner_id = ? AND name = ?", ownerID, name).First(&a)
	if result.Error != nil {
		return entity.Artist{}, fmt.Errorf("failed to get artist by name: %w", result.Error)
	}
	return a, nil
}
