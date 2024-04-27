package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) GetRating(ctx context.Context) (entity.RatingList, errs.ErrorInfo) {
	ratings, err := s.repo.GetRating(ctx)
	if err != nil {
		return entity.RatingList{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return ratings, emptyErrorInfo
}
