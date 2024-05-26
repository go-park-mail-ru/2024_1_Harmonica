package repository

import (
	"context"
	"harmonica/internal/entity"
	"time"

	"github.com/jackskj/carta"
)

const (
	QueryGetComments = `SELECT public.user.user_id, avatar_url, nickname, comment_id, text FROM public.comment
    INNER JOIN public.user ON public.comment.user_id=public.user.user_id
	WHERE public.comment.pin_id = $1
	ORDER BY public.comment.created_at DESC`

	QueryAddComment = `INSERT INTO public.comment ("user_id", "pin_id", "text") VALUES($1, $2, $3)`
)

func (r *DBRepository) AddComment(ctx context.Context, comment string, pinId entity.PinID, userId entity.UserID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryAddComment, userId, pinId, comment)
	LogDBQuery(r, ctx, QueryAddComment, time.Since(start))
	return err
}

func (r *DBRepository) GetComments(ctx context.Context, pinId entity.PinID) (entity.GetCommentsResponse, error) {
	var response entity.GetCommentsResponse
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetComments, pinId)
	LogDBQuery(r, ctx, QueryGetDraft, time.Since(start))
	if err != nil {
		return entity.GetCommentsResponse{}, err
	}
	err = carta.Map(rows, &response.Comments)
	if err != nil {
		return entity.GetCommentsResponse{}, err
	}
	return response, nil
}
