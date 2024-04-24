package service

import (
	"context"
	"database/sql"
	"errors"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"

	"github.com/lib/pq"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

const UniqueViolationErrCode = pq.ErrorCode("23505")

func (s *RepositoryService) CreateBoard(ctx context.Context, board entity.Board,
	userId entity.UserID) (entity.FullBoard, errs.ErrorInfo) {
	board.Sanitize()
	createdBoard, err := s.repo.CreateBoard(ctx, board, userId)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	fullBoard := entity.FullBoard{Board: createdBoard}
	author, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	fullBoard.BoardAuthors = append(fullBoard.BoardAuthors, entity.BoardAuthor{
		UserId:    userId,
		Nickname:  author.Nickname,
		AvatarURL: author.AvatarURL,
	})
	return fullBoard, emptyErrorInfo
}

func (s *RepositoryService) GetBoardById(ctx context.Context, boardId entity.BoardID,
	userId entity.UserID, limit, offset int) (entity.FullBoard, errs.ErrorInfo) {
	board, err := s.repo.GetBoardById(ctx, boardId)
	if err != nil {
		localErr := errs.ErrDBInternal
		if errors.Is(err, sql.ErrNoRows) {
			localErr = errs.ErrElementNotExist
		}
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   localErr,
		}
	}
	isAuthor, err := s.repo.CheckBoardAuthorExistence(ctx, userId, board.BoardID)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	board.IsOwner = isAuthor
	if board.VisibilityType == entity.VisibilityPrivate && !isAuthor {
		return entity.FullBoard{}, errs.ErrorInfo{
			LocalErr: errs.ErrPermissionDenied,
		}
	}
	authors, err := s.repo.GetBoardAuthors(ctx, boardId)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	//if board.VisibilityType == entity.VisibilityPrivate && !AuthorContains(authors, userId) {
	//	return entity.FullBoard{}, errs.ErrorInfo{
	//		LocalErr: errs.ErrPermissionDenied,
	//	}
	//} // по-моему это дубрирует "if board.VisibilityType == entity.VisibilityPrivate && !isAuthor {"

	//if board.VisibilityType == entity.VisibilityPrivate && (userId == 0 || !AuthorContains(authors, userId)) {
	//	return emptyFullBoard, errs.ErrorInfo{
	//		LocalErr: errs.ErrPermissionDenied,
	//	}
	//}

	pins, err := s.repo.GetBoardPins(ctx, boardId, limit, offset)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	fullBoard := entity.FullBoard{
		Board:        board,
		BoardAuthors: authors,
		Pins:         pins,
	}
	return fullBoard, emptyErrorInfo
}

func (s *RepositoryService) UpdateBoard(ctx context.Context, board entity.Board,
	userId entity.UserID) (entity.FullBoard, errs.ErrorInfo) {
	board.Sanitize()
	isAuthor, err := s.repo.CheckBoardAuthorExistence(ctx, userId, board.BoardID)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	// здесь же сразу проверяется существование доски (если доски нет, isAuthor = false)
	if !isAuthor {
		return entity.FullBoard{}, errs.ErrorInfo{
			LocalErr: errs.ErrPermissionDenied,
		}
	}
	updatedBoard, err := s.repo.UpdateBoard(ctx, board)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	authors, err := s.repo.GetBoardAuthors(ctx, board.BoardID)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	pins, err := s.repo.GetBoardPins(ctx, board.BoardID, defaultLimit, defaultOffset)
	if err != nil {
		return entity.FullBoard{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	fullBoard := entity.FullBoard{
		Board:        updatedBoard,
		BoardAuthors: authors,
		Pins:         pins,
	}
	return fullBoard, emptyErrorInfo
}

func (s *RepositoryService) AddPinToBoard(ctx context.Context, boardId entity.BoardID,
	pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	exists, err := s.repo.CheckPinExistence(ctx, pinId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !exists {
		return errs.ErrorInfo{
			LocalErr: errs.ErrElementNotExist,
		}
	}
	isAuthor, err := s.repo.CheckBoardAuthorExistence(ctx, userId, boardId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !isAuthor {
		return errs.ErrorInfo{
			LocalErr: errs.ErrPermissionDenied,
		}
	}
	err = s.repo.AddPinToBoard(ctx, boardId, pinId)
	if err != nil {
		localErr := errs.ErrDBInternal
		var pqErr *pq.Error
		ok := errors.As(err, &pqErr)
		if ok && (pqErr.Code == UniqueViolationErrCode) {
			localErr = errs.ErrDBUniqueViolation
		}
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   localErr,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) DeletePinFromBoard(ctx context.Context, boardId entity.BoardID,
	pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	exists, err := s.repo.CheckPinExistence(ctx, pinId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !exists {
		return errs.ErrorInfo{
			LocalErr: errs.ErrElementNotExist,
		}
	}
	isAuthor, err := s.repo.CheckBoardAuthorExistence(ctx, userId, boardId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !isAuthor {
		return errs.ErrorInfo{
			LocalErr: errs.ErrPermissionDenied,
		}
	}
	err = s.repo.DeletePinFromBoard(ctx, boardId, pinId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) DeleteBoard(ctx context.Context, boardId entity.BoardID,
	userId entity.UserID) errs.ErrorInfo {
	isAuthor, err := s.repo.CheckBoardAuthorExistence(ctx, userId, boardId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !isAuthor {
		return errs.ErrorInfo{
			LocalErr: errs.ErrPermissionDenied,
		}
	}
	err = s.repo.DeleteBoard(ctx, boardId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) GetUserBoards(ctx context.Context, authorNickname string,
	userId entity.UserID, limit, offset int) (entity.UserBoards, errs.ErrorInfo) {
	author, errInfo := s.GetUserByNickname(ctx, authorNickname)
	if errInfo != emptyErrorInfo {
		return entity.UserBoards{}, errInfo
	}
	if author == emptyUser {
		return entity.UserBoards{}, errs.ErrorInfo{
			LocalErr: errs.ErrUserNotExist,
		}
	}
	boards, err := s.repo.GetUserBoards(ctx, author.UserID, limit, offset)
	if err != nil {
		return entity.UserBoards{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if author.UserID != userId {
		var filteredBoards entity.UserBoards
		for _, board := range boards.Boards {
			if board.VisibilityType == entity.VisibilityPublic {
				filteredBoards.Boards = append(filteredBoards.Boards, board)
			}
			//if board.VisibilityType == entity.VisibilityPrivate {
			//	isAuthor, errCheck := s.repo.CheckBoardAuthorExistence(ctx, userId, board.BoardID)
			//	if errCheck != nil {
			//		return entity.UserBoards{}, errs.ErrorInfo{GeneralErr: err, LocalErr:   errs.ErrDBInternal}
			//	}
			//	if isAuthor {
			//		filteredBoards.Boards = append(filteredBoards.Boards, board)
			//	}
			//}
		}
		return filteredBoards, emptyErrorInfo
	}
	return boards, emptyErrorInfo
}

//func AuthorContains(authors []entity.BoardAuthor, userId entity.UserID) bool {
//	for _, author := range authors {
//		if author.UserId == userId {
//			return true
//		}
//	}
//	return false
//}
