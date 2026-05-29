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

type RegisterRequest struct {
	Username string `json:"username" example:"testuser"`
	Password string `json:"password" example:"password123"`
	Email    string `json:"email" example:"user@example.com"`
	Nickname string `json:"nickname" example:"Test User"`
}

type RegisterResponse struct {
	UserId   string `json:"user_id" example:"123456789"`
	Username string `json:"username" example:"testuser"`
	Email    string `json:"email" example:"user@example.com"`
	Nickname string `json:"nickname" example:"Test User"`
}

type LoginRequest struct {
	Username string `json:"username" example:"testuser"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserId   string `json:"user_id" example:"123456789"`
	Username string `json:"username" example:"testuser"`
}

type UserInfo struct {
	UserId   string `json:"user_id" example:"123456789"`
	Username string `json:"username" example:"testuser"`
	Email    string `json:"email" example:"user@example.com"`
	Nickname string `json:"nickname" example:"Test User"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Status   int32  `json:"status" example:"1"`
}

type FriendInfo struct {
	FriendId string `json:"friend_id" example:"123456789"`
	Username string `json:"username" example:"frienduser"`
	Nickname string `json:"nickname" example:"My Friend"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Status   int32  `json:"status" example:"1"`
}

type FriendsListResponse struct {
	Friends []FriendInfo `json:"friends"`
}

type AddFriendRequest struct {
	FriendID string `json:"friend_id" example:"123456"`
	Reason   string `json:"reason" example:"想加你为好友"`
}

type AddFriendResponse struct {
	RequestId string `json:"request_id" example:"987654321"`
	Status    string `json:"status" example:"pending"`
}

type RemoveFriendResponse struct {
	Success bool `json:"success" example:"true"`
}

type ReplyFriendRequestBody struct {
	Accept bool `json:"accept" example:"true"`
}

type ReplyFriendResponse struct {
	Success bool   `json:"success" example:"true"`
	Status  string `json:"status" example:"accepted"`
}

type FriendRequestInfo struct {
	RequestId  string `json:"request_id" example:"987654321"`
	FromUserId string `json:"from_user_id" example:"123456789"`
	FromName   string `json:"from_name" example:"sender"`
	Reason     string `json:"reason" example:"想加你为好友"`
	Status     string `json:"status" example:"pending"`
	CreatedAt  string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type FriendRequestsResponse struct {
	Requests []FriendRequestInfo `json:"requests"`
}

type CreateGroupRequest struct {
	Name        string `json:"name" example:"技术交流群"`
	Description string `json:"description" example:"这是一个技术交流群"`
}

type CreateGroupResponse struct {
	GroupId     string `json:"group_id" example:"789012"`
	Name        string `json:"name" example:"技术交流群"`
	Description string `json:"description" example:"这是一个技术交流群"`
	OwnerId     string `json:"owner_id" example:"123456789"`
}

type GroupInfo struct {
	GroupId     string `json:"group_id" example:"789012"`
	Name        string `json:"name" example:"技术交流群"`
	Description string `json:"description" example:"这是一个技术交流群"`
	OwnerId     string `json:"owner_id" example:"123456789"`
	MemberCount int32  `json:"member_count" example:"50"`
	CreatedAt   string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type GroupMemberInfo struct {
	MemberId string `json:"member_id" example:"123456789"`
	Username string `json:"username" example:"member1"`
	Nickname string `json:"nickname" example:"群成员1"`
	Role     string `json:"role" example:"member"`
	JoinedAt string `json:"joined_at" example:"2024-01-01T00:00:00Z"`
}

type GroupMembersResponse struct {
	Members []GroupMemberInfo `json:"members"`
}

type ChangeGroupMemberRequest struct {
	Role int32 `json:"role" example:"1"`
}

type ChangeGroupMemberResponse struct {
	Success bool   `json:"success" example:"true"`
	Role    string `json:"role" example:"admin"`
}

type RemoveGroupMemberResponse struct {
	Success bool `json:"success" example:"true"`
}

type UserGroupsResponse struct {
	Groups []GroupInfo `json:"groups"`
}

type JoinGroupRequest struct {
	GroupID string `json:"group_id" example:"789012"`
	Reason  string `json:"reason" example:"想加入群组"`
}

type JoinGroupResponse struct {
	RequestId string `json:"request_id" example:"987654321"`
	Status    string `json:"status" example:"pending"`
}

type GetPendingGroupJoinRequestsBody struct {
	GroupID string `json:"group_id" example:"789012"`
}

type GroupJoinRequestInfo struct {
	RequestId string `json:"request_id" example:"987654321"`
	UserId    string `json:"user_id" example:"123456789"`
	Username  string `json:"username" example:"applicant"`
	Reason    string `json:"reason" example:"想加入群组"`
	Status    string `json:"status" example:"pending"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type GroupJoinRequestsResponse struct {
	Requests []GroupJoinRequestInfo `json:"requests"`
}

type ReplyGroupJoinRequestBody struct {
	Accept bool `json:"accept" example:"true"`
}

type ReplyGroupJoinResponse struct {
	Success bool   `json:"success" example:"true"`
	Status  string `json:"status" example:"accepted"`
}

type LeaveGroupResponse struct {
	Success bool `json:"success" example:"true"`
}

type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id" example:"123456"`
	Content    string `json:"content" example:"你好"`
	MsgType    string `json:"msg_type" example:"text"`
}

