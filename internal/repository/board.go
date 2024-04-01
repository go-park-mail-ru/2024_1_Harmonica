package repository

import (
	"context"
	"github.com/jackskj/carta"
	"harmonica/internal/entity"
)

//board
//board_author
//board_pin

const (
	QueryCreateBoard = `INSERT INTO public.board (title, description, cover_url, visibility_type) 
    VALUES ($1, $2, $3, $4) RETURNING public.board.board_id`

	QueryInsertBoardAuthor = `INSERT INTO public.board_author (board_id, author_id) VALUES ($1, $2)`

	QueryGetBoardById = ` SELECT public.user.user_id, public.user.nickname, public.user.avatar_url, 
    public.board.board_id, public.board.title, public.board.created_at, public.board.description, 
    public.board.cover_url, public.board.visibility_type FROM public.board 
    INNER JOIN public.board_author ON public.board.board_id = public.board_author.board_id
    INNER JOIN public.user ON public.board_author.author_id = public.user.user_id 
    WHERE public.board.board_id = $1`

	QueryGetUserBoards = `SELECT board_id, title, created_at, description, cover_url, visibility_type
	FROM public.board INNER JOIN public.user ON public.board_author.author_id=public.user.user_id 
	WHERE public.board_author.author_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
)

var (
	emptyUserBoards = entity.UserBoards{}
	emptyFullBoard  = entity.FullBoard{}
)

func (r *DBRepository) CreateBoard(ctx context.Context, board entity.FullBoard) (entity.BoardID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var boardID entity.BoardID
	err = tx.QueryRowContext(ctx, QueryCreateBoard, board.Board.Title, board.Board.Description,
		board.Board.CoverURL, board.Board.VisibilityType).Scan(&boardID)
	if err != nil {
		return 0, err
	}
	_, err = tx.ExecContext(ctx, QueryInsertBoardAuthor, boardID, board.BoardAuthor.UserId)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return boardID, nil
}

func (r *DBRepository) GetBoardById(ctx context.Context, boardId entity.BoardID) (entity.FullBoard, error) {
	result := entity.FullBoard{}
	err := r.db.QueryRowxContext(ctx, QueryGetBoardById, boardId).StructScan(&result)
	if err != nil {
		return emptyFullBoard, err
	}
	return result, nil
}

func (r *DBRepository) GetUserBoards(ctx context.Context, authorId entity.UserID,
	limit, offset int) (entity.UserBoards, error) {
	result := entity.UserBoards{}
	rows, err := r.db.QueryContext(ctx, QueryGetUserBoards, authorId, limit, offset)
	if err != nil {
		return emptyUserBoards, err
	}
	err = carta.Map(rows, &result.Boards)
	if err != nil {
		return emptyUserBoards, err
	}
	return result, nil
}
