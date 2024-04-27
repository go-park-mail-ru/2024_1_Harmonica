package repository

import (
	"context"
	"github.com/jackskj/carta"
	"harmonica/internal/entity"
)

const (
	//QuerySetRating = `INSERT INTO public.rating ("")`
	QueryGetRating = `SELECT 'title', 'rating_count', 'user' FROM public.rating`
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
