package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) UpdateDraft(ctx context.Context, draft entity.Draft) errs.ErrorInfo {
	err := s.repo.UpdateDraft(ctx, draft)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}
