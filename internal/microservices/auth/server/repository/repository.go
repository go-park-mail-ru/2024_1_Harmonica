package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DBRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewDBRepository(db *sqlx.DB, l *zap.Logger) *DBRepository {
	return &DBRepository{db: db, logger: l}
}

type Repository struct {
	IRepository
}

func NewRepository(c *Connector, l *zap.Logger) *Repository {
	return &Repository{
		IRepository: NewDBRepository(c.db, l),
	}
}
