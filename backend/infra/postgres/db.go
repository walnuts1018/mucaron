package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/walnuts1018/mucaron/config"
)

type PostgresClient struct {
	db *sqlx.DB
}

func NewPostgres(cfg config.Config) (*PostgresClient, error) {
	db, err := sqlx.Open("postgres", cfg.PSQLDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	c := &PostgresClient{
		db: db,
	}

	return c, nil
}

func (s *PostgresClient) Close() error {
	return s.db.Close()
}

// var _ subjects.SubjectRepository = &PostgresClient{}
