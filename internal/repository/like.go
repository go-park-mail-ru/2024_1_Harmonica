package repository

import (
	"context"
	"harmonica/internal/entity"
	"time"

	"github.com/jackskj/carta"
)

const (
	QuerySetLike       = `INSERT INTO public.like ("pin_id", "user_id") VALUES ($1, $2) ON CONFLICT DO NOTHING`
	QueryClearLike     = `DELETE FROM public.like WHERE pin_id=$1 AND user_id=$2`
	QueryGetUsersLiked = `SELECT nickname, avatar_url FROM public.user WHERE(user_id IN (SELECT user_id FROM public.like WHERE pin_id=$1)) LIMIT $2`
	QueryFindLike      = `SELECT EXISTS(SELECT pin_id, user_id FROM public.like WHERE pin_id=$1 AND user_id=$2)`
	QueryIsLiked       = `SELECT EXISTS(SELECT 1 FROM public.like WHERE user_id=$2 AND pin_id=$1)`
)

func (r *DBRepository) SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QuerySetLike, pinId, userId)
	LogDBQuery(r, ctx, QuerySetLike, time.Since(start))
	return err
}

func (r *DBRepository) ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryClearLike, pinId, userId)
	LogDBQuery(r, ctx, QueryClearLike, time.Since(start))
	return err
}

func (r *DBRepository) GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, error) {
	result := entity.UserList{}
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetUsersLiked, pinId, limit)
	LogDBQuery(r, ctx, QueryGetUsersLiked, time.Since(start))
	if err != nil {
		return entity.UserList{}, err
	}
	err = carta.Map(rows, &result.Users)
	if err != nil {
		return entity.UserList{}, err
	}
	return result, nil
}

func (r *DBRepository) CheckIsLiked(ctx context.Context, pinId entity.PinID, userId entity.UserID) (bool, error) {
	var res bool
	err := r.db.QueryRowxContext(ctx, QueryIsLiked, pinId, userId).Scan(&res)
	if err != nil {
		return false, err
	}
	return res, nil
}
