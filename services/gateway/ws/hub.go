// Package ws 提供 WebSocket Hub 与 Client 管理，支持广播、点对点、多设备连接及优雅关闭。
package ws

import (
	"IM/pkg/logger"
	"sync"
	"time"
)

// Client 表示一个 WebSocket 客户端连接。
type Client struct {
	hub    *Hub
	conn   *websocketConn
	send   chan []byte
	userID string
}

// Hub 管理所有活跃客户端，支持注册、注销、广播、点对点发送及优雅关闭。
type Hub struct {
	clients    map[string]map[*Client]struct{} // userID -> clients set (支持多设备)
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	sendToUser chan userMessage
	shutdown   chan struct{}
	mu         sync.RWMutex
}

type userMessage struct {
	userID  string
	message []byte
}

// NewHub 创建并返回一个 Hub 实例。
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]struct{}),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sendToUser: make(chan userMessage, 256),
		shutdown:   make(chan struct{}),
	}
}

// Run 启动 Hub 事件循环。
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]struct{})
			}
			h.clients[client.userID][client] = struct{}{}
			h.mu.Unlock()
			logger.Infow("User connected", "component", "gateway_ws", "user", client.userID, "clients", len(h.clients[client.userID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				if _, ok2 := clients[client]; ok2 {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()
			logger.Infow("User disconnected", "component", "gateway_ws", "user", client.userID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, clients := range h.clients {
				for client := range clients {
					select {
					case client.send <- message:
					default:
						h.closeClient(client)
					}
				}
			}
			h.mu.RUnlock()

		case um := <-h.sendToUser:
			h.mu.RLock()
			clients := h.clients[um.userID]
			for client := range clients {
				select {
				case client.send <- um.message:
				default:
					h.closeClient(client)
				}
			}
			h.mu.RUnlock()

		case <-h.shutdown:
			h.mu.Lock()
			for _, clients := range h.clients {
				for client := range clients {
					close(client.send)
				}
			}
			h.clients = make(map[string]map[*Client]struct{})
			h.mu.Unlock()
			return
		}
	}
}

// Broadcast 向所有在线客户端广播消息。
func (h *Hub) Broadcast(message []byte) {
	select {
	case h.broadcast <- message:
	case <-time.After(5 * time.Second):
		logger.Warnw("Broadcast timeout", "component", "gateway_ws")
	}
}

// SendToUser 向指定用户的所有连接发送消息。
func (h *Hub) SendToUser(userID string, message []byte) {
	select {
	case h.sendToUser <- userMessage{userID: userID, message: message}:
	case <-time.After(5 * time.Second):
		logger.Warnw("SendToUser timeout", "component", "gateway_ws", "user", userID)
	}
}

// IsOnline 判断用户是否在线。
func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}

// GetOnlineUsers 返回当前在线用户列表。
func (h *Hub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]string, 0, len(h.clients))
	for userID, clients := range h.clients {
		if len(clients) > 0 {
			users = append(users, userID)
		}
	}
	return users
}

// OnlineCount 返回当前在线连接数。
func (h *Hub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	count := 0
	for _, clients := range h.clients {
		count += len(clients)
	}
	return count
}

// Stop 优雅关闭 Hub。
func (h *Hub) Stop() {
	close(h.shutdown)
}

func (h *Hub) closeClient(client *Client) {
	select {
	case <-client.send:
	default:
		close(client.send)
	}
}

// RegisterClient 将客户端注册到 Hub。
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient 从 Hub 注销客户端。
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}
