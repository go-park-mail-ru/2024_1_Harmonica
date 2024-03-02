package db

import "time"

var SQLStatements = map[string]string{
	"RegisterUser":   `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`,
	"GetUserByEmail": `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`,
	"GetPins":        `SELECT * FROM public.pins INNER JOIN public.users ON public.pins.author_id=public.users.user_id ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
}

// ------------ Users ------------

// ------------ Pins ------------
func (connector *DBConnector) GetPins(limit, offset int) ([]Pin, error) {
	result := []Pin{}
	rows, err := connector.db.Queryx(SQLStatements["GetPins"], limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		pinItem := new(Pin)
		err := rows.StructScan(pinItem)
		if err != nil {
			return nil, err
		}
		// Выглядит плохо, что делать?)
		pinItem.Password = ""
		pinItem.Email = ""
		pinItem.RegisterAt = time.Time{}

		result = append(result, *pinItem)
	}
	return result, nil
}
