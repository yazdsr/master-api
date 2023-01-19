package postgres

import (
	"context"
	"fmt"

	"github.com/yazdsr/master-api/internal/config"
	"github.com/yazdsr/master-api/internal/pkg/logger"
	"github.com/yazdsr/master-api/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postgres struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

func New(cfg config.Psql, logger logger.Logger) (repository.Postgres, error) {
	db, err := pgxpool.Connect(context.Background(), url(cfg))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &postgres{db: db, logger: logger}, nil
}

func url(cfg config.Psql) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}
