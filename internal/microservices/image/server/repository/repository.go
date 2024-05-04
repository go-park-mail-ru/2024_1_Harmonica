package repository

import (
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type MinioRepository struct {
	s3     *minio.Client
	logger *zap.Logger
}

func NewMinioRepository(s3 *minio.Client, l *zap.Logger) *MinioRepository {
	return &MinioRepository{s3: s3, logger: l}
}

type Repository struct {
	IRepository
}

func NewRepository(c *Connector, l *zap.Logger) *Repository {
	return &Repository{
		IRepository: NewMinioRepository(c.s3, l),
	}
}
