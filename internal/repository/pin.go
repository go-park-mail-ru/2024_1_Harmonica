package repository

import (
	"context"
	"harmonica/internal/entity"
	"time"

	"github.com/jackskj/carta"
)

const (
	QueryGetPinsFeed = `SELECT user_id, avatar_url, nickname, pin_id, content_url FROM public.pin
    INNER JOIN public.user ON public.pin.author_id=public.user.user_id ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	QueryGetUserPins = `SELECT pin_id, content_url FROM public.pin INNER JOIN public.user ON 
	public.pin.author_id=public.user.user_id WHERE public.pin.author_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	QueryGetPinById = `SELECT user_id, avatar_url, nickname, pin_id, title, "description", content_url, click_url, created_at, allow_comments, 
	(SELECT COUNT(*) FROM public.like WHERE public.like.pin_id=public.pin.pin_id) AS likes_count
	FROM public.pin INNER JOIN public.user ON public.pin.author_id=public.user.user_id WHERE pin_id=$1`

	QueryCreatePin = `INSERT INTO public.pin ("author_id", "content_id", "click_url", "title", "description", "allow_comments") 
	VALUES($1, $2, $3, $4, $5, $6) RETURNING pin_id`

	QueryUpdatePin = `UPDATE public.pin SET allow_comments=$2, title=$3, "description"=$4, click_url=$5 WHERE pin_id=$1`
	QueryDeletePin = `DELETE FROM public.pin WHERE pin_id=$1`

	QueryCheckPinExistence = `SELECT EXISTS(SELECT 1 FROM public.pin WHERE pin_id=$1)`
)

func (r *DBRepository) GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, error) {
	result := entity.FeedPins{}
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetPinsFeed, limit, offset)
	LogDBQuery(r, ctx, QueryGetPinsFeed, time.Since(start))
	if err != nil {
		return entity.FeedPins{}, err
	}
	err = carta.Map(rows, &result.Pins)
	if err != nil {
		return entity.FeedPins{}, err
	}
	return result, nil
}

func (r *DBRepository) GetUserPins(ctx context.Context, authorId entity.UserID, limit, offset int) (entity.UserPins, error) {
	result := entity.UserPins{}
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetUserPins, authorId, limit, offset)
	LogDBQuery(r, ctx, QueryGetUserPins, time.Since(start))
	if err != nil {
		return entity.UserPins{}, err
	}
	err = carta.Map(rows, &result.Pins)
	if err != nil {
		return entity.UserPins{}, err
	}
	return result, nil
}

func (r *DBRepository) GetPinById(ctx context.Context, pinId entity.PinID) (entity.PinPageResponse, error) {
	result := entity.PinPageResponse{}
	start := time.Now()
	err := r.db.QueryRowxContext(ctx, QueryGetPinById, pinId).StructScan(&result)
	LogDBQuery(r, ctx, QueryGetPinById, time.Since(start))
	if err != nil {
		return entity.PinPageResponse{}, err
	}
	return result, nil
}

func (r *DBRepository) CreatePin(ctx context.Context, pin entity.Pin) (entity.PinID, error) {
	res := entity.PinID(0)
	start := time.Now()
	err := r.db.QueryRowContext(ctx, QueryCreatePin, pin.AuthorId, pin.ContentId, pin.ClickUrl, pin.Title,
		pin.Description, pin.AllowComments).Scan(&res)
	LogDBQuery(r, ctx, QueryCreatePin, time.Since(start))
	return res, err
}

func (r *DBRepository) UpdatePin(ctx context.Context, pin entity.Pin) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdatePin, pin.PinId, pin.AllowComments, pin.Title, pin.Description, pin.ClickUrl)
	LogDBQuery(r, ctx, QueryUpdatePin, time.Since(start))
	return err
}

func (r *DBRepository) DeletePin(ctx context.Context, id entity.PinID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryDeletePin, id)
	LogDBQuery(r, ctx, QueryDeletePin, time.Since(start))
	return err
}

func (r *DBRepository) CheckPinExistence(ctx context.Context, id entity.PinID) (bool, error) {
	var exists bool
	start := time.Now()
	err := r.db.QueryRowContext(ctx, QueryCheckPinExistence, id).Scan(&exists)
	LogDBQuery(r, ctx, QueryCheckPinExistence, time.Since(start))
	return exists, err
}
