package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

type IService interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, errs.ErrorInfo)
	GetUserById(ctx context.Context, id entity.UserID) (entity.User, errs.ErrorInfo)
}
