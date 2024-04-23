package repository

import (
	"context"
	"github.com/jackskj/carta"
	"harmonica/internal/entity"
	"time"
)

const (
	QueryCreateMessage = `INSERT INTO public.message (sender_id, receiver_id, text) VALUES ($1, $2, $3);`
	QueryGetMessages   = `SELECT sender_id, receiver_id, text, status, sent_at FROM public.message
	WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1) ORDER BY sent_at DESC;`
)

func (r *DBRepository) CreateMessage(ctx context.Context, message entity.Message) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryCreateMessage, message.SenderId, message.ReceiverId, message.Text)
	LogDBQuery(r, ctx, QueryCreateMessage, time.Since(start))
	return err
}

func (r *DBRepository) GetMessages(ctx context.Context, firstUserId, secondUserId entity.UserID) (entity.Messages, error) {
	messages := entity.Messages{}
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetMessages, firstUserId, secondUserId)
	LogDBQuery(r, ctx, QueryGetMessages, time.Since(start))
	if err != nil {
		return entity.Messages{}, err
	}
	err = carta.Map(rows, &messages.Messages)
	if err != nil {
		return entity.Messages{}, err
	}
	return messages, nil
}
