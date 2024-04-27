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

func (s *RepositoryService) CreateRating(ctx context.Context, rating entity.Rating) errs.ErrorInfo {
	err := s.repo.RatingCreate(ctx, rating)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err, // добавила эту ошибку, ранее возвращалось nil
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}
