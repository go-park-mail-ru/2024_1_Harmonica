package db

import (
	"harmonica/models"

	"github.com/jackskj/carta"
)

var SQLStatements = map[string]string{
	"RegisterUser":   `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`,
	"GetUserByEmail": `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`,
	"GetPins":        `SELECT user_id, nickname, pin_id, caption, content_url, click_url, created_at FROM public.pins INNER JOIN public.users ON public.pins.author_id=public.users.user_id ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
}

// ------------ Users ------------

// ------------ Pins ------------
func (connector *DBConnector) GetPins(limit, offset int) (models.Pins, error) {
	result := models.Pins{}
	rows, err := connector.db.Query(SQLStatements["GetPins"], limit, offset)
	if err != nil {
		return models.Pins{}, err
	}
	err = carta.Map(rows, &result.Pins)
	if err != nil {
		return models.Pins{}, err
	}
	return result, nil
}
