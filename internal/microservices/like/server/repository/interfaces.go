package repository

import (
	"context"
	"harmonica/internal/entity"
)

type IRepository interface {
	SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error
	ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error
	GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, error)
	CheckIsLiked(ctx context.Context, pinId entity.PinID, userId entity.UserID) (bool, error)
	GetFavorites(ctx context.Context, userId entity.UserID, limit, offset int) (entity.FeedPins, error)
}
