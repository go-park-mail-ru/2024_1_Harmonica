package repository

import (
	image "harmonica/internal/microservices/image/proto"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DBRepository struct {
	db           *sqlx.DB
	ImageService image.ImageClient
	logger       *zap.Logger
}

func NewDBRepository(db *sqlx.DB, s3 image.ImageClient, l *zap.Logger) *DBRepository {
	return &DBRepository{db: db, ImageService: s3, logger: l}
}

type Repository struct {
	IRepository
}

func NewRepository(c *Connector, l *zap.Logger) *Repository {
	return &Repository{
		IRepository: NewDBRepository(c.db, c.s3, l),
	}
}
