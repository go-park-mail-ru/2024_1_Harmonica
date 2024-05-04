package repository

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/microservices/image/proto"
	"strings"
	"time"

	"github.com/jackskj/carta"
)

const (
	QueryCreateBoard = `INSERT INTO public.board (title, description, visibility_type) 
    VALUES ($1, $2, $3) RETURNING public.board.board_id, public.board.created_at, public.board.title, 
    public.board.description, public.board.cover_url, public.board.visibility_type`

	QueryInsertBoardAuthor = `INSERT INTO public.board_author (board_id, author_id) VALUES ($1, $2)`

	QueryGetBoardById = `SELECT board_id, title, created_at, description, cover_url, visibility_type
	FROM public.board WHERE board_id=$1`

	QueryGetBoardAuthors = `SELECT public.user.user_id, public.user.nickname, public.user.avatar_url FROM public.user
	INNER JOIN public.board_author ON public.user.user_id = public.board_author.author_id 
	WHERE public.board_author.board_id=$1`

	QueryGetBoardPins = `SELECT public.pin.pin_id, public.pin.content_url, public.user.user_id, public.user.nickname, 
    public.user.avatar_url FROM public.pin INNER JOIN public.board_pin ON public.pin.pin_id = public.board_pin.pin_id 
	INNER JOIN public.user ON public.pin.author_id = public.user.user_id WHERE public.board_pin.board_id=$1
	ORDER BY public.pin.created_at DESC LIMIT $2 OFFSET $3`

	QueryUpdateBoard = `UPDATE public.board SET title=$2, description=$3, cover_url=$4, visibility_type=$5 
    WHERE board_id=$1 RETURNING public.board.board_id, public.board.created_at, public.board.title, 
    public.board.description, public.board.cover_url, public.board.visibility_type`

	QueryAddPinToBoard = `INSERT INTO public.board_pin (board_id, pin_id) VALUES ($1, $2)`

	QueryDeletePinFromBoard = `DELETE FROM public.board_pin WHERE board_id=$1 AND pin_id=$2`

	QueryDeleteBoard = `DELETE FROM public.board WHERE board_id=$1`

	QueryGetUserBoards = `SELECT b.board_id, b.title, b.created_at, b.description, b.cover_url, b.visibility_type,
    ARRAY (SELECT p.content_url FROM public.pin p INNER JOIN public.board_pin bp ON p.pin_id = bp.pin_id
	WHERE bp.board_id = b.board_id ORDER BY bp.added_at DESC LIMIT 3) AS recent_pins
	FROM public.board b INNER JOIN public.board_author ba ON b.board_id = ba.board_id 
	WHERE ba.author_id = $1 AND ($1 = $2 OR b.visibility_type = 'public')
	ORDER BY b.created_at DESC LIMIT $3 OFFSET $4;`

	//OldQueryGetUserBoards = `SELECT public.board.board_id, public.board.title, public.board.created_at,
	//public.board.description, public.board.cover_url, public.board.visibility_type FROM public.board
	//INNER JOIN public.board_author ON public.board.board_id = public.board_author.board_id
	//WHERE public.board_author.author_id=$1 ORDER BY public.board.created_at DESC LIMIT $2 OFFSET $3`

	QueryCheckBoardAuthorExistence = `SELECT EXISTS(SELECT 1 FROM public.board_author WHERE author_id=$1 AND board_id=$2)`
)

func (r *DBRepository) CreateBoard(ctx context.Context, board entity.Board,
	userId entity.UserID) (entity.Board, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return entity.Board{}, err
	}
	defer tx.Rollback()

	var createdBoard entity.Board
	start := time.Now()
	err = tx.QueryRowxContext(ctx, QueryCreateBoard, board.Title, board.Description,
		board.VisibilityType).StructScan(&createdBoard)
	LogDBQuery(r, ctx, QueryCreateBoard, time.Since(start))
	if err != nil {
		return entity.Board{}, err
	}

	res, err := r.ImageService.GetImageBounds(ctx, &proto.GetImageBoundsRequest{Url: createdBoard.CoverURL})
	if err != nil {
		return entity.Board{}, err
	}
	createdBoard.CoverDX = res.Dx
	createdBoard.CoverDY = res.Dy

	start = time.Now()
	_, err = tx.ExecContext(ctx, QueryInsertBoardAuthor, createdBoard.BoardID, userId)
	LogDBQuery(r, ctx, QueryInsertBoardAuthor, time.Since(start))
	if err != nil {
		return entity.Board{}, err
	}

	if err = tx.Commit(); err != nil {
		return entity.Board{}, err
	}
	return createdBoard, nil
}

func (r *DBRepository) GetBoardById(ctx context.Context, boardId entity.BoardID) (entity.Board, error) {
	board := entity.Board{}
	start := time.Now()
	err := r.db.QueryRowxContext(ctx, QueryGetBoardById, boardId).StructScan(&board)
	LogDBQuery(r, ctx, QueryGetBoardById, time.Since(start))
	if err != nil {
		return entity.Board{}, err
	}

	res, err := r.ImageService.GetImageBounds(ctx, &proto.GetImageBoundsRequest{Url: board.CoverURL})
	if err != nil {
		return entity.Board{}, err
	}
	board.CoverDX = res.Dx
	board.CoverDY = res.Dy

	return board, nil
}

