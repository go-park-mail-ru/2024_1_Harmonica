package repository

import (
	"context"
	"github.com/jackskj/carta"
	"harmonica/internal/entity"
)

const (
	QueryCreateRating = `INSERT INTO public.rating ("title", "rating_count", "user_nickname") VALUES ($1, $2, $3)`
	QueryGetRating    = `SELECT title, rating_count, user_nickname FROM public.rating`
)

func (r *DBRepository) GetRating(ctx context.Context) (entity.RatingList, error) {
	result := entity.RatingList{}
	rows, err := r.db.QueryContext(ctx, QueryGetRating)
	if err != nil {
		return entity.RatingList{}, err
	}
	err = carta.Map(rows, &result.Ratings)
	if err != nil {
		return entity.RatingList{}, err
	}
	return result, nil
}

func (r *DBRepository) RatingCreate(ctx context.Context, rating entity.Rating) error {
	err := r.db.QueryRowContext(ctx, QueryCreateRating, rating.Title, rating.RatingCount, rating.User).Scan()
	return err
}
