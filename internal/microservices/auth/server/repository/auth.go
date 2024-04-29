package repository

import (
	"context"
	"harmonica/internal/entity"
	"time"
)

const (
	QueryGetUserByEmail = `SELECT user_id, email, nickname, "password", avatar_url FROM public.user WHERE email=$1`
	QueryGetUserById    = `SELECT user_id, email, nickname, "password", avatar_url FROM public.user WHERE user_id=$1`
)

func (r *DBRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	start := time.Now()
	rows, err := r.db.QueryxContext(ctx, QueryGetUserByEmail, email)
	LogDBQuery(r, ctx, QueryGetUserByEmail, time.Since(start))

	if err != nil {
		return entity.User{}, err
	}
	var user entity.User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return entity.User{}, err
		}
	}
	return user, nil
}

func (r *DBRepository) GetUserById(ctx context.Context, id entity.UserID) (entity.User, error) {
	start := time.Now()
	rows, err := r.db.QueryxContext(ctx, QueryGetUserById, id)
	LogDBQuery(r, ctx, QueryGetUserById, time.Since(start))
	if err != nil {
		return entity.User{}, err
	}
	var user entity.User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return entity.User{}, err
		}
	}
	return user, nil
}
