package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
	"gorm.io/gorm/clause"
)

type TestObject struct {
	gormmodel.UUIDModel
	Name string
}

func (p *PostgresClient) CreateTestObject(ctx context.Context, to TestObject) error {
	result := p.DB(ctx).Create(&to)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) UpdateTestObject(ctx context.Context, to TestObject) error {
	result := p.DB(ctx).Save(&to)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) DeleteTestObjects(ctx context.Context, to []TestObject) error {
	result := p.DB(ctx).Select(clause.Associations).Delete(&to)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) GetTestObjectByID(ctx context.Context, id uuid.UUID) (TestObject, error) {
	var to TestObject
	result := p.DB(ctx).First(&to, id)
	if result.Error != nil {
		return TestObject{}, fmt.Errorf("failed to get test object by id: %w", result.Error)
	}
	return to, nil
}

func (p *PostgresClient) GetTestObjectByIDs(ctx context.Context, ids []uuid.UUID) ([]TestObject, error) {
	to := make([]TestObject, 0, len(ids))
	result := p.DB(ctx).Find(&to, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get test object by ids: %w", result.Error)
	}
	return to, nil
}
