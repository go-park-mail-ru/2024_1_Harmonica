package repository

import (
	"context"
	"harmonica/internal/entity"
)

const (
	QueryGetUserByEmail     = `SELECT user_id, email, nickname, "password" FROM public.users WHERE email=$1`
	QueryGetUserByNickname  = `SELECT user_id, email, nickname, "password" FROM public.users WHERE nickname=$1`
	QueryGetUserById        = `SELECT user_id, email, nickname, "password" FROM public.users WHERE user_id=$1`
	QueryRegisterUser       = `INSERT INTO public.users ("email", "nickname", "password") VALUES($1, $2, $3)`
	QueryUpdateUserNickname = `UPDATE public.users SET nickname=$2 WHERE user_id=$1`
	QueryUpdateUserPassword = `UPDATE public.users SET "password"=$2 WHERE user_id=$1`
)

var emptyUser = entity.User{}

func (r *DBRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	rows, err := r.db.QueryxContext(ctx, QueryGetUserByEmail, email)
	//чтобы запрос не продолжался, если пользователь ушел

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

func (r *DBRepository) GetUserByNickname(ctx context.Context, nickname string) (entity.User, error) {
	rows, err := r.db.QueryxContext(ctx, QueryGetUserByNickname, nickname)
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

func (r *DBRepository) GetUserById(ctx context.Context, id int64) (entity.User, error) {
	rows, err := r.db.QueryxContext(ctx, QueryGetUserById, id)
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

func (r *DBRepository) RegisterUser(ctx context.Context, user entity.User) error {
	_, err := r.db.ExecContext(ctx, QueryRegisterUser, user.Email, user.Nickname, user.Password)
	return err
}

func (r *DBRepository) UpdateUser(ctx context.Context, user entity.User) error {
	if user.Nickname != "" {
		_, err := r.db.ExecContext(ctx, QueryUpdateUserNickname, user.UserID, user.Nickname)
		if err != nil {
			return err
		}
	}
	if user.Password != "" {
		_, err := r.db.ExecContext(ctx, QueryUpdateUserPassword, user.UserID, user.Password)
		if err != nil {
			return err
		}
	}
	return nil
}
