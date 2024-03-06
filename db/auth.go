package db

var sqlAuthStatements = map[string]string{
	"GetUserByEmail": `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`,
	"GetUserById":    `SELECT user_id, email, nickname, "password" FROM public.users WHERE user_id=$1`,
	"RegisterUser":   `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`,
}

func (connector *Connector) GetUserByEmail(email string) (User, error) {
	rows, err := connector.db.Queryx(sqlAuthStatements["GetUserByEmail"], email)
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

func (connector *Connector) GetUserById(id int64) (User, error) {
	rows, err := connector.db.Queryx(sqlAuthStatements["GetUserById"], id)
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

func (connector *Connector) RegisterUser(user User) error {
	_, err := connector.db.Exec(sqlAuthStatements["RegisterUser"], user.Email, user.Nickname, user.Password)
	return err
}
