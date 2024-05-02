package repository

import (
	"context"
	"harmonica/internal/entity"
	"time"

	"github.com/jackskj/carta"
)

const (
	QuerySearchForUser = `SELECT user_id, nickname, avatar_url FROM public."user" u WHERE LOWER(u.nickname) LIKE LOWER($1)
	ORDER BY u.register_at DESC LIMIT $2`

	QuerySearchForPin = `SELECT pin_id, title, content_url FROM public.pin p WHERE TO_TSVECTOR(p.title) @@ PLAINTO_TSQUERY($1) OR
	TO_TSVECTOR(p.description) @@ PLAINTO_TSQUERY($1) ORDER BY p.created_at DESC LIMIT $2`

	QuerySearchForBoard = `SELECT board_id, title, cover_url FROM public.board b WHERE b.visibility_type='public'
	AND (TO_TSVECTOR(b.title) @@ PLAINTO_TSQUERY($1) OR TO_TSVECTOR(b.description) @@ PLAINTO_TSQUERY($1))
	ORDER BY b.created_at DESC LIMIT $2`
)

const SEARCH_LIMIT = 10

func (r *DBRepository) SearchForUsers(ctx context.Context, query string) ([]entity.SearchUser, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QuerySearchForUser, query, SEARCH_LIMIT)
	LogDBQuery(r, ctx, QuerySearchForUser, time.Since(start))
	if err != nil {
		return []entity.SearchUser{}, err
	}
	var res []entity.SearchUser
	err = carta.Map(rows, &res)
	return res, err
}

func (r *DBRepository) SearchForPins(ctx context.Context, query string) ([]entity.SearchPin, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QuerySearchForPin, query, SEARCH_LIMIT)
	LogDBQuery(r, ctx, QuerySearchForPin, time.Since(start))

	if err != nil {
		return []entity.SearchPin{}, err
	}
	var res []entity.SearchPin
	err = carta.Map(rows, &res)
	return res, err
}

func (r *DBRepository) SearchForBoards(ctx context.Context, query string) ([]entity.SearchBoard, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QuerySearchForBoard, query, SEARCH_LIMIT)
	LogDBQuery(r, ctx, QuerySearchForBoard, time.Since(start))

	if err != nil {
		return []entity.SearchBoard{}, err
	}
	var res []entity.SearchBoard
	err = carta.Map(rows, &res)
	return res, err
}
