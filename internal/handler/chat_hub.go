package handler

import (
	"fmt"
	"harmonica/internal/entity"
	"sync"
)

type Hub struct {
	//clients map[*Client]bool // Registered clients.
	clients    map[entity.UserID][]*Client
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

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.mu.Lock()
			hub.clients[client.userId] = append(hub.clients[client.userId], client)
			hub.mu.Unlock()
		case client := <-hub.unregister:
			hub.mu.Lock()
			if clients, ok := hub.clients[client.userId]; ok {
				for i, c := range clients {
					if c == client {
						hub.clients[client.userId] = append(clients[:i], clients[i+1:]...)
						close(client.message)
						break
					}
				}
			}
			hub.mu.Unlock()
		case chatMessage := <-hub.broadcast:

			fmt.Println("3")

			hub.mu.Lock()
			//for client := range hub.clients {
			//	//select ... тут был дефолтно
			//	if chatMessage.ReceiverId == client.userId {
			//		fmt.Println("YES")
			//		client.message <- chatMessage
			//	}
			//}
			if clients, ok := hub.clients[chatMessage.ReceiverId]; ok {
				for _, client := range clients {
					client.message <- chatMessage
				}
			}
			// это для того, чтобы сообщение, отправленное юзером, отображалось во всех его вкладках
			if clients, ok := hub.clients[chatMessage.SenderId]; ok {
				for _, client := range clients {
					client.message <- chatMessage
				}
			}
			hub.mu.Unlock()
		}
	}
}
