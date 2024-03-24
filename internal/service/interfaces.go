package service

import (
	"context"
	"harmonica/internal/entity"
)

type IService interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (entity.User, error)
	GetUserById(ctx context.Context, id int64) (entity.User, error)
	RegisterUser(ctx context.Context, user entity.User) []error
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	GetPins(ctx context.Context, limit, offset int) (entity.Pins, error)
}
