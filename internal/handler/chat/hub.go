package chat

import (
	"harmonica/internal/entity"
	"sync"
)

type Hub struct {
	//clients map[*Client]bool // Registered clients.
	clients    map[entity.UserID][]*Client // ПРИКОЛЬНАЯ ИДЕЯ
	mu         sync.Mutex
	broadcast  chan *ChatMessage // Inbound messages from the clients.
	register   chan *Client      // Register requests from the clients.
	unregister chan *Client      // Unregister requests from clients.
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *ChatMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[entity.UserID][]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			//h.clients[client] = true
			h.mu.Lock()
			h.clients[client.userId] = append(h.clients[client.userId], client)
			h.mu.Unlock()
			//h.clients.Store(client.userId, client)
		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userId]; ok {
				for i, c := range clients {
					if c == client {
						h.clients[client.userId] = append(clients[:i], clients[i+1:]...)
						close(client.message)
						break
					}
				}
				//delete(h.clients, client.userId)
				//close(client.message)
			}
			h.mu.Unlock()
			//h.clients.Delete(client)
			//close(client.send)
		case chatMessage := <-h.broadcast:
			h.mu.Lock()
			//for client := range h.clients {
			//	//select ... тут был дефолтно
			//	if chatMessage.ReceiverId == client.userId {
			//		fmt.Println("YES")
			//		client.message <- chatMessage
			//	}
			//}
			if clients, ok := h.clients[chatMessage.ReceiverId]; ok {
				for _, client := range clients {
					client.message <- chatMessage
				}
			}
			h.mu.Unlock()
		}
	}
}
