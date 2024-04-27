package chat

import "fmt"

type Hub struct {
	clients map[*Client]bool // Registered clients.
	// TODO надо либо мьютексы, либо sync.Map

	//clients2 sync.Map

	//clients map[int]map[string]*Client - ПРИКОЛЬНАЯ ИДЕЯ
	broadcast  chan *ChatMessage // Inbound messages from the clients.
	register   chan *Client      // Register requests from the clients.
	unregister chan *Client      // Unregister requests from clients.
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *ChatMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			//h.clients.Store(client, true)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.message)
			}
			//h.clients.Delete(client)
			//close(client.send)
		case chatMessage := <-h.broadcast:
			for client := range h.clients {
				//select {
				//case client.message <- message:
				//default:
				//	close(client.message)
				//	delete(h.clients, client)
				//}

				if chatMessage.ReceiverId == client.userId {

					fmt.Println("YES")

					client.message <- chatMessage
				}
			}
		}
	}
}
