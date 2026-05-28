// Package handler 提供 Gateway 的 HTTP API 处理函数，主要负责参数校验、鉴权转发以及调用后端 gRPC 服务。
package handler

import (
	"IM/api/gen/group"
	"IM/api/gen/message"
	"IM/api/gen/user"
	"IM/services/gateway/proxy"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 处理用户注册请求并转发到用户服务。
func Register(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()

		resp, err := p.UserClient().Register(ctx, &req)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// Login 处理用户登录请求并转发到用户服务。
func Login(p *proxy.ServiceProxy, jwtUtil interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()

		resp, err := p.UserClient().Login(ctx, &req)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// GetUser 获取指定用户的信息。
func GetUser(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		ctx := c.Request.Context()

		resp, err := p.UserClient().GetUser(ctx, &user.GetUserRequest{UserId: userID})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetFriends 获取用户好友列表。
func GetFriends(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.UserClient().GetFriends(ctx, &user.GetFriendsRequest{UserId: userID.(string)})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// AddFriend 发送好友请求（转发到用户服务）。
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
		ctx := c.Request.Context()

		resp, err := p.UserClient().AddFriend(ctx, &user.AddFriendRequest{
			UserId:   userID.(string),
			TargetId: req.FriendID,
			Reason:   req.Reason,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func RemoveFriend(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FriendID string `json:"friend_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.UserClient().RemoveFriend(ctx, &user.RemoveFriendRequest{
			UserId:   userID.(string),
			TargetId: req.FriendID,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// AcceptFriendRequest 接受或拒绝好友请求。
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
		ctx := c.Request.Context()

		resp, err := p.UserClient().ReplyFriend(ctx, &user.ReplyFriendRequest{
			UserId:    userID.(string),
			RequestId: req.RequestID,
			Accept:    req.Accept,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetPendingFriendRequests 获取当前用户的待处理好友请求。
func GetPendingFriendRequests(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.UserClient().GetPendingFriendRequests(ctx, &user.GetPendingFriendRequestsRequest{UserId: userID.(string)})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// CreateGroup 创建群组。
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
		ctx := c.Request.Context()

		resp, err := p.GroupClient().CreateGroup(ctx, &group.CreateGroupRequest{
			OwnerId:     userID.(string),
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetGroup 获取群组信息。
func GetGroup(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().GetGroup(ctx, &group.GetGroupRequest{GroupId: groupID})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetGroupMembers 获取群组成员列表。
func GetGroupMembers(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().GetMembers(ctx, &group.GetMembersRequest{GroupId: groupID})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
func RemoveGroupMember(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().RemoveMember(ctx, &group.RemoveMemberRequest{GroupId: groupID})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetUserGroups 获取用户加入的群组列表。
func GetUserGroups(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().GetUserGroups(ctx, &group.GetUserGroupsRequest{UserId: userID.(string)})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// JoinGroup 发送加入群组请求。
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
		ctx := c.Request.Context()

		resp, err := p.GroupClient().JoinGroup(ctx, &group.JoinGroupRequest{
			UserId:  userID.(string),
			GroupId: req.GroupID,
			Reason:  req.Reason,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
func GetPendingGroupJoinRequests(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			GroupID string `json:"group_id"`
		}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().GetPendingGroupRequests(ctx, &group.GetPendingGroupRequestsRequest{
			OwnerId: userID.(string),
			GroupId: req.GroupID,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// ReplyGroupJoinRequest 群主接受或拒绝加入群组的请求。
func ReplyGroupJoinRequest(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			GroupID string `json:"group_id"`
			Accept  bool   `json:"accept"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().ReplyGroupJoin(ctx, &group.ReplyGroupJoinRequest{
			OwnerId:   userID.(string),
			RequestId: req.GroupID,
			Accept:    req.Accept,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// LeaveGroup 退出群组。
func LeaveGroup(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		userID, _ := c.Get("user_id")

		ctx := c.Request.Context()

		resp, err := p.GroupClient().LeaveGroup(ctx, &group.LeaveGroupRequest{
			UserId:  userID.(string),
			GroupId: groupID,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// SendMessage 发送消息（转发到消息服务）。
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
		ctx := c.Request.Context()

		resp, err := p.MessageClient().SendMessage(ctx, &message.SendMessageRequest{
			SenderId:   userID.(string),
			ReceiverId: req.ReceiverID,
			Content:    req.Content,
			MsgType:    req.MsgType,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetOfflineMessages 获取离线消息列表。
func GetOfflineMessages(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.MessageClient().GetOfflineMessages(ctx, &message.GetOfflineMessagesRequest{
			UserId: userID.(string),
			Limit:  50,
			Offset: 0,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// MarkAsRead 标记指定消息为已读。
func MarkAsRead(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := c.Param("id")

		ctx := c.Request.Context()

		resp, err := p.MessageClient().MarkAsRead(ctx, &message.MarkAsReadRequest{
			MsgId: msgID,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// MarkAllAsRead 将当前用户的所有消息标记为已读。
func MarkAllAsRead(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.MessageClient().MarkAllAsRead(ctx, &message.MarkAllAsReadRequest{
			UserId: userID.(string),
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetUnreadCount 获取当前用户未读消息数。
func GetUnreadCount(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.MessageClient().GetUnreadCount(ctx, &message.GetUnreadCountRequest{UserId: userID.(string)})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
