package handler

type RegisterRequest struct {
	Tele     string `json:"tele" binding:"required" example:"13800138000"`
	ID       string `json:"id" example:"user123"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Name     string `json:"name" binding:"required" example:"testuser"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type RegisterResponse struct {
	UserId string `json:"user_id" example:"123456789"`
	Name   string `json:"name" example:"testuser"`
}

type LoginRequest struct {
	Tele     string `json:"tele" example:"13800138000"`
	Password string `json:"password" binding:"required" example:"password123"`
	ID       string `json:"id" example:"user123"`
	Email    string `json:"email" example:"user@example.com"`
}

type LoginResponse struct {
	UserId string `json:"user_id" example:"123456789"`
	Name   string `json:"name" example:"testuser"`
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type UserInfo struct {
	ID        string `json:"id" example:"123456789"`
	Name      string `json:"name" example:"testuser"`
	Tele      string `json:"tele" example:"13800138000"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	Status    int32  `json:"status" example:"1"`
	CreatedAt int64  `json:"created_at" example:"1704067200"`
	UpdatedAt int64  `json:"updated_at" example:"1704067200"`
}

type FriendInfo struct {
	ID        string `json:"id" example:"123456789"`
	Name      string `json:"name" example:"frienduser"`
	Tele      string `json:"tele" example:"13800138000"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	Status    int32  `json:"status" example:"1"`
}

type FriendsListResponse struct {
	Friends []FriendInfo `json:"friends"`
}

type AddFriendRequest struct {
	FriendID string `json:"friend_id" binding:"required" example:"123456"`
	Reason   string `json:"reason" example:"想加你为好友"`
}

type CommonResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"操作成功"`
	Code    int32  `json:"code" example:"0"`
}

type ReplyFriendRequestBody struct {
	Accept bool `json:"accept" example:"true"`
}

type FriendRequestInfo struct {
	ID        string    `json:"id" example:"987654321"`
	FromUID   string    `json:"from_uid" example:"123456789"`
	ToUID     string    `json:"to_uid" example:"987654321"`
	Reason    string    `json:"reason" example:"想加你为好友"`
	Status    string    `json:"status" example:"pending"`
	CreatedAt int64     `json:"created_at" example:"1704067200"`
	FromUser  *UserInfo `json:"from_user,omitempty"`
}

type FriendRequestsResponse struct {
	Requests []FriendRequestInfo `json:"requests"`
}

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required" example:"技术交流群"`
	Description string `json:"description" example:"这是一个技术交流群"`
	ImageURL    string `json:"image_url" example:"https://example.com/group.jpg"`
}

type CreateGroupResponse struct {
	GroupId string `json:"group_id" example:"789012"`
	Name    string `json:"name" example:"技术交流群"`
}

type GroupInfo struct {
	ID          string `json:"id" example:"789012"`
	Name        string `json:"name" example:"技术交流群"`
	Description string `json:"description" example:"这是一个技术交流群"`
	OwnerID     string `json:"owner_id" example:"123456789"`
	Type        string `json:"type" example:"public"`
	ImageURL    string `json:"image_url" example:"https://example.com/group.jpg"`
	CreatedAt   int64  `json:"created_at" example:"1704067200"`
	UpdatedAt   int64  `json:"updated_at" example:"1704067200"`
	MemberCount int32  `json:"member_count" example:"50"`
}

type MemberInfo struct {
	UserID    string `json:"user_id" example:"123456789"`
	Name      string `json:"name" example:"member1"`
	Nickname  string `json:"nickname" example:"群成员1"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	Role      string `json:"role" example:"member"`
	JoinedAt  int64  `json:"joined_at" example:"1704067200"`
}

type GroupMembersResponse struct {
	Members []MemberInfo `json:"members"`
}

type ChangeGroupMemberRequest struct {
	Role int32 `json:"role" binding:"required" example:"1"`
}

type JoinGroupRequest struct {
	GroupID string `json:"group_id" binding:"required" example:"789012"`
	Reason  string `json:"reason" example:"想加入群组"`
}

type GetPendingGroupJoinRequestsBody struct {
	GroupID string `json:"group_id" binding:"required" example:"789012"`
}

type GroupJoinRequestInfo struct {
	ID         string `json:"id" example:"987654321"`
	UserID     string `json:"user_id" example:"123456789"`
	GroupID    string `json:"group_id" example:"789012"`
	Reason     string `json:"reason" example:"想加入群组"`
	Status     string `json:"status" example:"pending"`
	CreatedAt  int64  `json:"created_at" example:"1704067200"`
	UserName   string `json:"user_name" example:"applicant"`
	UserAvatar string `json:"user_avatar" example:"https://example.com/avatar.jpg"`
}

type GroupJoinRequestsResponse struct {
	Requests []GroupJoinRequestInfo `json:"requests"`
}

type ReplyGroupJoinRequestBody struct {
	Accept bool `json:"accept" example:"true"`
}

type UserGroupsResponse struct {
	Groups []GroupInfo `json:"groups"`
}

type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id" binding:"required" example:"123456"`
	Content    string `json:"content" binding:"required" example:"你好"`
	MsgType    string `json:"msg_type" binding:"required" example:"text"`
	Timestamp  int64  `json:"timestamp" example:"1704067200"`
}

type SendMessageResponse struct {
	MsgId     string `json:"msg_id" example:"msg123456"`
	Success   bool   `json:"success" example:"true"`
	Timestamp int64  `json:"timestamp" example:"1704067200"`
}

type MessageInfo struct {
	MsgId      string `json:"msg_id" example:"msg123456"`
	SenderID   string `json:"sender_id" example:"123456789"`
	ReceiverID string `json:"receiver_id" example:"987654321"`
	Content    string `json:"content" example:"你好"`
	MsgType    string `json:"msg_type" example:"text"`
	Timestamp  int64  `json:"timestamp" example:"1704067200"`
	IsRead     bool   `json:"is_read" example:"false"`
	IsRevoked  bool   `json:"is_revoked" example:"false"`
	ReadAt     int64  `json:"read_at" example:"1704067200"`
}

type OfflineMessagesResponse struct {
	Messages []MessageInfo `json:"messages"`
	Total    int32         `json:"total"`
}

type UnreadCountResponse struct {
	Count       int64            `json:"count" example:"5"`
	GroupCounts map[string]int64 `json:"group_counts"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

type HealthResponse struct {
	Status      string `json:"status" example:"healthy"`
	Service     string `json:"service" example:"gateway"`
	OnlineUsers int    `json:"online_users" example:"100"`
}