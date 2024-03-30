package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DBRepository struct {
	db *sqlx.DB
}

func NewDBRepository(db *sqlx.DB) *DBRepository {
	return &DBRepository{db: db}
}

type Repository struct {
	IRepository
	Logger *zap.Logger
}

func NewRepository(c *Connector, logger *zap.Logger) *Repository {
	return &Repository{
		IRepository: NewDBRepository(c.db),
		Logger:      logger,
	}
}
