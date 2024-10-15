package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gorm.io/gorm/clause"
)

func (p *PostgresClient) CreateUser(u entity.User) error {
	result := p.db.Create(&u)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateUser(u entity.User) error {
	result := p.db.Save(&u)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteUser(u entity.User) error {
	result := p.db.Select(clause.Associations).Delete(&u)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetUserByIDs(userIDs []uuid.UUID) ([]entity.User, error) {
	u := make([]entity.User, 0, len(userIDs))
	result := p.db.Find(&u, userIDs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music by ids: %w", result.Error)
	}
	return u, nil
}

func (p *PostgresClient) GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var u entity.User
	result := p.db.WithContext(ctx).First(&u, userID)
	if result.Error != nil {
		return entity.User{}, fmt.Errorf("failed to get user by id: %w", result.Error)
	}
	return u, nil
}

func (p *PostgresClient) GetUserByName(userName string) (entity.User, error) {
	var u entity.User
	result := p.db.Where("user_name = ?", userName).First(&u)
	if result.Error != nil {
		return entity.User{}, fmt.Errorf("failed to get user by name: %w", result.Error)
	}
	return u, nil
}