func (r *DBRepository) GetBoardAuthors(ctx context.Context, boardId entity.BoardID) ([]entity.BoardAuthor, error) {
	var authors []entity.BoardAuthor
	start := time.Now()
	err := r.db.SelectContext(ctx, &authors, QueryGetBoardAuthors, boardId)
	LogDBQuery(r, ctx, QueryGetBoardAuthors, time.Since(start))
	if err != nil {
		return []entity.BoardAuthor{}, err
	}
	for i, author := range authors {
		res, err := r.ImageService.GetImageBounds(ctx, &proto.GetImageBoundsRequest{Url: author.AvatarURL})
		if err != nil {
			return []entity.BoardAuthor{}, err
		}
		author.AvatarDX = res.Dx
		author.AvatarDY = res.Dy
		authors[i] = author
	}
	return authors, nil
}

func (r *DBRepository) GetBoardPins(ctx context.Context, boardId entity.BoardID, limit, offset int) ([]entity.BoardPinResponse, error) {
	var pins []entity.BoardPinResponse
	start := time.Now()
	err := r.db.SelectContext(ctx, &pins, QueryGetBoardPins, boardId, limit, offset)
	LogDBQuery(r, ctx, QueryGetBoardPins, time.Since(start))
	if err != nil {
		return []entity.BoardPinResponse{}, err
	}
	for i, pin := range pins {
		res, err := r.ImageService.GetImageBounds(ctx, &proto.GetImageBoundsRequest{Url: pin.ContentUrl})
		if err != nil {
			return []entity.BoardPinResponse{}, err
		}
		pin.ContentDX = res.Dx
		pin.ContentDY = res.Dy
		pins[i] = pin
	}
	return pins, nil
}

func (r *DBRepository) UpdateBoard(ctx context.Context, board entity.Board) (entity.Board, error) {
	var updatedBoard entity.Board
	start := time.Now()
	err := r.db.QueryRowxContext(ctx, QueryUpdateBoard, board.BoardID, board.Title, board.Description,
		board.CoverURL, board.VisibilityType).StructScan(&updatedBoard)
	if err != nil {
		return entity.Board{}, err
	}
	LogDBQuery(r, ctx, QueryUpdateBoard, time.Since(start))

	res, err := r.ImageService.GetImageBounds(ctx, &proto.GetImageBoundsRequest{Url: board.CoverURL})
	if err != nil {
		return entity.Board{}, err
	}
	board.CoverDX = res.Dx
	board.CoverDY = res.Dy
	return updatedBoard, err
}

func (r *DBRepository) AddPinToBoard(ctx context.Context, boardId entity.BoardID, pinId entity.PinID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryAddPinToBoard, boardId, pinId)
	LogDBQuery(r, ctx, QueryAddPinToBoard, time.Since(start))
	return err
}

func (r *DBRepository) DeletePinFromBoard(ctx context.Context, boardId entity.BoardID, pinId entity.PinID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryDeletePinFromBoard, boardId, pinId)
	LogDBQuery(r, ctx, QueryDeletePinFromBoard, time.Since(start))
	return err
}

func (r *DBRepository) DeleteBoard(ctx context.Context, boardId entity.BoardID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryDeleteBoard, boardId)
	LogDBQuery(r, ctx, QueryDeleteBoard, time.Since(start))
	return err
}

func (r *DBRepository) GetUserBoards(ctx context.Context, authorId, userId entity.UserID,
	limit, offset int) (entity.UserBoards, error) {
	boards := entity.UserBoards{}
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetUserBoards, authorId, userId, limit, offset)
	LogDBQuery(r, ctx, QueryGetUserBoards, time.Since(start))
	if err != nil {
		return entity.UserBoards{}, err
	}
	err = carta.Map(rows, &boards.Boards)
	if err != nil {
		return entity.UserBoards{}, err
	}
	for i := range boards.Boards {
		if boards.Boards[i].RecentPinContentUrls[0] == "{}" {
			boards.Boards[i].RecentPinContentUrls = nil
			continue
		}
		recentPins := boards.Boards[i].RecentPinContentUrls
		recentPins[0] = recentPins[0][1 : len(recentPins[0])-1] // убираем фигурные скобки {}
		recentPins = strings.Split(recentPins[0], ",")
		boards.Boards[i].RecentPinContentUrls = recentPins
	}

	return boards, nil
}

func (r *DBRepository) CheckBoardAuthorExistence(ctx context.Context, userId entity.UserID,
	boardId entity.BoardID) (bool, error) {
	var exists bool
	start := time.Now()
	err := r.db.QueryRowContext(ctx, QueryCheckBoardAuthorExistence, userId, boardId).Scan(&exists)
	LogDBQuery(r, ctx, QueryCheckBoardAuthorExistence, time.Since(start))
	return exists, err
}
