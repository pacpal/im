package handler

import (
	"IM/services/message-service/application/service"
	"IM/services/message-service/interface/dto"
	ws "IM/services/message-service/infrastructure/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type MessageHttpHandler struct {
	msgSvc *service.MessageApplicationService
	upgrader websocket.Upgrader
}

func NewMessageHttpHandler(msgSvc *service.MessageApplicationService) *MessageHttpHandler {
	return &MessageHttpHandler{
		msgSvc: msgSvc,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *MessageHttpHandler) SendMessage(c *gin.Context) {
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	senderID, _ := c.Get("user_id")

	msg, err := h.msgSvc.SendMessage(c.Request.Context(), senderID.(string), req.ReceiverID, req.Content, req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "message sent",
		"data":    dto.ToMessageResponse(msg),
	})
}

func (h *MessageHttpHandler) GetMessage(c *gin.Context) {
	messageID := c.Param("id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "message id is required"})
		return
	}

	msg, err := h.msgSvc.GetMessage(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": dto.ToMessageResponse(msg),
	})
}

func (h *MessageHttpHandler) GetOfflineMessages(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limit := 50
	offset := 0

	msgs, err := h.msgSvc.GetOfflineMessages(c.Request.Context(), userID.(string), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"messages": dto.ToMessageResponseList(msgs),
	})
}

func (h *MessageHttpHandler) MarkAsRead(c *gin.Context) {
	messageID := c.Param("id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "message id is required"})
		return
	}

	err := h.msgSvc.MarkAsRead(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "message marked as read",
	})
}

func (h *MessageHttpHandler) MarkAllAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")

	err := h.msgSvc.MarkAllAsRead(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "all messages marked as read",
	})
}

func (h *MessageHttpHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	count, err := h.msgSvc.GetUnreadCount(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"unread_count": count,
	})
}

func (h *MessageHttpHandler) HandleWebSocket(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		ID:   userID,
		Send: make(chan *ws.Message, 256),
	}

	hub := h.msgSvc.GetHub()
	hub.RegisterClient(client)

	go func() {
		defer func() {
			hub.UnregisterClient(client)
			conn.Close()
		}()

		for {
			var msg ws.Message
			if err := conn.ReadJSON(&msg); err != nil {
				break
			}

			_, err := h.msgSvc.SendMessage(c.Request.Context(), userID, msg.ReceiverID, msg.Content, msg.Type)
			if err != nil {
				conn.WriteJSON(gin.H{"success": false, "error": err.Error()})
			}
		}
	}()

	go func() {
		defer conn.Close()
		for {
			select {
			case msg, ok := <-client.Send:
				if !ok {
					return
				}
				if err := conn.WriteJSON(msg); err != nil {
					return
				}
			}
		}
	}()
}

func (h *MessageHttpHandler) RegisterRoutes(router *gin.RouterGroup, msgSvc *service.MessageApplicationService) {
	h.msgSvc = msgSvc

	msgGroup := router.Group("/messages")
	{
		msgGroup.POST("", h.SendMessage)
		msgGroup.GET("/:id", h.GetMessage)
		msgGroup.GET("/offline", h.GetOfflineMessages)
		msgGroup.PUT("/:id/read", h.MarkAsRead)
		msgGroup.PUT("/read/all", h.MarkAllAsRead)
		msgGroup.GET("/unread/count", h.GetUnreadCount)
		msgGroup.GET("/ws", h.HandleWebSocket)
	}
}

type Message struct {
	ReceiverID string `json:"receiver_id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
}