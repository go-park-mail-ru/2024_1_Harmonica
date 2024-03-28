package repository

import (
	"context"
	"harmonica/internal/entity"

	"github.com/jackskj/carta"
)

const (
	QuerySetLike       = `INSERT INTO public.like ("pin_id", "user_id") VALUES ($1, $2) ON CONFLICT DO NOTHING`
	QueryClearLike     = `DELETE FROM public.like WHERE pin_id=$1 AND user_id=$2`
	QueryGetUsersLiked = `SELECT nickname, avatar_url FROM public.user WHERE(user_id IN (SELECT user_id FROM public.like WHERE pin_id=$1)) LIMIT $2`
	QueryFindLike      = `SELECT EXISTS(SELECT pin_id, user_id FROM public.like WHERE pin_id=$1 AND user_id=$2)`
)

func (r *DBRepository) SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	_, err := r.db.ExecContext(ctx, QuerySetLike, pinId, userId)
	return err
}

func (r *DBRepository) ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	_, err := r.db.ExecContext(ctx, QueryClearLike, pinId, userId)
	return err
}

func (r *DBRepository) GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, error) {
	result := entity.UserList{}
	rows, err := r.db.QueryContext(ctx, QueryGetUsersLiked, pinId, limit)
	if err != nil {
		return entity.UserList{}, err
	}
	err = carta.Map(rows, &result.Users)
	if err != nil {
		return entity.UserList{}, err
	}
	return result, nil
}
