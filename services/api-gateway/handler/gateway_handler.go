package handler

import (
	"IM/services/api-gateway/config"
	"IM/services/api-gateway/middleware"
	"IM/services/api-gateway/proxy"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GatewayHandler struct {
	proxy *proxy.ServiceProxy
	cfg   *config.Config
}

func NewGatewayHandler(p *proxy.ServiceProxy, cfg *config.Config) *GatewayHandler {
	return &GatewayHandler{
		proxy: p,
		cfg:   cfg,
	}
}

func (h *GatewayHandler) RegisterRoutes(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	router.GET("/health", h.HealthCheck)

	api := router.Group("/api/v1")
	{
		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users/:id", h.GetUser)
			protected.GET("/users/:id/friends", h.GetFriends)
			protected.POST("/users/friends", h.AddFriend)
			protected.POST("/users/friend-requests/accept", h.AcceptFriendRequest)
			protected.GET("/users/friend-requests", h.GetPendingFriendRequests)

			protected.POST("/groups", h.CreateGroup)
			protected.GET("/groups/:id", h.GetGroup)
			protected.GET("/groups/:id/members", h.GetGroupMembers)
			protected.GET("/users/:id/groups", h.GetUserGroups)
			protected.POST("/groups/join", h.JoinGroup)
			protected.POST("/groups/join/accept", h.AcceptGroupJoinRequest)
			protected.DELETE("/groups/:id/leave", h.LeaveGroup)

			protected.POST("/messages", h.SendMessage)
			protected.GET("/messages/:id", h.GetMessage)
			protected.GET("/messages/offline", h.GetOfflineMessages)
			protected.PUT("/messages/:id/read", h.MarkAsRead)
			protected.PUT("/messages/read/all", h.MarkAllAsRead)
			protected.GET("/messages/unread/count", h.GetUnreadCount)
			protected.GET("/messages/ws", h.WebSocket)
		}
	}
}

func (h *GatewayHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "api-gateway",
	})
}

func (h *GatewayHandler) Register(c *gin.Context) {
	var req struct {
		Tele     string `json:"tele"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(err.Error()))
		return
	}

	resp, err := h.proxy.Register(c.Request.Context(), req.Tele, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.User))
}

func (h *GatewayHandler) Login(c *gin.Context) {
	var req struct {
		Tele     string `json:"tele"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(err.Error()))
		return
	}

	resp, err := h.proxy.Login(c.Request.Context(), req.Tele, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusUnauthorized, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.User))
}

func (h *GatewayHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse("user id is required"))
		return
	}

	resp, err := h.proxy.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusNotFound, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.User))
}

func (h *GatewayHandler) GetFriends(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse("user id is required"))
		return
	}

	resp, err := h.proxy.GetFriends(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.Users))
}

func (h *GatewayHandler) AddFriend(c *gin.Context) {
	var req struct {
		FriendID string `json:"friend_id"`
		Reason   string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")

	resp, err := h.proxy.AddFriend(c.Request.Context(), userID.(string), req.FriendID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(nil))
}

func (h *GatewayHandler) AcceptFriendRequest(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")

	resp, err := h.proxy.AcceptFriendRequest(c.Request.Context(), req.RequestID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(nil))
}

func (h *GatewayHandler) GetPendingFriendRequests(c *gin.Context) {
	userID, _ := c.Get("user_id")
	_ = userID

	c.JSON(http.StatusOK, proxy.SuccessResponse([]interface{}{}))
}

func (h *GatewayHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")

	resp, err := h.proxy.CreateGroup(c.Request.Context(), userID.(string), req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.Group))
}

func (h *GatewayHandler) GetGroup(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse("group id is required"))
		return
	}

	resp, err := h.proxy.GetGroup(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusNotFound, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.Group))
}

func (h *GatewayHandler) GetGroupMembers(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse("group id is required"))
		return
	}

	resp, err := h.proxy.GetGroupMembers(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.Members))
}

func (h *GatewayHandler) GetUserGroups(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse("user id is required"))
		return
	}

	_ = userID
	c.JSON(http.StatusOK, proxy.SuccessResponse([]interface{}{}))
}

func (h *GatewayHandler) JoinGroup(c *gin.Context) {
	var req struct {
		GroupID string `json:"group_id"`
		Reason  string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")

	resp, err := h.proxy.JoinGroup(c.Request.Context(), userID.(string), req.GroupID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(nil))
}

func (h *GatewayHandler) AcceptGroupJoinRequest(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(err.Error()))
		return
	}

	_ = req.RequestID
	c.JSON(http.StatusOK, proxy.SuccessResponse(nil))
}

func (h *GatewayHandler) LeaveGroup(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse("group id is required"))
		return
	}

	userID, _ := c.Get("user_id")

	resp, err := h.proxy.LeaveGroup(c.Request.Context(), userID.(string), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(nil))
}

func (h *GatewayHandler) SendMessage(c *gin.Context) {
	var req struct {
		ReceiverID string `json:"receiver_id"`
		Content    string `json:"content"`
		Type       string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")

	resp, err := h.proxy.SendMessage(c.Request.Context(), userID.(string), req.ReceiverID, req.Content, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.MessageInfo))
}

func (h *GatewayHandler) GetMessage(c *gin.Context) {
	messageID := c.Param("id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse("message id is required"))
		return
	}

	_ = messageID
	c.JSON(http.StatusOK, proxy.SuccessResponse(nil))
}

func (h *GatewayHandler) GetOfflineMessages(c *gin.Context) {
	userID, _ := c.Get("user_id")

	resp, err := h.proxy.GetOfflineMessages(c.Request.Context(), userID.(string), 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(resp.Messages))
}

func (h *GatewayHandler) MarkAsRead(c *gin.Context) {
	messageID := c.Param("id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse("message id is required"))
		return
	}

	resp, err := h.proxy.MarkAsRead(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(nil))
}

func (h *GatewayHandler) MarkAllAsRead(c *gin.Context) {
	c.JSON(http.StatusOK, proxy.SuccessResponse(nil))
}

func (h *GatewayHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	resp, err := h.proxy.GetUnreadCount(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(err.Error()))
		return
	}

	if !resp.Success {
		c.JSON(http.StatusInternalServerError, proxy.ErrorResponse(resp.Message))
		return
	}

	c.JSON(http.StatusOK, proxy.SuccessResponse(gin.H{"unread_count": resp.UnreadCount}))
}

func (h *GatewayHandler) WebSocket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "websocket endpoint - connect to message service directly"})
}
