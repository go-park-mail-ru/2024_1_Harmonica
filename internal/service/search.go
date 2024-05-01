package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) Search(ctx context.Context, request entity.SearchRequest) (entity.SearchResult, errs.ErrorInfo) {
	return entity.SearchResult{}, errs.ErrorInfo{}
}
