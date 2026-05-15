package handler

import (
	"IM/api/gen/group"
	"IM/api/gen/message"
	"IM/api/gen/user"
	"IM/services/gateway/proxy"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := p.UserClient().Register(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func Login(p *proxy.ServiceProxy, jwtUtil interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := p.UserClient().Login(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func GetUser(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		resp, err := p.UserClient().GetUser(c.Request.Context(), &user.GetUserRequest{UserId: userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetFriends(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		resp, err := p.UserClient().GetFriends(c.Request.Context(), &user.GetFriendsRequest{UserId: userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func AddFriend(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FriendID string `json:"friend_id"`
			Reason   string `json:"reason"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		resp, err := p.UserClient().AddFriend(c.Request.Context(), &user.AddFriendRequest{
			UserId:   userID.(string),
			TargetId: req.FriendID,
			Reason:   req.Reason,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func AcceptFriendRequest(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RequestID string `json:"request_id"`
			Accept    bool   `json:"accept"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		resp, err := p.UserClient().ReplyFriendRequest(c.Request.Context(), &user.ReplyFriendRequest{
			UserId:    userID.(string),
			RequestId: req.RequestID,
			Accept:    req.Accept,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetPendingFriendRequests(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		resp, err := p.UserClient().GetPendingFriendRequests(c.Request.Context(), &user.GetPendingFriendRequestsRequest{UserId: userID.(string)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func CreateGroup(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		resp, err := p.GroupClient().CreateGroup(c.Request.Context(), &group.CreateGroupRequest{
			OwnerId:     userID.(string),
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetGroup(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		resp, err := p.GroupClient().GetGroup(c.Request.Context(), &group.GetGroupRequest{GroupId: groupID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetGroupMembers(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		resp, err := p.GroupClient().GetMembers(c.Request.Context(), &group.GetMembersRequest{GroupId: groupID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetUserGroups(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		resp, err := p.GroupClient().GetUserGroups(c.Request.Context(), &group.GetUserGroupsRequest{UserId: userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func JoinGroup(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			GroupID string `json:"group_id"`
			Reason  string `json:"reason"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		resp, err := p.GroupClient().JoinGroup(c.Request.Context(), &group.JoinGroupRequest{
			UserId:  userID.(string),
			GroupId: req.GroupID,
			Reason:  req.Reason,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func AcceptGroupJoinRequest(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RequestID string `json:"request_id"`
			Accept    bool   `json:"accept"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		resp, err := p.GroupClient().ReplyGroupJoinRequest(c.Request.Context(), &group.ReplyGroupJoinRequest{
			OwnerId:   userID.(string),
			RequestId: req.RequestID,
			Accept:    req.Accept,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func LeaveGroup(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		userID, _ := c.Get("user_id")

		resp, err := p.GroupClient().LeaveGroup(c.Request.Context(), &group.LeaveGroupRequest{
			UserId:  userID.(string),
			GroupId: groupID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func SendMessage(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ReceiverID string `json:"receiver_id"`
			Content    string `json:"content"`
			MsgType    string `json:"msg_type"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		resp, err := p.MessageClient().SendMessage(c.Request.Context(), &message.SendMessageRequest{
			SenderId:   userID.(string),
			ReceiverId: req.ReceiverID,
			Content:    req.Content,
			MsgType:    req.MsgType,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetMessage(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := c.Param("id")
		resp, err := p.MessageClient().GetMessage(c.Request.Context(), &message.GetMessageRequest{MsgId: msgID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetOfflineMessages(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		resp, err := p.MessageClient().GetOfflineMessages(c.Request.Context(), &message.GetOfflineMessagesRequest{
			UserId: userID.(string),
			Limit:  50,
			Offset: 0,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func MarkAsRead(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := c.Param("id")
		userID, _ := c.Get("user_id")

		resp, err := p.MessageClient().MarkAsRead(c.Request.Context(), &message.MarkAsReadRequest{
			MsgId:  msgID,
			UserId: userID.(string),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func MarkAllAsRead(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TargetID string `json:"target_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		resp, err := p.MessageClient().MarkAllAsRead(c.Request.Context(), &message.MarkAllAsReadRequest{
			UserId:   userID.(string),
			TargetId: req.TargetID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetUnreadCount(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		resp, err := p.MessageClient().GetUnreadCount(c.Request.Context(), &message.GetUnreadCountRequest{UserId: userID.(string)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
