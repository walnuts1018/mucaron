package postgres

import (
	"fmt"

	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (p *PostgresClient) CreateMusic(m entity.Music) error {
	result := p.db.Create(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to insert: %w", result.Error)
	}
	return nil
}
