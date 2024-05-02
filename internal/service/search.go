package service

import (
	"context"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"strings"
)

func (s *RepositoryService) Search(ctx context.Context, request entity.SearchRequest) (entity.SearchResult, errs.ErrorInfo) {
	request.SearchQuery = strings.Trim(request.SearchQuery, " ")
	var result entity.SearchResult
	users, err := s.repo.SearchForUsers(ctx, fmt.Sprintf(`%s%s%s`, "%", request.SearchQuery, "%"))
	if err != nil {
		return entity.SearchResult{}, errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: err}
	}
	result.Users = append(result.Users, users...)

	pins, err := s.repo.SearchForPins(ctx, request.SearchQuery)
	if err != nil {
		return entity.SearchResult{}, errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: err}
	}
	result.Pins = append(result.Pins, pins...)

	boards, err := s.repo.SearchForBoards(ctx, request.SearchQuery)
	if err != nil {
		return entity.SearchResult{}, errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: err}
	}
	result.Boards = append(result.Boards, boards...)

	return result, errs.ErrorInfo{}
}
