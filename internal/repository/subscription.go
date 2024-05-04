package repository

import (
	"context"
	"github.com/jackskj/carta"
	"harmonica/internal/entity"
	"time"
)

const (
	QueryAddSubscriptionToUser = `INSERT INTO public.subscribe_on_person (user_id, followed_user_id) VALUES ($1, $2)`

	QueryDeleteSubscriptionToUser = `DELETE FROM public.subscribe_on_person WHERE user_id=$1 AND followed_user_id=$2`

	GetSubscriptionsInfo = `SELECT (SELECT COUNT(*) FROM public.subscribe_on_person WHERE followed_user_id = $1) 
    AS subscribers_count, (SELECT COUNT(*) FROM public.subscribe_on_person WHERE user_id = $1) AS subscriptions_count,
    CASE WHEN EXISTS (SELECT 1 FROM public.subscribe_on_person WHERE user_id = $2 AND followed_user_id = $1) 
    THEN true ELSE false END AS is_subscribed;`

	QueryGetUserSubscribers = `SELECT u.user_id, u.email, u.nickname, u.avatar_url FROM public.user u
	JOIN public.subscribe_on_person s ON u.user_id = s.user_id WHERE s.followed_user_id = $1 ORDER BY u.user_id DESC;`

	QueryGetUserSubscriptions = `SELECT u.user_id, u.email, u.nickname, u.avatar_url FROM public.user u
	JOIN public.subscribe_on_person s ON u.user_id = s.followed_user_id WHERE s.user_id = $1 ORDER BY u.user_id DESC;`
)

func (r *DBRepository) AddSubscriptionToUser(ctx context.Context, userId, subscribeUserId entity.UserID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryAddSubscriptionToUser, userId, subscribeUserId)
	LogDBQuery(r, ctx, QueryAddSubscriptionToUser, time.Since(start))
	return err
}

func (r *DBRepository) DeleteSubscriptionToUser(ctx context.Context, userId, unsubscribeUserId entity.UserID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryDeleteSubscriptionToUser, userId, unsubscribeUserId)
	LogDBQuery(r, ctx, QueryDeleteSubscriptionToUser, time.Since(start))
	return err
}

func (r *DBRepository) GetSubscriptionsInfo(ctx context.Context, userToGetInfoId, userId entity.UserID) (entity.UserProfileResponse, error) {
	start := time.Now()
	var userProfile entity.UserProfileResponse
	err := r.db.QueryRowxContext(ctx, GetSubscriptionsInfo, userToGetInfoId, userId).StructScan(&userProfile)
	LogDBQuery(r, ctx, GetSubscriptionsInfo, time.Since(start))
	return userProfile, err
}

func (r *DBRepository) GetUserSubscribers(ctx context.Context, userId entity.UserID) (entity.UserSubscribers, error) {
	start := time.Now()
	var subscribers entity.UserSubscribers
	rows, err := r.db.QueryContext(ctx, QueryGetUserSubscribers, userId)
	LogDBQuery(r, ctx, QueryGetUserSubscribers, time.Since(start))
	if err != nil {
		return entity.UserSubscribers{}, err
	}
	err = carta.Map(rows, &subscribers.Subscribers)
	return subscribers, err
}

func (r *DBRepository) GetUserSubscriptions(ctx context.Context, userId entity.UserID) (entity.UserSubscriptions, error) {
	start := time.Now()
	var subscriptions entity.UserSubscriptions
	rows, err := r.db.QueryContext(ctx, QueryGetUserSubscriptions, userId)
	LogDBQuery(r, ctx, QueryGetUserSubscriptions, time.Since(start))
	if err != nil {
		return entity.UserSubscriptions{}, err
	}
	err = carta.Map(rows, &subscriptions.Subscriptions)
	return subscriptions, err
}
