package handler

import (
	"harmonica/internal/entity"
	"sync"
)

type Hub struct {
	clients    map[entity.UserID][]*Client
	mu         sync.Mutex
	broadcast  chan *entity.ChatMessage // Inbound messages from the clients.
	register   chan *Client             // Register requests from the clients.
	unregister chan *Client             // Unregister requests from clients.
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *entity.ChatMessage),
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
			deleteClientFromHub(hub, client)
			hub.mu.Unlock()
		case chatMessage := <-hub.broadcast:
			senderId := chatMessage.Payload.SenderId
			receiverId := chatMessage.Payload.ReceiverId
			hub.mu.Lock()
			if clients, ok := hub.clients[receiverId]; ok {
				for _, client := range clients {
					select {
					case client.message <- chatMessage:
						// отправка прошла успешно
					default:
						close(client.message)
						deleteClientFromHub(hub, client)
					}
				}
			}
			// это для того, чтобы сообщение, отправленное юзером, отображалось во всех его вкладках
			clients, ok := hub.clients[senderId]
			if ok && receiverId != senderId {
				for _, client := range clients {
					select {
					case client.message <- chatMessage:
						// отправка прошла успешно
					default:
						close(client.message)
						deleteClientFromHub(hub, client)
					}
				}
			}
			hub.mu.Unlock()
		}
	}
}

func deleteClientFromHub(hub *Hub, client *Client) {
	if clients, ok := hub.clients[client.userId]; ok {
		for i, c := range clients {
			if c == client {
				hub.clients[client.userId] = append(clients[:i], clients[i+1:]...)
				close(client.message)
				break
			}
		}
		if len(hub.clients[client.userId]) == 0 {
			delete(hub.clients, client.userId)
		}
	}
}
