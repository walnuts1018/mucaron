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

type dbController struct {
	db *gorm.DB
}

func newDBController(db *gorm.DB) dbControllerInterface {
	return &dbController{db: db}
}

func (c *dbController) DB(ctx context.Context) *gorm.DB {
	return c.db
}

type dbControllerInterface interface {
	DB(ctx context.Context) *gorm.DB
}

type PostgresClient struct {
	dbControllerInterface
}

var Entities = []any{&entity.Album{}, &entity.Artist{}, &entity.Genre{}, &entity.Music{}, &entity.Playlist{}, &entity.User{}, &entity.RawMusicMetadata{}, &entity.RawMusicMetadataTag{}}

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

	if err := c.DB(ctx).AutoMigrate(Entities...); err != nil {
		return nil, fmt.Errorf("failed to automigrate: %v", err)
	}

	return c, nil
}
