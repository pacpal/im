// Package ws 提供 WebSocket 处理器，用于升级连接、认证、心跳及消息读写泵。
package ws

import (
	"IM/pkg/auth"
	"IM/pkg/logger"
	"context"
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

// websocketConn 是对 *websocket.Conn 的简单包装，避免在 Hub 中使用 interface{}。
type websocketConn struct {
	*websocket.Conn
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
		logger.Errorw("Failed to upgrade connection", "component", "gateway_ws", "err", err)
		return
	}

	// 设置在线状态到 Redis
	ctx := r.Context()
	if h.redis != nil {
		if err := h.redis.SAdd(ctx, "online_users", userID).Err(); err != nil {
			logger.Warnw("Failed to set online status", "component", "gateway_ws", "user", userID, "err", err)
		}
	}

	client := &Client{
		hub:    h.hub,
		conn:   &websocketConn{conn},
		send:   make(chan []byte, 256),
		userID: userID,
	}

	h.hub.RegisterClient(client)

	go client.writePump()
	go client.readPump(h.redis)
}

// readPump 从 WebSocket 连接读取消息并在断开时注销客户端。
func (c *Client) readPump(redisClient *redis.Client) {
	defer func() {
		c.hub.UnregisterClient(c)
		if redisClient != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			// 如果该用户没有其它连接，再从 Redis 移除
			if !c.hub.IsOnline(c.userID) {
				if err := redisClient.SRem(ctx, "online_users", c.userID).Err(); err != nil {
					logger.Warnw("Failed to remove online status", "component", "gateway_ws", "user", c.userID, "err", err)
				}
			}
		}
	}()

	c.conn.SetReadLimit(65536)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Warnw("WebSocket unexpected close", "component", "gateway_ws", "user", c.userID, "err", err)
			}
			break
		}

		logger.Infow("Received WS message", "component", "gateway_ws", "user", c.userID, "message", string(message))
		// TODO: 根据业务解析消息并处理（如心跳 ack、发送消息等）
	}
}

// writePump 负责按心跳向客户端发送 Ping，并将 Hub 下发的消息写入连接。
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 合并待发送的后续消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
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
