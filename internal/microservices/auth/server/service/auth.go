package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) GetUserByEmail(ctx context.Context, email string) (entity.User, errs.ErrorInfo) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return entity.User{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return user, errs.ErrorInfo{}
}

func (s *RepositoryService) GetUserById(ctx context.Context, id entity.UserID) (entity.User, errs.ErrorInfo) {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return entity.User{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return user, errs.ErrorInfo{}
}
