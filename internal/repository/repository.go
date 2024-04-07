package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type DBRepository struct {
	db *sqlx.DB
	s3 *minio.Client
}

func NewDBRepository(db *sqlx.DB, s3 *minio.Client) *DBRepository {
	return &DBRepository{db: db, s3: s3}
}

type Repository struct {
	IRepository
}

func NewRepository(c *Connector) *Repository {
	return &Repository{
		IRepository: NewDBRepository(c.db, c.s3),
	}
}
