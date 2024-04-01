package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

var emptyFullBoard = entity.FullBoard{}

func (s *RepositoryService) CreateBoard(ctx context.Context, board entity.FullBoard) (entity.FullBoard, errs.ErrorInfo) {
	boardId, err := s.repo.CreateBoard(ctx, board)

	if err != nil {
		return emptyFullBoard, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	res, errInfo := s.GetBoardById(ctx, boardId)

	if err != nil {
		return emptyFullBoard, errInfo
	}
	return res, emptyErrorInfo
}

func (s *RepositoryService) GetBoardById(ctx context.Context, boardId entity.BoardID) (entity.FullBoard, errs.ErrorInfo) {
	board, err := s.repo.GetBoardById(ctx, boardId)
	if err != nil {
		return emptyFullBoard, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return board, emptyErrorInfo
}

func (s *RepositoryService) GetUserBoards(ctx context.Context, authorNickname string,
	limit, offset int) (entity.UserBoards, errs.ErrorInfo) {
	user, errInfo := s.GetUserByNickname(ctx, authorNickname)
	if errInfo.GeneralErr != nil {
		return entity.UserBoards{}, errInfo
	}
	pins, err := s.repo.GetUserBoards(ctx, user.UserID, limit, offset)
	if err != nil {
		return entity.UserBoards{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return pins, emptyErrorInfo
}
