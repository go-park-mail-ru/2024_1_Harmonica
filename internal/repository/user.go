package repository

import (
	"harmonica/internal/entity"
)

const (
	QueryGetUserByEmail    = `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`
	QueryGetUserByNickname = `SELECT user_id, email, nickname, "password" FROM public.users WHERE nickname=$1`
	QueryGetUserById       = `SELECT user_id, email, nickname, "password" FROM public.users WHERE user_id=$1`
	QueryRegisterUser      = `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`
)

func (r *DBRepository) GetUserByEmail(email string) (entity.User, error) {
	//rows, err := r.db.QueryxContext(ctx, QueryGetUserByEmail, email) // добавить
	//(чтобы запрос не продолжался, если пользователь ушел)
	rows, err := r.db.Queryx(QueryGetUserByEmail, email)
	emptyUser := entity.User{}

	if err != nil {
		return emptyUser, err
	}

	var user entity.User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return emptyUser, err
		}
	}

	return user, nil
}

func (r *DBRepository) GetUserByNickname(nickname string) (entity.User, error) {
	rows, err := r.db.Queryx(QueryGetUserByNickname, nickname)
	emptyUser := entity.User{}
	if err != nil {
		return emptyUser, err
	}

	var user entity.User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return emptyUser, err
		}
	}
	return user, nil
}

func (r *DBRepository) GetUserById(id int64) (entity.User, error) {
	rows, err := r.db.Queryx(QueryGetUserById, id)
	emptyUser := entity.User{}
	if err != nil {
		return emptyUser, err
	}

	var user entity.User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return emptyUser, err
		}
	}
	return user, nil
}

func (r *DBRepository) RegisterUser(user entity.User) error {
	_, err := r.db.Exec(QueryRegisterUser, user.Email, user.Nickname, user.Password)
	return err
}
