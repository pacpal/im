package handler

import (
	"IM/services/user-service/application/service"
	"IM/services/user-service/interface/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHttpHandler struct {
	userSvc *service.UserApplicationService
}

func NewUserHttpHandler(userSvc *service.UserApplicationService) *UserHttpHandler {
	return &UserHttpHandler{userSvc: userSvc}
}

func (h *UserHttpHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := h.userSvc.Register(c.Request.Context(), req.Tele, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "register success",
		"user":    dto.ToUserResponse(user),
	})
}

func (h *UserHttpHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := h.userSvc.Login(c.Request.Context(), req.Tele, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login success",
		"user":    dto.ToUserResponse(user),
	})
}

func (h *UserHttpHandler) GetUser(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "user id is required"})
		return
	}

	user, err := h.userSvc.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    dto.ToUserResponse(user),
	})
}

func (h *UserHttpHandler) GetFriends(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "user id is required"})
		return
	}

	friends, err := h.userSvc.GetFriends(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"friends": dto.ToUserResponseList(friends),
	})
}

func (h *UserHttpHandler) AddFriend(c *gin.Context) {
	var req dto.AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err := h.userSvc.AddFriend(c.Request.Context(), userID.(string), req.ToUID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "friend request sent",
	})
}

func (h *UserHttpHandler) AcceptFriendRequest(c *gin.Context) {
	var req dto.FriendRequestActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err := h.userSvc.AcceptFriendRequest(c.Request.Context(), req.RequestID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "friend request accepted",
	})
}

func (h *UserHttpHandler) GetPendingFriendRequests(c *gin.Context) {
	userID, _ := c.Get("user_id")

	requests, err := h.userSvc.GetPendingFriendRequests(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"requests": dto.ToFriendRequestResponseList(requests),
	})
}

func (h *UserHttpHandler) RegisterRoutes(router *gin.RouterGroup, userSvc *service.UserApplicationService) {
	h.userSvc = userSvc

	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", h.Register)
		userGroup.POST("/login", h.Login)
		userGroup.GET("/:id", h.GetUser)
		userGroup.GET("/:id/friends", h.GetFriends)
		userGroup.POST("/friends", h.AddFriend)
		userGroup.POST("/friend-requests/accept", h.AcceptFriendRequest)
		userGroup.GET("/friend-requests", h.GetPendingFriendRequests)
	}
}