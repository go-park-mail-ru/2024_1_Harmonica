package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

var emptyErrorInfo = errs.ErrorInfo{}

func (s *RepositoryService) SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	err := s.repo.SetLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	err := s.repo.ClearLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, errs.ErrorInfo) {
	res, err := s.repo.GetUsersLiked(ctx, pinId, limit)
	if err != nil {
		return entity.UserList{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return res, emptyErrorInfo
}

func (s *RepositoryService) GetFavorites(ctx context.Context, limit, offset int) (entity.FeedPins, errs.ErrorInfo) {
	res, err := s.repo.GetFavorites(ctx, ctx.Value("user_id").(entity.UserID), limit, offset)
	if err != nil {
		return entity.FeedPins{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return res, errs.ErrorInfo{}
}