type SendMessageResponse struct {
	MsgId      string `json:"msg_id" example:"msg123456"`
	SenderId   string `json:"sender_id" example:"123456789"`
	ReceiverId string `json:"receiver_id" example:"987654321"`
	Content    string `json:"content" example:"你好"`
	MsgType    string `json:"msg_type" example:"text"`
	CreatedAt  string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type OfflineMessageInfo struct {
	MsgId      string `json:"msg_id" example:"msg123456"`
	SenderId   string `json:"sender_id" example:"123456789"`
	SenderName string `json:"sender_name" example:"sender"`
	Content    string `json:"content" example:"你好"`
	MsgType    string `json:"msg_type" example:"text"`
	CreatedAt  string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type OfflineMessagesResponse struct {
	Messages []OfflineMessageInfo `json:"messages"`
}

type MarkAsReadResponse struct {
	Success bool `json:"success" example:"true"`
}

type MarkAllAsReadResponse struct {
	Success bool  `json:"success" example:"true"`
	Count   int32 `json:"count" example:"10"`
}

type UnreadCountResponse struct {
	Count int32 `json:"count" example:"5"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

type HealthResponse struct {
	Status      string `json:"status" example:"healthy"`
	Service     string `json:"service" example:"gateway"`
	OnlineUsers int    `json:"online_users" example:"100"`
}

// Register 处理用户注册请求并转发到用户服务。
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册请求参数"
// @Success 200 {object} RegisterResponse "注册成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /register [post]
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
// @Summary 用户登录
// @Description 用户登录获取 JWT Token
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求参数"
// @Success 200 {object} LoginResponse "登录成功，返回 Token"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /login [post]
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
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "用户ID"
// @Success 200 {object} UserInfo "用户信息"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /users/{id} [get]
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
// @Summary 获取好友列表
// @Description 获取当前用户的所有好友
// @Tags 好友
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} FriendsListResponse "好友列表"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /friends [get]
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
// @Summary 添加好友
// @Description 发送好友请求给指定用户
// @Tags 好友
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddFriendRequest true "添加好友请求参数"
// @Success 200 {object} AddFriendResponse "好友请求已发送"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /friends [post]
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

// RemoveFriend 删除好友。
// @Summary 删除好友
// @Description 从好友列表中删除指定好友
// @Tags 好友
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param friend_id path string true "好友ID"
// @Success 200 {object} RemoveFriendResponse "删除成功"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /friends/{friend_id} [delete]
func RemoveFriend(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		friendID := c.Param("friend_id")
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.UserClient().RemoveFriend(ctx, &user.RemoveFriendRequest{
			UserId:   userID.(string),
			TargetId: friendID,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// ReplyFriendRequest 接受或拒绝好友请求。
// @Summary 回复好友请求
// @Description 接受或拒绝好友请求
// @Tags 好友
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request_id path string true "好友请求ID"
// @Param request body ReplyFriendRequestBody true "回复参数"
// @Success 200 {object} ReplyFriendResponse "回复成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /friend_requests/{request_id} [put]
func ReplyFriendRequest(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Param("request_id")
		var req struct {
			Accept bool `json:"accept"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.UserClient().ReplyFriend(ctx, &user.ReplyFriendRequest{
			UserId:    userID.(string),
			RequestId: requestID,
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
// @Summary 获取待处理好友请求
// @Description 获取当前用户收到的所有待处理好友请求
// @Tags 好友
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} FriendRequestsResponse "好友请求列表"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /friend_requests [get]
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
// @Summary 创建群组
// @Description 创建新的群组
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateGroupRequest true "创建群组请求参数"
// @Success 200 {object} CreateGroupResponse "群组创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /groups [post]
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
// @Summary 获取群组信息
// @Description 根据群组ID获取群组详细信息
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "群组ID"
// @Success 200 {object} GroupInfo "群组信息"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /groups/{id} [get]
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
// @Summary 获取群组成员列表
// @Description 获取指定群组的所有成员
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "群组ID"
// @Success 200 {object} GroupMembersResponse "群组成员列表"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /groups/{id}/members [get]
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

// ChangeGroupMember 修改群组成员角色。
// @Summary 修改群组成员角色
// @Description 群主修改群组成员的角色（管理员/普通成员）
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "群组ID"
// @Param member_id path string true "成员ID"
// @Param request body ChangeGroupMemberRequest true "角色参数"
// @Success 200 {object} ChangeGroupMemberResponse "修改成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /group/{id}/members/{member_id} [post]
func ChangeGroupMember(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Role int32 `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ownerID, _ := c.Get("user_id")
		groupID := c.Param("id")
		memberID := c.Param("member_id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().ChangeMember(ctx, &group.ChangeMemberRequest{
			GroupId:  groupID,
			OwnerId:  ownerID.(string),
			MemberId: memberID,
			Role:     group.MemberRole(req.Role),
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// RemoveGroupMember 移除群组成员。
// @Summary 移除群组成员
// @Description 群主或管理员移除群组成员
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "群组ID"
// @Param member_id path string true "成员ID"
// @Success 200 {object} RemoveGroupMemberResponse "移除成功"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /groups/{id}/members/{member_id} [delete]
func RemoveGroupMember(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		memberID := c.Param("member_id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().RemoveMember(ctx, &group.RemoveMemberRequest{GroupId: groupID, MemberId: memberID})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetUserGroups 获取用户加入的群组列表。
// @Summary 获取用户群组列表
// @Description 获取当前用户加入的所有群组
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserGroupsResponse "群组列表"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /groups [get]
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
// @Summary 申请加入群组
// @Description 发送加入群组的申请请求
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body JoinGroupRequest true "加入群组请求参数"
// @Success 200 {object} JoinGroupResponse "申请已发送"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /group_requests [post]
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

// GetPendingGroupJoinRequests 获取待处理的群组加入请求。
// @Summary 获取待处理群组加入请求
// @Description 群主获取待处理的群组加入请求列表
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GetPendingGroupJoinRequestsBody true "群组ID"
// @Success 200 {object} GroupJoinRequestsResponse "加入请求列表"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /group_requests [get]
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
// @Summary 回复群组加入请求
// @Description 群主接受或拒绝加入群组的请求
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request_id path string true "请求ID"
// @Param request body ReplyGroupJoinRequestBody true "回复参数"
// @Success 200 {object} ReplyGroupJoinResponse "回复成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /group_requests/{request_id} [put]
func ReplyGroupJoinRequest(p *proxy.ServiceProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Param("request_id")
		var req struct {
			Accept bool `json:"accept"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()

		resp, err := p.GroupClient().ReplyGroupJoin(ctx, &group.ReplyGroupJoinRequest{
			OwnerId:   userID.(string),
			RequestId: requestID,
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
// @Summary 退出群组
// @Description 用户退出指定的群组
// @Tags 群组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "群组ID"
// @Success 200 {object} LeaveGroupResponse "退出成功"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /groups/{id}/members/me [delete]
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
// @Summary 发送消息
// @Description 发送消息给指定用户或群组
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SendMessageRequest true "发送消息请求参数"
// @Success 200 {object} SendMessageResponse "消息发送成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /messages [post]
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
// @Summary 获取离线消息
// @Description 获取当前用户的离线消息列表
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} OfflineMessagesResponse "离线消息列表"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /messages/offline [get]
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
// @Summary 标记消息已读
// @Description 将指定消息标记为已读状态
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "消息ID"
// @Success 200 {object} MarkAsReadResponse "标记成功"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /messages/{id}/read [put]
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
// @Summary 标记所有消息已读
// @Description 将当前用户的所有消息标记为已读状态
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MarkAllAsReadResponse "标记成功"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /messages/read [put]
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
// @Summary 获取未读消息数
// @Description 获取当前用户的未读消息数量
// @Tags 消息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UnreadCountResponse "未读消息数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /messages/unread/count [get]
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
