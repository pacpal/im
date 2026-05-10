package api

import (
	auth "IM/server/gateway/auth"
	userService "IM/server/userservice"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSv  *userservice.UserService
	groupSv *userservice.GroupService
	userSv  *userService.UserService
	groupSv *userService.GroupService
func NewUserHandler(userSv *userservice.UserService, groupSv *userservice.GroupService) *UserHandler {
	return &UserHandler{userSv: userSv, groupSv: groupSv}
func NewUserHandler(userSv userService.UserService) *UserHandler {
	return &UserHandler{userSv: &userSv}
func (h *UserHandler) Register(c *gin.Context) {
func (h *UserHandler) register(c *gin.Context) {
		Name     string `json:"name"`
		tele     string
		password string
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong register info"})
	if err := c.ShouldBindJSON(&req); err != nil || req.tele == "" || req.password == "" {
	}
	resp, err := h.userSv.Register(context.Background(), req.Tele, req.Name, req.Password)
	if err != nil {
	resp, err := h.userSv.Register(context.Background(), req.tele, req.password)
		return
	c.JSON(http.StatusOK, gin.H{"message": resp.Name, "user_id": resp.ID})
}
	c.JSON(http.StatusOK, gin.H{"message": resp.Name})
func (h *UserHandler) Login(c *gin.Context) {
func (h *UserHandler) login(c *gin.Context) {
		Tele     string `json:"tele"`
		name     string
		tele     string
		id       string
		password string
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown user"})
	if err := c.ShouldBindJSON(&req); err != nil || req.id == "" && req.tele == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknow user"})
	user, err := h.userSv.Login(context.Background(), req.ID, req.Tele, req.Password)
	if err != nil {
	err := h.userSv.Login(
		context.Background(),
		req.id,
		req.tele,
		req.password)
		return
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	if err != nil {
	token, err := auth.GenerateToken(req.name, req.id)
		return
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, 24*60*60, "/ws", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "login success", "name": user.Name, "user_id": user.ID})
}
	c.JSON(http.StatusOK, gin.H{"message": "login success", "name": req.name})
func (h *UserHandler) GetFriends(c *gin.Context) {
	if !ok {
	uid, ok := c.Get("uid")
		return
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
	resp, err := h.userSv.GetFriends(context.Background(), uid.(string))
	if err != nil {
	resp, err := h.userSv.GetFriends(
		context.Background(),
		uid.(string))
		return
	}
}

func (h *UserHandler) AddFriend(c *gin.Context) {
	var req struct {
		TargetID string `json:"target_id"`
		Reason   string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}
	err := h.userSv.AddFriend(context.Background(), userID.(string), req.TargetID, req.Reason)
	if err != nil {
	err := h.userSv.AddFriend(
		context.Background(),
		userID.(string),
		req.TargetID,
		req.Reason)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "friend request sent"})
}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
func (h *UserHandler) RemoveFriend(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		TargetID string `json:"target_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}
	err := h.userSv.RemoveFriend(context.Background(), userID.(string), req.TargetID)
	if err != nil {

	err := h.userSv.RemoveFriend(
		context.Background(),
		userID.(string),
		req.TargetID)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "friend removed"})
}
func (h *UserHandler) ReplyFriendAdd(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		RequestID string `json:"request_id"`
		Reply     string `json:"reply"`
		TargetID string `json:"target_id"`
		Reply    string `json:"reply"`
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}
	err := h.userSv.ReplyFriendAdd(context.Background(), userID.(string), req.RequestID, req.Reply)
	if err != nil {

	err := h.userSv.ReplyFriendAdd(
		context.Background(),
		userID.(string),
		req.TargetID,
		req.Reply)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "reply processed"})
}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": err})
func (h *UserHandler) CreateGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Name string `json:"name"`
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}
	group, err := h.groupSv.CreateGroup(context.Background(), userID.(string), req.Name, req.Description)
	if err != nil {

	group, err := h.groupSv.CreateGroup(context.Background(), userID.(string), req.Name)
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
	c.JSON(http.StatusOK, gin.H{"message": "join request sent"})
}

	c.JSON(http.StatusOK, gin.H{"message": "already post"})
func (h *UserHandler) LeaveGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		GroupID string `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}
	err := h.groupSv.LeaveGroup(context.Background(), userID.(string), req.GroupID)
	if err != nil {
	err := h.groupSv.LeaveGroup(
		context.Background(),
		userID.(string),
		req.GroupID,
	)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "left group"})
}
	c.JSON(http.StatusOK, gin.H{"message": "already leave"})
func (h *UserHandler) ReplyGroupAdd(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		RequestID string `json:"request_id"`
		Reply     string `json:"reply"`
		UserID  string `json:"user_id"`
		GroupID string `json:"group_id"`
		Reply   string `json:"reply"`
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong info"})
		return
	}
	err := h.groupSv.ReplyGroupAdd(context.Background(), userID.(string), req.RequestID, req.Reply)
	if err != nil {

	err := h.groupSv.ReplyGroupAdd(
		context.Background(),
		userID.(string),
		req.UserID,
		req.GroupID,
		req.Reply,
	)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "reply processed"})
}
	c.JSON(http.StatusOK, gin.H{"message": "already post"})