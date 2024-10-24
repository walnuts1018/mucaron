package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
	"gorm.io/gorm/clause"
)

type testObject struct {
	gormmodel.UUIDModel
	Name string
}

func (p *PostgresClient) createTestObject(ctx context.Context, to testObject) error {
	result := p.DB(ctx).Create(&to)
	if result.Error != nil {
		return fmt.Errorf("failed to create: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) updateTestObject(ctx context.Context, to testObject) error {
	result := p.DB(ctx).Save(&to)
	if result.Error != nil {
		return fmt.Errorf("failed to update: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) deleteTestObjects(ctx context.Context, to []testObject) error {
	result := p.DB(ctx).Select(clause.Associations).Delete(&to)
	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (p *PostgresClient) getTestObjectByID(ctx context.Context, id uuid.UUID) (testObject, error) {
	var to testObject
	result := p.DB(ctx).First(&to, id)
	if result.Error != nil {
		return testObject{}, fmt.Errorf("failed to get test object by id: %w", result.Error)
	}
	return to, nil
}

func (p *PostgresClient) getTestObjectByIDs(ctx context.Context, ids []uuid.UUID) ([]testObject, error) {
	to := make([]testObject, 0, len(ids))
	result := p.DB(ctx).Find(&to, ids)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get test object by ids: %w", result.Error)
	}
	return to, nil
}
