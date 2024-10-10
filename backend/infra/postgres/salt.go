package postgres

import (
	"fmt"

	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (p *PostgresClient) CreateHashSalt(salt entity.HashSalt) error {
	result := p.db.Create(&salt)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetHashSalt()
