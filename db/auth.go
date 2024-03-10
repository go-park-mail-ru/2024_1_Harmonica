package db

const (
	QueryGetUserByEmail    = `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`
	QueryGetUserByNickname = `SELECT user_id, email, nickname, "password" FROM public.users WHERE nickname=$1`
	QueryGetUserById       = `SELECT user_id, email, nickname, "password" FROM public.users WHERE user_id=$1`
	QueryRegisterUser      = `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`
)

func (connector *Connector) GetUserByEmail(email string) (User, error) {
	rows, err := connector.db.Queryx(QueryGetUserByEmail, email)
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

func (connector *Connector) GetUserByNickname(nickname string) (User, error) {
	rows, err := connector.db.Queryx(QueryGetUserByNickname, nickname)
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
	rows, err := connector.db.Queryx(QueryGetUserById, id)
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
	_, err := connector.db.Exec(QueryRegisterUser, user.Email, user.Nickname, user.Password)
	return err
}
