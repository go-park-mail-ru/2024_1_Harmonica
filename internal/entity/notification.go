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
	NotificationId    NotificationID   `db:"notification_id" json:"notification_id"`
	UserId            UserID           `db:"user_id" json:"user_id"`
	Type              NotificationType `db:"type" json:"type"`
	TriggeredByUserId UserID           `db:"triggered_by_user_id" json:"triggered_by_user_id"`
	PinId             PinID            `db:"pin_id" json:"pin_id"`
	CommentId         CommentID        `db:"comment_id" json:"comment_id"`
	MessageId         int64            `db:"message_id" json:"message_id"`
	CreatedAt         time.Time        `db:"created_at" json:"created_at"`
}

type TriggeredByUser struct {
	UserId    UserID `db:"user_id" json:"user_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
} // заполняется всегда, вне зависимости от типа уведомления;
// в случае, если тип - "subscription", заполняется только оно
// единственное, это и есть новый подписчик

type PinNotificationResponse struct {
	PinId      PinID  `db:"pin_id" json:"pin_id"`
	ContentUrl string `db:"content_url" json:"content_url"`
}

type CommentNotificationResponse struct {
	CommentId CommentID `db:"comment_id" json:"comment_id"`
	Text      string    `db:"text" json:"text"`
}

type MessageNotificationResponse struct {
	Text string `db:"text" json:"text"`
}

type NotificationResponse struct {
	NotificationId  NotificationID              `db:"notification_id" json:"notification_id"`
	UserId          UserID                      `db:"user_id" json:"user_id"`
	Type            NotificationType            `db:"type" json:"type"`
	TriggeredByUser TriggeredByUser             `db:"triggered_by_user" json:"triggered_by_user"`
	Pin             PinNotificationResponse     `db:"pin" json:"pin"`
	Comment         CommentNotificationResponse `db:"comment" json:"comment"`
	Message         MessageNotificationResponse `db:"message" json:"message"`
	CreatedAt       time.Time                   `db:"created_at" json:"created_at"`
}

type Notifications struct {
	Notifications []NotificationResponse `db:"notifications" json:"notifications"`
}
