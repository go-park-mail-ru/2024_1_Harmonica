package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

var emptyFullBoard = entity.FullBoard{}

func (s *RepositoryService) CreateBoard(ctx context.Context, board entity.Board,
	userId entity.UserID) (entity.FullBoard, errs.ErrorInfo) {

	boardId, err := s.repo.CreateBoard(ctx, board, userId)
	if err != nil {
		return emptyFullBoard, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	res, errInfo := s.GetBoardById(ctx, boardId, userId)
	if err != nil {
		return emptyFullBoard, errInfo
	}
	return res, emptyErrorInfo
}

func (s *RepositoryService) GetBoardById(ctx context.Context, boardId entity.BoardID,
	userId entity.UserID) (entity.FullBoard, errs.ErrorInfo) {

	board, err := s.repo.GetBoardById(ctx, boardId)
	if err != nil {
		return emptyFullBoard, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	authors, err := s.repo.GetBoardAuthors(ctx, boardId)
	if err != nil {
		return emptyFullBoard, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	if board.VisibilityType == "private" && !authorContains(authors, userId) {
		return emptyFullBoard, errs.ErrorInfo{
			LocalErr: errs.ErrPermissionDenied,
		}
	}

	pins, err := s.repo.GetBoardPins(ctx, boardId)
	if err != nil {
		return emptyFullBoard, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	var fullBoard entity.FullBoard
	fullBoard.Board = board
	fullBoard.BoardAuthors = authors
	fullBoard.Pins = pins

	return fullBoard, emptyErrorInfo
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

func authorContains(authors []entity.BoardAuthor, userId entity.UserID) bool {
	for _, author := range authors {
		if author.UserId == userId {
			return true
		}
	}
	return false
}
