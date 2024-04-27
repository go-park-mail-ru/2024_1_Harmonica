package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"harmonica/internal/entity"
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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

func NewChatMessage(text string, senderID, receiverID entity.UserID) *ChatMessage {
	return &ChatMessage{
		Text:       text,
		SenderId:   senderID,
		ReceiverId: receiverID,
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	//send chan []byte
	message chan *ChatMessage

	userId    entity.UserID
	wsConnKey string
}

func NewClient(hub *Hub, conn *websocket.Conn, userId entity.UserID, wsConnKey string) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		//send:      make(chan []byte, 256), // создаем буферизованный канал для исходящих сообщений
		message:   make(chan *ChatMessage, 10),
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

		//fmt.Println("22222")
		//
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
		//case chatMessage, ok := <-c.message:
		//
		//	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		//
		//	if !ok && chatMessage.ReceiverId != c.userId {
		//
		//		fmt.Println("99999")
		//
		//		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		//		return
		//	}
		//
		//	for {
		//		select {
		//		case msg, ok := <-c.message:
		//			if !ok {
		//				return
		//			}
		//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		//			err := c.conn.WriteJSON(msg)
		//			if err != nil {
		//				log.Printf("error: %s", err)
		//				return
		//			}
		//		default:
		//			break
		//		}
		//	}
		//
		//	//if chatMessage.ReceiverId == c.userId {
		//	//	// Если да, игнорируем это сообщение
		//	//	continue
		//	//}
		//
		//	//c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		//
		//	err := c.conn.WriteJSON(chatMessage)
		//	if err != nil {
		//
		//		fmt.Println("77777")
		//
		//		log.Printf("error: %s", err)
		//		return
		//	}
		//
		//	// Add queued chat messages to the current websocket message.
		//	n := len(c.message)
		//	for i := 0; i < n; i++ {
		//		//w.Write(newline)
		//		//w.Write(<-c.message)
		//		chatMessage = <-c.message
		//		if chatMessage.ReceiverId == c.userId {
		//			c.conn.WriteJSON(<-c.message)
		//		}
		//	}
		//
		//	if errClose := c.conn.Close(); errClose != nil {
		//		return
		//	}

		case chatMessage, ok := <-c.message:
			if !ok {

				fmt.Println("99999")

				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
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

			// Add queued chat messages to the current websocket message.
			//for {
			//	select {
			//	case msg, ok := <-c.message:
			//		if !ok {
			//			return
			//		}
			//		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			//		err := c.conn.WriteJSON(msg)
			//		if err != nil {
			//			log.Printf("error: %s", err)
			//			return
			//		}
			//	default:
			//		break
			//	}
			//}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error: %v", err)
				return
			}
		}
	}
}
