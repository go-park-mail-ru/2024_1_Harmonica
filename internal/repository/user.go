package repository

import (
	"github.com/jmoiron/sqlx"
	"harmonica/internal/entity"
)

const (
	QueryGetUserByEmail    = `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`
	QueryGetUserByNickname = `SELECT user_id, email, nickname, "password" FROM public.users WHERE nickname=$1`
	QueryGetUserById       = `SELECT user_id, email, nickname, "password" FROM public.users WHERE user_id=$1`
	QueryRegisterUser      = `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`
)

type UserDB struct {
	db *sqlx.DB
}

func NewUserDB(db *sqlx.DB) *UserDB {
	return &UserDB{db: db}
}

func (u *UserDB) GetUserByEmail(email string) (entity.User, error) {
	rows, err := u.db.Queryx(QueryGetUserByEmail, email)
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

func (u *UserDB) GetUserByNickname(nickname string) (entity.User, error) {
	rows, err := u.db.Queryx(QueryGetUserByNickname, nickname)
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

func (u *UserDB) GetUserById(id int64) (entity.User, error) {
	rows, err := u.db.Queryx(QueryGetUserById, id)
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

func (u *UserDB) RegisterUser(user entity.User) error {
	_, err := u.db.Exec(QueryRegisterUser, user.Email, user.Nickname, user.Password)
	return err
}
