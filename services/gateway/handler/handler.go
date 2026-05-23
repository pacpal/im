package handler

import (
	"IM/api/gen/group"
	"IM/api/gen/message"
	"IM/api/gen/user"
	"IM/services/gateway/proxy"
	"net/http"

	"IM/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func Register(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.UserClient().Register(ctx, &req)
		if err != nil {
			c.Error(err)
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

		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.UserClient().Login(ctx, &req)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func GetUser(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				if len(ts) > 8 {
					logger.Infof("Forwarding token prefix=%s len=%d", ts[:8], len(ts))
				} else {
					logger.Infof("Forwarding token len=%d", len(ts))
				}
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.UserClient().GetUser(ctx, &user.GetUserRequest{UserId: userID})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetFriends(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.UserClient().GetFriends(ctx, &user.GetFriendsRequest{UserId: userID})
		if err != nil {
			c.Error(err)
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
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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

func GetPendingFriendRequests(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.UserClient().GetPendingFriendRequests(ctx, &user.GetPendingFriendRequestsRequest{UserId: userID.(string)})
		if err != nil {
			c.Error(err)
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
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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

func GetGroup(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.GroupClient().GetGroup(ctx, &group.GetGroupRequest{GroupId: groupID})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetGroupMembers(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.GroupClient().GetMembers(ctx, &group.GetMembersRequest{GroupId: groupID})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetUserGroups(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.GroupClient().GetUserGroups(ctx, &group.GetUserGroupsRequest{UserId: userID})
		if err != nil {
			c.Error(err)
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
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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

func AcceptGroupJoinRequest(p *proxy.ServiceProxy) gin.HandlerFunc {
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
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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

func LeaveGroup(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		userID, _ := c.Get("user_id")

		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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

func GetOfflineMessages(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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

func MarkAsRead(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		msgID := c.Param("id")

		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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

func MarkAllAsRead(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

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

func GetUnreadCount(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()
		if t, ok := c.Get("token"); ok {
			if ts, ok2 := t.(string); ok2 {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+ts)
			}
		}

		resp, err := p.MessageClient().GetUnreadCount(ctx, &message.GetUnreadCountRequest{UserId: userID.(string)})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
