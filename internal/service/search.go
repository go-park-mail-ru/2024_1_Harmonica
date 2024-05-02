package service

import (
	"context"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"strings"
)

func (s *RepositoryService) Search(ctx context.Context, query string) (entity.SearchResult, errs.ErrorInfo) {
	query = strings.Trim(query, " ")
	var result entity.SearchResult
	users, err := s.repo.SearchForUsers(ctx, fmt.Sprintf(`%s%s%s`, "%", query, "%"))
	if err != nil {
		return entity.SearchResult{}, errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: err}
	}
	result.Users = append(result.Users, users...)

	pins, err := s.repo.SearchForPins(ctx, query)
	if err != nil {
		return entity.SearchResult{}, errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: err}
	}
	result.Pins = append(result.Pins, pins...)

	boards, err := s.repo.SearchForBoards(ctx, query)
	if err != nil {
		return entity.SearchResult{}, errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: err}
	}
	result.Boards = append(result.Boards, boards...)

	return result, errs.ErrorInfo{}
}
