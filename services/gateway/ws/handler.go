// Package ws 提供 WebSocket 处理器，用于升级连接并将客户端注册到 Hub。
package ws

import (
	"IM/pkg/auth"
	"IM/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler 负责处理 WebSocket 协议升级、认证与 Client 的注册。
type Handler struct {
	hub     *Hub
	jwtUtil *auth.JWTUtil
	redis   *redis.Client
}

// NewHandler 创建一个 Handler 实例。
func NewHandler(hub *Hub, jwtUtil *auth.JWTUtil, redis *redis.Client) *Handler {
	return &Handler{
		hub:     hub,
		jwtUtil: jwtUtil,
		redis:   redis,
	}
}

// HandleWebSocket 处理来自 HTTP 的 WebSocket 升级请求并注册客户端到 Hub。
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	userID, err := h.jwtUtil.ParseToken(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		hub:    h.hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}

	h.hub.RegisterClient(client)

	go client.writePump()
	go client.readPump()
}

// readPump 从 WebSocket 连接读取消息并在断开时注销客户端。
func (c *Client) readPump() {
	defer func() {
		c.hub.UnregisterClient(c)
	}()

	conn := c.conn.(*websocket.Conn)
	conn.SetReadLimit(65536)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		logger.Infof("Received message from %s: %s", c.userID, string(message))
	}
}

// writePump 负责按心跳向客户端发送 Ping，并将 Hub 广播的消息写入连接。
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		conn := c.conn.(*websocket.Conn)
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			conn := c.conn.(*websocket.Conn)
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			conn := c.conn.(*websocket.Conn)
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WebSocketHandler 返回一个 Gin 处理器，用于将请求交由 Handler 处理 WebSocket 连接。
func WebSocketHandler(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.HandleWebSocket(c.Writer, c.Request)
	}
}
