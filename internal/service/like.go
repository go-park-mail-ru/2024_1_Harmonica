package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

var emptyErrorInfo = errs.ErrorInfo{}

func (r *RepositoryService) SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	err := r.repo.SetLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (r *RepositoryService) ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	err := r.repo.ClearLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (r *RepositoryService) GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, errs.ErrorInfo) {
	res, err := r.repo.GetUsersLiked(ctx, pinId, limit)
	if err != nil {
		return entity.UserList{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return res, emptyErrorInfo
}
