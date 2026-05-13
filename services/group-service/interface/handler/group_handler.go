package handler

import (
	"IM/services/group-service/application/service"
	"IM/services/group-service/interface/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GroupHttpHandler struct {
	groupSvc *service.GroupApplicationService
}

func NewGroupHttpHandler(groupSvc *service.GroupApplicationService) *GroupHttpHandler {
	return &GroupHttpHandler{groupSvc: groupSvc}
}

func (h *GroupHttpHandler) CreateGroup(c *gin.Context) {
	var req dto.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ownerID, _ := c.Get("user_id")

	group, err := h.groupSvc.CreateGroup(c.Request.Context(), ownerID.(string), req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "group created",
		"group":   dto.ToGroupResponse(group),
	})
}

func (h *GroupHttpHandler) GetGroup(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "group id is required"})
		return
	}

	group, err := h.groupSvc.GetGroup(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"group":   dto.ToGroupResponse(group),
	})
}

func (h *GroupHttpHandler) GetGroupMembers(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "group id is required"})
		return
	}

	members, err := h.groupSvc.GetGroupMembers(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"members": dto.ToGroupMemberResponseList(members),
	})
}

func (h *GroupHttpHandler) GetUserGroups(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "user id is required"})
		return
	}

	groups, err := h.groupSvc.GetUserGroups(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"groups":  dto.ToGroupResponseList(groups),
	})
}

func (h *GroupHttpHandler) JoinGroup(c *gin.Context) {
	var req dto.JoinGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err := h.groupSvc.JoinGroup(c.Request.Context(), userID.(string), req.GroupID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "join request sent",
	})
}

func (h *GroupHttpHandler) AcceptJoinRequest(c *gin.Context) {
	var req dto.JoinRequestActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err := h.groupSvc.AcceptJoinRequest(c.Request.Context(), req.RequestID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "join request accepted",
	})
}

func (h *GroupHttpHandler) LeaveGroup(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "group id is required"})
		return
	}

	userID, _ := c.Get("user_id")

	err := h.groupSvc.LeaveGroup(c.Request.Context(), userID.(string), groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "left group",
	})
}

func (h *GroupHttpHandler) RegisterRoutes(router *gin.RouterGroup, groupSvc *service.GroupApplicationService) {
	h.groupSvc = groupSvc

	groupGroup := router.Group("/groups")
	{
		groupGroup.POST("", h.CreateGroup)
		groupGroup.GET("/:id", h.GetGroup)
		groupGroup.GET("/:id/members", h.GetGroupMembers)
		groupGroup.GET("/users/:id/groups", h.GetUserGroups)
		groupGroup.POST("/join", h.JoinGroup)
		groupGroup.POST("/join/accept", h.AcceptJoinRequest)
		groupGroup.DELETE("/:id/leave", h.LeaveGroup)
	}
}