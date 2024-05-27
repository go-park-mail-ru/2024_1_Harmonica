package handler

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
	"time"
)

const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 8000                // Maximum message size allowed from peer.
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	//CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	message   chan *entity.WSMessage //send chan []byte
	userId    entity.UserID
	wsConnKey string
	logger    *zap.Logger
}

func NewClient(hub *Hub, conn *websocket.Conn, userId entity.UserID, wsConnKey string, l *zap.Logger) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		message:   make(chan *entity.WSMessage, 100), // буферизованный канал для исходящих сообщений
		userId:    userId,
		wsConnKey: wsConnKey,
		logger:    l,
	}
}

func (h *APIHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	// не забывать возвращать куки вместо r.URL.Query().Get("user_id") после постмана!!!
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrWSConnectionUpgrade))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrWSConnectionUpgrade))
		return
	}

	wsConnKey := r.URL.Query().Get("ws_conn_key")
	client := NewClient(h.hub, conn, userId, wsConnKey, h.logger)
	client.hub.register <- client

	go client.WriteMessage()
}

func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case messageFromChan, ok := <-c.message:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{}) // the hub closed the channel
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			// отправка сообщения
			messageToSend, senderId, receiverId := configureMessageToSend(messageFromChan)
			err := c.conn.WriteJSON(messageToSend)

			if err != nil {
				// возникла ошибка при отправке json -> простое сообщение мб отправится
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.logger.Error(
					errs.ErrWSConnectionClosed.Error(),
					zap.Int("local_error_code", errs.ErrorCodes[errs.ErrWSConnectionClosed].LocalCode),
					zap.String("general_error", err.Error()),
				)
				return
			}

			// add queued chat messages to the current websocket message
			n := len(c.message)
			for i := 0; i < n; i++ {
				messageFromChan = <-c.message
				messageToSend, senderId, receiverId = configureMessageToSend(messageFromChan)
				action := messageFromChan.Action

				//if (receiverId == c.userId || senderId == c.userId) && senderId != receiverId {
				if (action == entity.WSActionChatMessage && (receiverId == c.userId || senderId == c.userId) && senderId != receiverId) ||
					(action == entity.WSActionChatDraft && senderId == c.userId) ||
					(action != entity.WSActionChatMessage && action != entity.WSActionChatDraft && receiverId == c.userId) {

					// отправка сообщения
					err = c.conn.WriteJSON(messageToSend)

					if err != nil {
						//возникла ошибка при отправке json -> простое сообщение мб отправится
						_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
						c.logger.Error(
							errs.ErrWSConnectionClosed.Error(),
							zap.Int("local_error_code", errs.ErrorCodes[errs.ErrWSConnectionClosed].LocalCode),
							zap.String("general_error", err.Error()),
						)
						return
					}
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Error(
					errs.ErrWSConnectionClosed.Error(),
					zap.Int("local_error_code", errs.ErrorCodes[errs.ErrWSConnectionClosed].LocalCode),
					zap.String("general_error", err.Error()),
				)
				return
			}
		}
	}
}

func configureMessageToSend(message *entity.WSMessage) (entity.WSMessageToSend, entity.UserID, entity.UserID) {
	senderId := message.Payload.TriggeredByUser.UserId
	receiverId := message.Payload.UserId
	messageToSend := entity.WSMessageToSend{Action: message.Action}
	switch message.Action {
	case entity.WSActionChatMessage:
		messageToSend.Payload = entity.WSChatMessagePayload{
			SenderId:   senderId,
			ReceiverId: receiverId,
			Text:       message.Payload.Message.Text,
		}
	case entity.WSActionChatDraft:
		messageToSend.Payload = entity.WSChatMessagePayload{
			SenderId:   senderId,
			ReceiverId: receiverId,
			Text:       message.Payload.Message.Text,
		}
	case entity.WSActionNotificationSubscription:
		messageToSend.Payload = entity.WSSubscriptionNotificationPayload{
			UserId:          receiverId,
			TriggeredByUser: message.Payload.TriggeredByUser,
			CreatedAt:       message.Payload.CreatedAt,
		}
	case entity.WSActionNotificationNewPin:
		messageToSend.Payload = entity.WSNewPinNotificationPayload{
			UserId:          receiverId,
			TriggeredByUser: message.Payload.TriggeredByUser,
			Pin:             message.Payload.Pin,
			CreatedAt:       message.Payload.CreatedAt,
		}
	case entity.WSActionNotificationComment:
		messageToSend.Payload = entity.WSCommentNotificationPayload{
			UserId:          receiverId,
			TriggeredByUser: message.Payload.TriggeredByUser,
			Comment:         message.Payload.Comment,
			Pin:             message.Payload.Pin,
			CreatedAt:       message.Payload.CreatedAt,
		}
	}
	return messageToSend, senderId, receiverId
}
