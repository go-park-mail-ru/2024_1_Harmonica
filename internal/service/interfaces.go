package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type IService interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, errs.ErrorInfo)
	GetUserByNickname(ctx context.Context, nickname string) (entity.User, errs.ErrorInfo)
	GetUserById(ctx context.Context, id entity.UserID) (entity.User, errs.ErrorInfo)
	RegisterUser(ctx context.Context, user entity.User) []errs.ErrorInfo
	UpdateUser(ctx context.Context, user entity.User) (entity.User, errs.ErrorInfo)

	GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, errs.ErrorInfo)
	GetUserPins(ctx context.Context, authorNickname string, limit, offset int) (entity.UserPins, errs.ErrorInfo)
	GetPinById(ctx context.Context, id entity.PinID) (entity.PinPageResponse, errs.ErrorInfo)
	CreatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo)
	UpdatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo)
	DeletePin(ctx context.Context, pin entity.Pin) errs.ErrorInfo

	SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo
	ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo
	GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, errs.ErrorInfo)

	UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	GetImage(ctx context.Context, name string) (*minio.Object, error)
}
