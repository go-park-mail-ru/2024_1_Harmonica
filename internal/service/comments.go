package service

import (
	"context"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) AddComment(ctx context.Context, comment string, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	exists, err := s.repo.CheckPinExistence(ctx, pinId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !exists {
		return errs.ErrorInfo{LocalErr: errs.ErrNotFound}
	}
	if len(comment) <= 0 {
		return errs.ErrorInfo{
			LocalErr: errs.ErrEmptyComment,
		}
	}
	err = s.repo.AddComment(ctx, comment, pinId, userId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return errs.ErrorInfo{}
}

func (s *RepositoryService) GetComments(ctx context.Context, pinId entity.PinID) (entity.GetCommentsResponse, errs.ErrorInfo) {
	exists, err := s.repo.CheckPinExistence(ctx, pinId)
	if err != nil {
		return entity.GetCommentsResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !exists {
		return entity.GetCommentsResponse{}, errs.ErrorInfo{LocalErr: errs.ErrNotFound}
	}
	res, err := s.repo.GetComments(ctx, pinId)
	fmt.Println(err)
	if err != nil {
		return entity.GetCommentsResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return res, errs.ErrorInfo{}
}
