package repository

import (
	"harmonica/internal/entity"

	"github.com/jackskj/carta"
)

//type PinDB struct {
//	db *sqlx.DB
//}
//
//func NewPinDB(db *sqlx.DB) *PinDB {
//	return &PinDB{db: db}
//}

const QueryGetPins = `SELECT user_id, nickname, pin_id, caption, content_url, click_url, created_at FROM public.pins
    INNER JOIN public.users ON public.pins.author_id=public.users.user_id ORDER BY created_at DESC LIMIT $1 OFFSET $2`

func (r *DBRepository) GetPins(limit, offset int) (entity.Pins, error) {
	result := entity.Pins{}
	rows, err := r.db.Query(QueryGetPins, limit, offset)
	if err != nil {
		return entity.Pins{}, err
	}
	err = carta.Map(rows, &result.Pins)
	if err != nil {
		return entity.Pins{}, err
	}
	return result, nil
}
