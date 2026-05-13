package websocket

import (
	"IM/services/message-service/domain/entity"
	"sync"
)

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan *entity.Message
	Register   chan *Client
	Unregister chan *Client
	mutex      sync.RWMutex
}

type Client struct {
	ID   string
	Send chan *entity.Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan *entity.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client.ID] = client
			h.mutex.Unlock()

		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mutex.Unlock()

		case message := <-h.Broadcast:
			h.mutex.RLock()
			if client, ok := h.Clients[message.ReceiverID]; ok {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) GetOnlineClient(userID string) (*Client, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	client, ok := h.Clients[userID]
	return client, ok
}

func (h *Hub) RegisterClient(client *Client) {
	h.Register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.Unregister <- client
}

func (h *Hub) BroadcastMessage(msg *entity.Message) {
	h.Broadcast <- msg
}