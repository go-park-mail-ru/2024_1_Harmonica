package db

import (
	"harmonica/models"

	"github.com/jackskj/carta"
)

var SQLStatements = map[string]string{
	"RegisterUser":   `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`,
	"GetUserByEmail": `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`,
	"GetPins":        `SELECT user_id, nickname, pin_id, caption, content_url, click_url, created_at FROM public.pins INNER JOIN public.users ON public.pins.author_id=public.users.user_id ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
	"GetUserById":    `SELECT user_id, email, nickname, "password" FROM public.users WHERE user_id=$1`,
}

// ------------ Users ------------

func (connector *DBConnector) GetUserByEmail(email string) (User, error) {
	rows, err := connector.db.Queryx(SQLStatements["GetUserByEmail"], email)
	emptyUser := User{}
	if err != nil {
		return emptyUser, err
	}

	var user User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return emptyUser, err
		}
	}
	return user, nil
}

// дублирование кода, но мне кажется, что так лучше для понятности (?)
func (connector *DBConnector) GetUserById(id int64) (User, error) {
	rows, err := connector.db.Queryx(SQLStatements["GetUserById"], id)
	emptyUser := User{}
	if err != nil {
		return emptyUser, err
	}

	var user User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return emptyUser, err
		}
	}
	return user, nil
}

func (connector *DBConnector) RegisterUser(user User) error {
	_, err := connector.db.Exec(SQLStatements["RegisterUser"], user.Email, user.Nickname, user.Password)
	return err
}

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
