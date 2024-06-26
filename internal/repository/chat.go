package repository

import (
	"context"
	"github.com/jackskj/carta"
	"harmonica/internal/entity"
	"time"
)

const (
	QueryCreateMessage = `INSERT INTO public.message (sender_id, receiver_id, text) VALUES ($1, $2, $3);`

	QueryGetMessages = `SELECT sender_id, text,
    CASE WHEN status = 'read' THEN true ELSE false END AS message_read, sent_at FROM public.message
	WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1) ORDER BY sent_at DESC;`

	QueryGetUserChats = `SELECT user_user_id, user_nickname, user_avatar_url,
    chat_last_message_sender_id, chat_last_message_text, chat_last_message_message_read, chat_last_message_sent_at
	FROM (
		SELECT
             u.user_id AS user_user_id, u.nickname AS user_nickname, u.avatar_url AS user_avatar_url,
             m.sender_id AS chat_last_message_sender_id,
             m.text AS chat_last_message_text,
             (m.status = 'read') AS chat_last_message_message_read,
             m.sent_at AS chat_last_message_sent_at,
             ROW_NUMBER() OVER (
                 PARTITION BY
                     LEAST(m.sender_id, m.receiver_id),
                     GREATEST(m.sender_id, m.receiver_id)
                 ORDER BY m.sent_at DESC
                 ) AS message_rank
         FROM public.user u
                  JOIN public.message m ON u.user_id IN (m.sender_id, m.receiver_id)
         	--WHERE (m.sender_id = 2 OR m.receiver_id = 2) AND u.user_id != 2
         	WHERE ((m.sender_id = $1 OR m.receiver_id = $1) AND u.user_id != $1) OR  (m.sender_id = $1 AND m.receiver_id = m.sender_id)
     	) AS ranked_messages
	WHERE message_rank = 1
	ORDER BY chat_last_message_sent_at DESC;`
)

func (r *DBRepository) CreateMessage(ctx context.Context, message entity.Message) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryCreateMessage, message.SenderId, message.ReceiverId, message.Text)
	LogDBQuery(r, ctx, QueryCreateMessage, time.Since(start))
	return err
}

func (r *DBRepository) GetMessages(ctx context.Context, dialogUserId, authUserId entity.UserID) (entity.Messages, error) {
	messages := entity.Messages{}
	dialogUser, err := r.GetUserById(ctx, dialogUserId)
	if err != nil {
		return entity.Messages{}, err
	}
	messages.User = entity.UserFromChat{
		UserID:    dialogUser.UserID,
		Nickname:  dialogUser.Nickname,
		AvatarURL: dialogUser.AvatarURL,
	}
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetMessages, dialogUserId, authUserId)
	LogDBQuery(r, ctx, QueryGetMessages, time.Since(start))
	if err != nil {
		return entity.Messages{}, err
	}
	err = carta.Map(rows, &messages.Messages)
	return messages, err
}

func (r *DBRepository) GetUserChats(ctx context.Context, userId entity.UserID) (entity.UserChats, error) {
	chats := entity.UserChats{}
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetUserChats, userId)
	LogDBQuery(r, ctx, QueryGetUserChats, time.Since(start))
	if err != nil {
		return entity.UserChats{}, err
	}
	defer rows.Close()

	// carta.Map() тут не сработала
	for rows.Next() {
		var userChat entity.UserChat
		err = rows.Scan(
			&userChat.User.UserID,
			&userChat.User.Nickname,
			&userChat.User.AvatarURL,
			&userChat.LastMessage.SenderId,
			&userChat.LastMessage.Text,
			&userChat.LastMessage.IsRead,
			&userChat.LastMessage.SentAt,
		)
		if err != nil {
			return entity.UserChats{}, err
		}
		chats.Chats = append(chats.Chats, userChat)
	}

	if err = rows.Err(); err != nil {
		return entity.UserChats{}, err
	}

	return chats, nil
}
