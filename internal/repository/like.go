package repository

import (
	"context"
	"harmonica/internal/entity"

	"github.com/jackskj/carta"
)

const (
	QueryCreateLike   = `INSERT INTO public.like ("pin_id", "user_id") VALUES ($1, $2)`
	QueryDeleteLike   = `DELETE FROM public.like WHERE pin_id=$1 AND user_id=$2`
	QueryGetLikedPins = `SELECT user_id, avatar_url, nickname, pin_id, content_url FROM public.pin
    INNER JOIN public.user ON public.pin.author_id=public.user.user_id WHERE pin_id IN (SELECT pin_id FROM public.like WHERE user_id=$1) ORDER BY created_at DESC LIMIT $2`
	QueryGetUsersLiked = `SELECT nickname, avatar_url FROM public.user WHERE(user_id IN (SELECT user_id FROM public.like WHERE pin_id=$1)) LIMIT $2`
	QueryGetLikesCount = `SELECT COUNT(*) FROM public.like WHERE pin_id=$1`
	QueryFindLike      = `SELECT EXISTS(SELECT pin_id, user_id FROM public.like WHERE pin_id=$1 AND user_id=$2)`
)

func (r *DBRepository) CreateLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	_, err := r.db.ExecContext(ctx, QueryCreateLike, pinId, userId)
	return err
}

func (r *DBRepository) DeleteLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	_, err := r.db.ExecContext(ctx, QueryDeleteLike, pinId, userId)
	return err
}

func (r *DBRepository) GetLikedPins(ctx context.Context, userId entity.UserID, limit int) (entity.FeedPins, error) {
	result := entity.FeedPins{}
	rows, err := r.db.QueryContext(ctx, QueryGetLikedPins, userId, limit)
	if err != nil {
		return entity.FeedPins{}, err
	}
	err = carta.Map(rows, &result.Pins)
	if err != nil {
		return entity.FeedPins{}, err
	}
	return result, nil
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

func (r *DBRepository) GetLikesCount(ctx context.Context, pinId entity.PinID) (uint64, error) {
	var res uint64
	err := r.db.QueryRowContext(ctx, QueryGetLikesCount, pinId).Scan(&res)
	return res, err
}

func (r *DBRepository) FindLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) (bool, error) { // true if found
	res := false
	err := r.db.QueryRowxContext(ctx, QueryFindLike, pinId, userId).Scan(&res)
	return res, err
}
