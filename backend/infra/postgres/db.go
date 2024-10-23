package postgres

import (
	"context"
	"fmt"

	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type MucaronTransactionKey struct{}

type dbController struct {
	db *gorm.DB
}

func newDBController(db *gorm.DB) dbControllerInterface {
	return &dbController{db: db}
}

func (c *dbController) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(MucaronTransactionKey{}).(*gorm.DB); ok {
		return tx
	}

	return c.db.WithContext(ctx)
}

func (c *dbController) Transaction(ctx context.Context, f func(ctx context.Context) error) error {
	tx := c.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	ctx = context.WithValue(ctx, MucaronTransactionKey{}, tx)
	if err := f(ctx); err != nil {
		if rerr := tx.Rollback().Error; rerr != nil {
			return fmt.Errorf("failed to rollback transaction: %w, original error: %v", rerr, err)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type dbControllerInterface interface {
	DB(ctx context.Context) *gorm.DB
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}

type PostgresClient struct {
	dbControllerInterface
}

var entities = []any{&entity.Album{}, &entity.Artist{}, &entity.Genre{}, &entity.Music{}, &entity.Playlist{}, &entity.User{}, &entity.RawMusicMetadata{}, &entity.RawMusicMetadataTag{}}

func NewPostgres(ctx context.Context, cfg config.Config) (*PostgresClient, error) {
	db, err := gorm.Open(postgresdriver.Open(cfg.PSQLDSN), &gorm.Config{
		Logger: NewLogger(cfg),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, fmt.Errorf("failed to use tracing plugin: %v", err)
	}

	c := &PostgresClient{
		newDBController(db),
	}

	if err := c.DB(ctx).AutoMigrate(entities...); err != nil {
		return nil, fmt.Errorf("failed to automigrate: %v", err)
	}

	return c, nil
}
