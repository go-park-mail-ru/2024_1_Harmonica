package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"log"
	"net/http"
	"strconv"
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
	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	//CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatMessage struct {
	Text       string        `json:"text"`
	SenderId   entity.UserID `json:"sender_id"`
	ReceiverId entity.UserID `json:"receiver_id"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	message   chan *ChatMessage //send chan []byte
	userId    entity.UserID
	wsConnKey string
	logger    *zap.Logger
}

func NewClient(hub *Hub, conn *websocket.Conn, userId entity.UserID, wsConnKey string, l *zap.Logger) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		message:   make(chan *ChatMessage, 100), //send: make(chan []byte, 256), // создаем буферизованный канал для исходящих сообщений
		userId:    userId,
		wsConnKey: wsConnKey,
		logger:    l,
	}
}

func (h *APIHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	//cookie := r.Header.Get("Cookie") //если понадобится, это можно будет передать в responseHeader
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrWSConnectionUpgrade))
		return
	}

	//userId, ok := ctx.Value("user_id").(entity.UserID)
	//if !ok {
	//	WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrWSConnectionUpgrade))
	//	return
	//}

	// в ws-протоколе нет кук!!!!! мб можно прокинуть?
	userIdString := r.URL.Query().Get("user_id")
	userIdInt, err := strconv.Atoi(userIdString)
	if err != nil {
		// TODO исправить
		log.Println(errs.ErrTypeConversion)
	}
	userId := entity.UserID(userIdInt)

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
		var chatMessage ChatMessage
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
		chatMessage.SenderId = c.userId
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
				// была идея отправлять так, но потом поняла, что это фигня. так ведь?
				//errResponse := errs.ErrorResponse{
				//	Code:    errs.ErrorCodes[errs.ErrWSConnectionClosed].LocalCode,
				//	Message: errs.ErrWSConnectionClosed.Error(),
				//}
				//c.conn.WriteJSON(errResponse)
				//return

				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			fmt.Println("2")

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteJSON(chatMessage)
			if err != nil {
				//возникла ошибка при отправке json -> простое сообщение мб отправится
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
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
				//if chatMessage.ReceiverId == c.userId {
				if chatMessage.ReceiverId == c.userId || chatMessage.SenderId == c.userId {
					err = c.conn.WriteJSON(<-c.message)
					if err != nil {
						//возникла ошибка при отправке json -> простое сообщение мб отправится
						c.conn.WriteMessage(websocket.CloseMessage, []byte{})
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
