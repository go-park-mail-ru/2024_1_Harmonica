package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type DBRepository struct {
	db     *sqlx.DB
	s3     *minio.Client
	logger *zap.Logger
}

func NewDBRepository(db *sqlx.DB, s3 *minio.Client, l *zap.Logger) *DBRepository {
	return &DBRepository{db: db, s3: s3, logger: l}
}

type Repository struct {
	IRepository
}

func NewRepository(c *Connector, l *zap.Logger) *Repository {
	return &Repository{
		IRepository: NewDBRepository(c.db, c.s3, l),
	}
}
