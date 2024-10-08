package postgres

import (
	"fmt"

	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (p *PostgresClient) CreateMusic(m entity.Music) error {
	result := p.db.Create(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UploadMusic(m entity.Music) error {
	result := p.db.Save(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteMusic(m entity.Music) error {
	result := p.db.Delete(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetMusicByIDs(ids []string) ([]entity.Music, error) {
	m := make([]entity.Music, 0, len(ids))
	result := p.db.Find(&m, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music by ids: %w", result.Error)
	}
	return m, nil
}
