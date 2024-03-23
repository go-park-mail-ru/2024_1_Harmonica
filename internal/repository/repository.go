package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"harmonica/internal/entity"
)

type IRepository interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (entity.User, error)
	GetUserById(ctx context.Context, id int64) (entity.User, error)
	RegisterUser(ctx context.Context, user entity.User) error
	UpdateUser(ctx context.Context, user entity.User) error
	GetPins(ctx context.Context, limit, offset int) (entity.Pins, error)
}

type DBRepository struct {
	db *sqlx.DB
}

func NewDBRepository(db *sqlx.DB) *DBRepository {
	return &DBRepository{db: db}
}

type Repository struct {
	IRepository
}

func NewRepository(c *Connector) *Repository {
	return &Repository{
		IRepository: NewDBRepository(c.db),
	}
}
