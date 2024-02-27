package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var SQLStatements = map[string]string{
	"RegisterUser":   `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`,
	"GetUserByEmail": `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`,
	"GetAllPins":     `SELECT * FROM public.pins`,
	"GetPinsOfUser":  `SELECT * FROM public.pins WHERE author_id=$1`,
}

// ------------ Users ------------

// ------------ Pins ------------
func getPinsByRows(rows *sqlx.Rows) ([]Pin, error) {
	var result []Pin
	var pin = &Pin{}
	for rows.Next() {
		err := rows.StructScan(&pin)
		if err != nil {
			return nil, err
		}
		result = append(result, *pin)
	}
	return result, nil
}

func (handler *DBConnector) GetAllPins() ([]Pin, error) {
	rows, err := handler.db.Queryx(SQLStatements["GetAllPins"])
	fmt.Println(rows)
	if err != nil {
		return nil, err
	}
	return getPinsByRows(rows)
}

func (handler *DBConnector) GetPinsOfUser(userId int) ([]Pin, error) {
	rows, err := handler.db.Queryx(SQLStatements["GetPinsOfUser"], userId)
	if err != nil {
		return nil, err
	}
	return getPinsByRows(rows)
}
