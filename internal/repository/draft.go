package repository

import (
	"context"
	"harmonica/internal/entity"
	"time"
)

const (
	QueryGetDraft = `SELECT text FROM public.draft WHERE sender_id = $2 AND receiver_id = $1;`

	QueryUpdateDraft = `INSERT INTO public.draft (sender_id, receiver_id, text) VALUES ($1, $2, $3) 
	ON CONFLICT (sender_id, receiver_id) DO UPDATE SET text = EXCLUDED.text;`
)

func (r *DBRepository) GetDraft(ctx context.Context, receiverId, senderId entity.UserID) (entity.DraftResponse, error) {
	draft := entity.DraftResponse{}
	start := time.Now()
	err := r.db.QueryRowxContext(ctx, QueryGetDraft, receiverId, senderId).StructScan(&draft)
	LogDBQuery(r, ctx, QueryGetDraft, time.Since(start))
	return draft, err
}

func (r *DBRepository) UpdateDraft(ctx context.Context, draft entity.Draft) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateDraft, draft.SenderId, draft.ReceiverId, draft.Text)
	LogDBQuery(r, ctx, QueryUpdateDraft, time.Since(start))
	return err
}
