package repository

import (
	"context"
	"harmonica/internal/entity"
)

type IRepository interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserById(ctx context.Context, id entity.UserID) (entity.User, error)
}
