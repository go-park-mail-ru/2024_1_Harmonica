package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

type IService interface {
	SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo
	ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo
	CheckIsLiked(ctx context.Context, pinId entity.PinID, userId entity.UserID) (bool, error)
	GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, errs.ErrorInfo)
	GetFavorites(ctx context.Context, userId entity.UserID, limit, offset int) (entity.FeedPins, errs.ErrorInfo)
}
