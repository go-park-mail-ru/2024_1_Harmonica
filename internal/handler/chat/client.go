package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"log"
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
	hub     *Hub
	conn    *websocket.Conn
	message chan *ChatMessage //send chan []byte

	userId    entity.UserID
	wsConnKey string
}

func NewClient(hub *Hub, conn *websocket.Conn, userId entity.UserID, wsConnKey string) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		message:   make(chan *ChatMessage, 10), //send: make(chan []byte, 256), // создаем буферизованный канал для исходящих сообщений
		userId:    userId,
		wsConnKey: wsConnKey,
	}
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {

		fmt.Println("11111")

		log.Println(err) // логгер поменять тут
		return
	}

	ctx := r.Context()
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		// TODO исправить
		//log.Println(errs.ErrTypeConversion)
		userId = 3
	}
	wsConnKey := r.URL.Query().Get("ws_conn_key")

	client := NewClient(hub, conn, userId, wsConnKey)
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

				fmt.Println("88888")

				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- &chatMessage
	}
}

func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case chatMessage, ok := <-c.message:
			if !ok {

				fmt.Println("99999")

				errResponse := errs.ErrorResponse{
					Code:    errs.ErrorCodes[errs.ErrWSConnectionClosed].LocalCode,
					Message: errs.ErrWSConnectionClosed.Error(),
				}
				c.conn.WriteJSON(errResponse)
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteJSON(chatMessage)
			if err != nil {

				fmt.Println("new 111")

				log.Printf("error: %s", err)
				return
			}

			n := len(c.message)
			for i := 0; i < n; i++ {
				//w.Write(newline)
				//w.Write(<-c.message)
				chatMessage = <-c.message
				if chatMessage.ReceiverId == c.userId {
					c.conn.WriteJSON(<-c.message)
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error: %v", err)

				// тут надо?
				//errResponse := errs.ErrorResponse{
				//	Code:    errs.ErrorCodes[errs.ErrWSConnectionClosed].LocalCode,
				//	Message: errs.ErrWSConnectionClosed.Error(),
				//}
				//c.conn.WriteJSON(errResponse)
				return
			}
		}
	}
}
