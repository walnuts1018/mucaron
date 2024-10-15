package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"gorm.io/gorm/clause"
)

func (p *PostgresClient) CreateUser(ctx context.Context, u entity.User) error {
	result := p.DB(ctx).Create(&u)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateUser(ctx context.Context, u entity.User) error {
	result := p.DB(ctx).Save(&u)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteUser(ctx context.Context, u entity.User) error {
	result := p.DB(ctx).Select(clause.Associations).Delete(&u)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetUserByIDs(ctx context.Context, userIDs []uuid.UUID) ([]entity.User, error) {
	u := make([]entity.User, 0, len(userIDs))
	result := p.DB(ctx).Find(&u, userIDs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get music by ids: %w", result.Error)
	}
	return u, nil
}

func (p *PostgresClient) GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var u entity.User
	result := p.DB(ctx).First(&u, userID)
	if result.Error != nil {
		return entity.User{}, fmt.Errorf("failed to get user by id: %w", result.Error)
	}
	return u, nil
}

func (p *PostgresClient) GetUserByName(ctx context.Context, userName string) (entity.User, error) {
	var u entity.User
	result := p.DB(ctx).Where("user_name = ?", userName).First(&u)
	if result.Error != nil {
		return entity.User{}, fmt.Errorf("failed to get user by name: %w", result.Error)
	}
	return u, nil
}
