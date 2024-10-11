package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (p *PostgresClient) CreateAlbum(a entity.Album) error {
	result := p.db.Create(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateAlbum(a entity.Album) error {
	result := p.db.Save(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteAlbum(a entity.Album) error {
	result := p.db.Delete(&a)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetAlbumByIDs(ids []string) ([]entity.Album, error) {
	a := make([]entity.Album, 0, len(ids))
	result := p.db.Find(&a, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get album by ids: %w", result.Error)
	}
	return a, nil
}

func (p *PostgresClient) GetAlbumByID(id string) (entity.Album, error) {
	var a entity.Album
	result := p.db.First(&a, id)
	if result.Error != nil {
		return entity.Album{}, fmt.Errorf("failed to get album by id: %w", result.Error)
	}
	return a, nil
}

func (p *PostgresClient) GetAlbumByName(owner uuid.UUID, albumName string) (entity.Album, error) {
	var a entity.Album
	result := p.db.Where("name = ?", name).First(&a)
	if result.Error != nil {
		return entity.Album{}, fmt.Errorf("failed to get album by name: %w", result.Error)
	}
	return a, nil
}
