package entity

import "time"

type NotificationType string

const (
	NotificationTypeSubscription NotificationType = "subscription"
	NotificationTypeNewPin       NotificationType = "new_pin"
	NotificationTypeComment      NotificationType = "comment"
	NotificationTypeMessage      NotificationType = "message"
)

var NotificationTypes = []NotificationType{
	NotificationTypeSubscription, NotificationTypeNewPin, NotificationTypeComment, NotificationTypeMessage,
}

type Notification struct {
	NotificationID    int64            `db:"notification_id" json:"notification_id"`
	UserID            UserID           `db:"user_id" json:"user_id"`
	Type              NotificationType `db:"type" json:"type"`
	TriggeredByUserID UserID           `db:"triggered_by_user_id" json:"triggered_by_user_id"`
	PinID             PinID            `db:"pin_id" json:"pin_id"`
	CommentID         int64            `db:"comment_id" json:"comment_id"`
	MessageID         int64            `db:"message_id" json:"message_id"`
	CreatedAt         time.Time        `db:"created_at" json:"created_at"`
}
