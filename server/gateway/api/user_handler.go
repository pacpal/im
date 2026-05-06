package api

import (
	auth "IM/server/gateway/auth"
	userService "IM/server/userservice"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSv  *userService.UserService
	groupSv *userService.GroupService
}

func NewUserHandler(userSv userService.UserService) *UserHandler {
	return &UserHandler{userSv: &userSv}
}
func (h *UserHandler) register(c *gin.Context) {
	var req struct {
		tele     string
		password string
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.tele == "" || req.password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong register info"})
		return
	}
	resp, err := h.userSv.Register(context.Background(), req.tele, req.password)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp.Name})
}
func (h *UserHandler) login(c *gin.Context) {
	var req struct {
		name     string
		tele     string
		id       string
		password string
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.id == "" && req.tele == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknow user"})
		return
	}
	err := h.userSv.Login(
		context.Background(),
		req.id,
		req.tele,
		req.password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	token, err := auth.GenerateToken(req.name, req.id)
	if err != nil {
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, 24*60*60, "/ws", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "login success", "name": req.name})
}
func (h *UserHandler) GetFriends(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
		return
	}
	resp, err := h.userSv.GetFriends(
		context.Background(),
		uid.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"friends": resp})
}
func (h *UserHandler) AddFriend(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		TargetID string `json:"target_id"`
		Reason   string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}
	err := h.userSv.AddFriend(
		context.Background(),
		userID.(string),
		req.TargetID,
		req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *UserHandler) RemoveFriend(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		TargetID string `json:"target_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}

	err := h.userSv.RemoveFriend(
		context.Background(),
		userID.(string),
		req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *UserHandler) ReplyFriendAdd(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		TargetID string `json:"target_id"`
		Reply    string `json:"reply"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}

	err := h.userSv.ReplyFriendAdd(
		context.Background(),
		userID.(string),
		req.TargetID,
		req.Reply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": err})
}

func (h *UserHandler) CreateGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}

	group, err := h.groupSv.CreateGroup(context.Background(), userID.(string), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"group_id": group.ID, "group_name": group.Name})
}

func (h *UserHandler) JoinGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		GroupID string `json:"group_id"`
		Reason  string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}

	err := h.groupSv.JoinGroup(context.Background(), userID.(string), req.GroupID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "already post"})
}

func (h *UserHandler) LeaveGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		GroupID string `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}
	err := h.groupSv.LeaveGroup(
		context.Background(),
		userID.(string),
		req.GroupID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "already leave"})
}

func (h *UserHandler) ReplyGroupAdd(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		UserID  string `json:"user_id"`
		GroupID string `json:"group_id"`
		Reply   string `json:"reply"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}

	err := h.groupSv.ReplyGroupAdd(
		context.Background(),
		userID.(string),
		req.UserID,
		req.GroupID,
		req.Reply,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "already post"})
}
