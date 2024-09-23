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

func NewPostgres(cfg config.Config) (*PostgresClient, error) {
	db, err := gorm.Open(postgresdriver.Open(cfg.PSQLDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	c := &PostgresClient{
		db: db,
	}

	if err := c.db.AutoMigrate(&entity.Album{}, &entity.Artist{}, &entity.Genre{}, &entity.Music{}, &entity.Playlist{}, &entity.User{}); err != nil {
		return nil, fmt.Errorf("failed to automigrate: %v", err)
	}

	return c, nil
}

// var _ subjects.SubjectRepository = &PostgresClient{}
