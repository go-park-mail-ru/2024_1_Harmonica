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
		message:   make(chan *entity.WSMessage, 100), // создаем буферизованный канал для исходящих сообщений
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
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrWSConnectionUpgrade))
		return
	}

	wsConnKey := r.URL.Query().Get("ws_conn_key")
	client := NewClient(h.hub, conn, userId, wsConnKey, h.logger)
	client.hub.register <- client

	go client.WriteMessage()
	go client.ReadMessage()
}

func (c *Client) ReadMessage() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var chatMessage entity.WSMessage
		err := c.conn.ReadJSON(&chatMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error(
					errs.ErrWSConnectionClosed.Error(),
					zap.Int("local_error_code", errs.ErrorCodes[errs.ErrWSConnectionClosed].LocalCode),
					zap.String("general_error", err.Error()),
				)
			}
			break
		}
		chatMessage.Payload.SenderId = c.userId
		c.hub.broadcast <- &chatMessage
	}
}

func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		//c.hub.unregister <- c //почему это тут не нужно?
		c.conn.Close()
	}()
	for {
		select {
		case chatMessage, ok := <-c.message:
			if !ok {
				// The hub closed the channel
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteJSON(chatMessage)
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

			// Add queued chat messages to the current websocket message.
			n := len(c.message)
			for i := 0; i < n; i++ {
				chatMessage = <-c.message
				senderId := chatMessage.Payload.SenderId
				receiverId := chatMessage.Payload.ReceiverId
				action := chatMessage.Action
				//if chatMessage.ReceiverId == c.userId {
				if (action == entity.ActionMessage && (receiverId == c.userId ||
					senderId == c.userId) && senderId != receiverId) ||
					(action == entity.ActionDraft && senderId == c.userId) {
					err = c.conn.WriteJSON(<-c.message)
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
				//не понятно, нужно и возможно ли тут как-то оповещать пользователя о закрытии соединения,
				//так как ошибка возникла как раз при попытке отправки сообщения
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
