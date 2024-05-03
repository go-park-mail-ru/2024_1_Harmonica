package repository

import (
	"context"
	"fmt"
	"harmonica/internal/entity"
	"time"

	"go.uber.org/zap"
)

const (
	QueryGetUserByEmail = `SELECT user_id, email, nickname, "password", avatar_url FROM public.user WHERE email=$1`

	QueryGetUserByNickname = `SELECT user_id, email, nickname, "password", avatar_url FROM public.user WHERE nickname=$1`

	QueryGetUserById = `SELECT user_id, email, nickname, "password", avatar_url FROM public.user WHERE user_id=$1`

	QueryRegisterUser = `INSERT INTO public.user ("email", "nickname", "password") VALUES($1, $2, $3)`

	QueryUpdateUserNickname = `UPDATE public.user SET nickname=$2 WHERE user_id=$1`

	QueryUpdateUserPassword = `UPDATE public.user SET "password"=$2 WHERE user_id=$1`

	QueryUpdateUserAvatar = `UPDATE public.user SET "avatar_url"=$2 WHERE user_id=$1`
)

func (r *DBRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	start := time.Now()
	rows, err := r.db.QueryxContext(ctx, QueryGetUserByEmail, email)
	LogDBQuery(r, ctx, QueryGetUserByEmail, time.Since(start))
	// TODO по-хорошему переписать на QueryRowxContext (сложность в том, чтобы переписать все методы, использующие этот)
	// так как он возвращает ошибку sql.ErrNoRows, а QueryxContext - нет
	if err != nil {
		return entity.User{}, err
	}
	var user entity.User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return entity.User{}, err
		}
		dx, dy, err := r.GetImageBounds(ctx, user.AvatarURL)
		if err != nil {
			fmt.Print("IM HER", err)
			return entity.User{}, err
		}
		user.AvatarDX = dx
		user.AvatarDY = dy
	}
	fmt.Print("IM HER")
	return user, nil
}

func (r *DBRepository) GetUserByNickname(ctx context.Context, nickname string) (entity.User, error) {
	start := time.Now()
	rows, err := r.db.QueryxContext(ctx, QueryGetUserByNickname, nickname)
	LogDBQuery(r, ctx, QueryGetUserByNickname, time.Since(start))
	if err != nil {
		return entity.User{}, err
	}
	var user entity.User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return entity.User{}, err
		}
		dx, dy, err := r.GetImageBounds(ctx, user.AvatarURL)
		if err != nil {
			return entity.User{}, err
		}
		user.AvatarDX = dx
		user.AvatarDY = dy
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
		dx, dy, err := r.GetImageBounds(ctx, user.AvatarURL)
		if err != nil {
			return entity.User{}, err
		}
		user.AvatarDX = dx
		user.AvatarDY = dy
	}
	return user, nil
}

func (r *DBRepository) RegisterUser(ctx context.Context, user entity.User) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryRegisterUser, user.Email, user.Nickname, user.Password)
	LogDBQuery(r, ctx, QueryRegisterUser, time.Since(start))
	return err
}

func (r *DBRepository) UpdateUser(ctx context.Context, user entity.User) error {
	if user.Nickname != "" {
		start := time.Now()
		_, err := r.db.ExecContext(ctx, QueryUpdateUserNickname, user.UserID, user.Nickname)
		LogDBQuery(r, ctx, QueryUpdateUserNickname, time.Since(start))
		if err != nil {
			return err
		}
	}
	if user.Password != "" {
		start := time.Now()
		_, err := r.db.ExecContext(ctx, QueryUpdateUserPassword, user.UserID, user.Password)
		LogDBQuery(r, ctx, QueryUpdateUserPassword, time.Since(start))
		if err != nil {
			return err
		}
	}
	if user.AvatarURL != "" {
		start := time.Now()
		_, err := r.db.ExecContext(ctx, QueryUpdateUserAvatar, user.UserID, user.AvatarURL)
		LogDBQuery(r, ctx, QueryUpdateUserAvatar, time.Since(start))
		if err != nil {
			return err
		}
	}
	return nil
}

func LogDBQuery(r *DBRepository, ctx context.Context, query string, duration time.Duration) {
	requestId := ctx.Value("request_id").(string)
	r.logger.Info("DB query handled",
		zap.String("request_id", requestId),
		zap.String("query", query),
		zap.String("duration", duration.String()),
	)
}
