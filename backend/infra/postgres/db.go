package postgres

import (
	"fmt"

	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresClient struct {
	db *gorm.DB
}

var Entities = []any{&entity.Album{}, &entity.Artist{}, &entity.Genre{}, &entity.Music{}, &entity.Playlist{}, &entity.User{}, &entity.RawMusicMetadata{}, &entity.RawMusicMetadataTag{}}

func NewPostgres(cfg config.Config) (*PostgresClient, error) {
	db, err := gorm.Open(postgresdriver.Open(cfg.PSQLDSN), &gorm.Config{
		Logger: NewLogger(cfg),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	c := &PostgresClient{
		db: db,
	}

	if err := c.db.AutoMigrate(Entities...); err != nil {
		return nil, fmt.Errorf("failed to automigrate: %v", err)
	}

	return c, nil
}
