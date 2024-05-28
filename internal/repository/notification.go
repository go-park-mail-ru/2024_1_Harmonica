package repository

import (
	"context"
	"github.com/jackskj/carta"
	"harmonica/internal/entity"
	"time"
)

const (
	QueryCreatePinNotification = `INSERT INTO public.notification (user_id, type, triggered_by_user_id, pin_id) 
	SELECT s.user_id, 'new_pin', $1, $2 FROM public.subscribe_on_person s WHERE s.followed_user_id = $1;`

	QueryCreateBaseNotification = `INSERT INTO public.notification (user_id, type, triggered_by_user_id, pin_id,
  	comment_id, message_id) VALUES ($1, $2, $3, NULLIF($4, 0), NULLIF($5, 0), NULLIF($6, 0)) RETURNING notification_id`

	QueryGetNotificationById = `SELECT notification_id, user_id, type, triggered_by_user_id, 
    COALESCE(pin_id, 0) AS pin_id, COALESCE(comment_id, 0) AS comment_id, COALESCE(message_id, 0) AS message_id, 
    created_at FROM public.notification WHERE notification_id = $1`

	QueryGetUnreadNotifications = `SELECT n.notification_id, n.user_id, n.type, COALESCE(n.pin_id, 0) AS pin_pin_id, 
       COALESCE(n.comment_id, 0) AS comment_comment_id, COALESCE(n.message_id, 0) AS message_message_id, 
       n.created_at, tu.user_id AS triggered_by_user_user_id, tu.nickname AS triggered_by_user_nickname, 
       tu.avatar_url AS triggered_by_user_avatar_url, COALESCE(p.content_url, '') AS pin_content_url,
       COALESCE(c.text, '') AS comment_text, COALESCE(m.sender_id, 0) AS message_sender_id, 
       COALESCE(m.receiver_id, 0) AS message_receiver_id, COALESCE(m.text, '') AS message_text
	FROM public.notification n
			LEFT JOIN public.user tu ON n.triggered_by_user_id = tu.user_id
			LEFT JOIN public.pin p ON n.pin_id = p.pin_id
			LEFT JOIN public.comment c ON n.comment_id = c.comment_id
			LEFT JOIN public.message m ON n.message_id = m.message_id
	WHERE n.user_id = $1 AND n.status = 'unread' ORDER BY n.created_at DESC;`

	QueryUpdateNotificationStatus = `UPDATE public.notification SET status = 'read' WHERE notification_id = $1`

	QueryUpdateAllNotificationsStatuses = `UPDATE public.notification SET status = 'read' 
    WHERE user_id = $1 AND status = 'unread'`
)

func (r *DBRepository) CreateNotification(ctx context.Context, notification entity.Notification) (entity.NotificationID, error) {
	notificationId := entity.NotificationID(0)
	if notification.Type == entity.NotificationTypeNewPin {
		start := time.Now()
		_, err := r.db.ExecContext(ctx, QueryCreatePinNotification, notification.TriggeredByUserId, notification.PinId)
		LogDBQuery(r, ctx, QueryCreatePinNotification, time.Since(start))
		// для случая нового пина (-> создается много уведомлений) возвращаем нулевой айди просто
		return notificationId, err
	}
	start := time.Now()
	err := r.db.QueryRowContext(ctx, QueryCreateBaseNotification, notification.UserId, notification.Type,
		notification.TriggeredByUserId, notification.PinId, notification.CommentId, notification.MessageId).
		Scan(&notificationId)
	LogDBQuery(r, ctx, QueryCreateBaseNotification, time.Since(start))
	return notificationId, err
}

func (r *DBRepository) GetNotificationById(ctx context.Context, notificationId entity.NotificationID) (entity.Notification, error) {
	notification := entity.Notification{}
	start := time.Now()
	err := r.db.QueryRowxContext(ctx, QueryGetNotificationById, notificationId).StructScan(&notification)
	LogDBQuery(r, ctx, QueryGetNotificationById, time.Since(start))
	return notification, err
}

func (r *DBRepository) GetUnreadNotifications(ctx context.Context, userId entity.UserID) (entity.Notifications, error) {
	notifications := entity.Notifications{}
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, QueryGetUnreadNotifications, userId)
	LogDBQuery(r, ctx, QueryGetUnreadNotifications, time.Since(start))
	if err != nil {
		return entity.Notifications{}, err
	}
	err = carta.Map(rows, &notifications.Notifications)
	if err != nil {
		return entity.Notifications{}, err
	}
	return notifications, nil
}

func (r *DBRepository) ReadNotification(ctx context.Context, notificationId entity.NotificationID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateNotificationStatus, notificationId)
	LogDBQuery(r, ctx, QueryUpdateNotificationStatus, time.Since(start))
	return err
}

func (r *DBRepository) ReadAllNotifications(ctx context.Context, userId entity.UserID) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateAllNotificationsStatuses, userId)
	LogDBQuery(r, ctx, QueryUpdateAllNotificationsStatuses, time.Since(start))
	return err
}
