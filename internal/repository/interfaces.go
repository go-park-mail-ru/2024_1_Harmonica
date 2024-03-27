package repository

import (
	"context"
	"harmonica/internal/entity"
)

type IRepository interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (entity.User, error)
	GetUserById(ctx context.Context, id entity.UserID) (entity.User, error)
	RegisterUser(ctx context.Context, user entity.User) error
	UpdateUser(ctx context.Context, user entity.User) error

	GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, error)
	GetUserPins(ctx context.Context, authorId entity.UserID, limit, offset int) (entity.UserPins, error)
	GetPinById(ctx context.Context, id entity.PinID) (entity.PinPageResponse, error)
	CreatePin(ctx context.Context, pin entity.Pin) (entity.PinID, error)
	UpdatePin(ctx context.Context, pin entity.Pin) error
	DeletePin(ctx context.Context, id entity.PinID) error

	CreateLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error
	DeleteLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error
	GetLikedPins(ctx context.Context, userId entity.UserID, limit int) (entity.FeedPins, error)
	GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, error)
	GetLikesCount(ctx context.Context, pinId entity.PinID) (uint64, error)
	FindLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) (bool, error)
}
