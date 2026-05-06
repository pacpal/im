package hub

import (
	"IM/server/model"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func newHub() *Hub {
	h := &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 5),
	}
	go h.Run()
	return h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
			}
			h.mu.Unlock()
		}
	}
}
func (h *Hub) Register(conn *websocket.Conn, uid string) *Client {
	c := &Client{
		UserID: uid,
		Conn:   conn,
		Send:   make(chan model.Message, 256)}
	h.register <- c
	return c
}
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}
func (h *Hub) GetOnlineClient(uid string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.clients[uid]
	return c, ok
}
func (h *Hub) GetGroupMembers(uid string) (*[]Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
}
