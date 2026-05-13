package entity

// User 用户实体
type User struct {
	ID        string
	Name      string
	Tele      string
	Password  string
	AvatarURL string
	Status    int
	CreatedAt int64
	UpdatedAt int64
}

// Friendship 好友关系
type Friendship struct {
	UserID    string
	FriendID  string
	Status    int
	CreatedAt int64
}

// FriendRequest 好友请求
type FriendRequest struct {
	ID        string
	FromUID   string
	ToUID     string
	Reason    string
	Status    string // pending, accepted, rejected
	CreatedAt int64
	UpdatedAt int64
}

const (
	FriendRequestPending  = "pending"
	FriendRequestAccepted = "accepted"
	FriendRequestRejected = "rejected"
)
