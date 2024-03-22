package repository

import (
	"github.com/jmoiron/sqlx"
	"harmonica/internal/entity"
)

type IRepository interface {
	GetUserByEmail(email string) (entity.User, error)
	GetUserByNickname(nickname string) (entity.User, error)
	GetUserById(id int64) (entity.User, error)
	RegisterUser(user entity.User) error
	GetPins(limit, offset int) (entity.Pins, error)
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
