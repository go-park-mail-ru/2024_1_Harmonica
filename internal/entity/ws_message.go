package entity

import "time"

type WSAction string

const (
	WSActionChatMessage              WSAction = "CHAT_MESSAGE"
	WSActionChatDraft                WSAction = "CHAT_DRAFT"
	WSActionNotificationSubscription WSAction = "NOTIFICATION_SUBSCRIPTION"
	WSActionNotificationNewPin       WSAction = "NOTIFICATION_NEW_PIN"
	WSActionNotificationComment      WSAction = "NOTIFICATION_COMMENT"
)

var WSActions = []WSAction{WSActionChatMessage, WSActionChatDraft, WSActionNotificationSubscription,
	WSActionNotificationNewPin, WSActionNotificationComment}

// структуры для внутренней отправки в канал broadcast

type WSMessagePayload struct {
	UserId          UserID                      `db:"user_id" json:"user_id"`
	TriggeredByUser TriggeredByUser             `db:"triggered_by_user" json:"triggered_by_user"`
	NotificationId  NotificationID              `db:"notification_id" json:"notification_id"`
	Pin             PinNotificationResponse     `db:"pin" json:"pin"`
	Comment         CommentNotificationResponse `db:"comment" json:"comment"`
	Message         MessageNotificationResponse `db:"message" json:"message"`
	CreatedAt       time.Time                   `db:"created_at" json:"created_at"`
}

type WSMessage struct {
	Action  WSAction         `json:"action"`
	Payload WSMessagePayload `json:"payload"`
}

// структуры для отправки непосредственно в вебсокет-соединение

type WSSubscriptionNotificationPayload struct {
	NotificationId  NotificationID  `db:"notification_id" json:"notification_id"`
	UserId          UserID          `db:"user_id" json:"user_id"`
	TriggeredByUser TriggeredByUser `db:"triggered_by_user" json:"triggered_by_user"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
}

type WSNewPinNotificationPayload struct {
	NotificationId  NotificationID          `db:"notification_id" json:"notification_id"`
	UserId          UserID                  `db:"user_id" json:"user_id"`
	TriggeredByUser TriggeredByUser         `db:"triggered_by_user" json:"triggered_by_user"`
	Pin             PinNotificationResponse `db:"pin" json:"pin"`
	CreatedAt       time.Time               `db:"created_at" json:"created_at"`
}

type WSCommentNotificationPayload struct {
	NotificationId  NotificationID              `db:"notification_id" json:"notification_id"`
	UserId          UserID                      `db:"user_id" json:"user_id"`
	TriggeredByUser TriggeredByUser             `db:"triggered_by_user" json:"triggered_by_user"`
	Comment         CommentNotificationResponse `db:"comment" json:"comment"`
	Pin             PinNotificationResponse     `db:"pin" json:"pin"` // пин, к которому написали комментарий
	CreatedAt       time.Time                   `db:"created_at" json:"created_at"`
}

type WSChatMessagePayload struct {
	Text       string `json:"text"`
	SenderId   UserID `json:"sender_id"`
	ReceiverId UserID `json:"receiver_id"`
}

type WSMessageToSend struct {
	Action  WSAction    `json:"action"`
	Payload interface{} `json:"payload"`
}
