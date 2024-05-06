package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	image "harmonica/internal/microservices/image/proto"
)

type DBRepository struct {
	db           *sqlx.DB
	imageService image.ImageClient
	logger       *zap.Logger
}

func NewDBRepository(db *sqlx.DB, l *zap.Logger, imCli image.ImageClient) *DBRepository {
	return &DBRepository{db: db, logger: l, imageService: imCli}
}

type Repository struct {
	IRepository
}

func NewRepository(c *Connector, l *zap.Logger, imCli image.ImageClient) *Repository {
	return &Repository{
		IRepository: NewDBRepository(c.db, l, imCli),
	}
}
