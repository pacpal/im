// Package persistence 提供基于 Redis 的缓存实现。
package persistence

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// MessageCache 封装 Redis 操作，用于在线用户集合与未读计数。
type MessageCache struct {
	client *redis.Client
}

func NewMessageCache(client *redis.Client) *MessageCache {
	return &MessageCache{client: client}
}

func (c *MessageCache) onlineUsersKey() string {
	return "online_users"
}

func (c *MessageCache) unreadCountKey(userID string) string {
	return "unread:" + userID
}

// GetOnlineUsers 返回当前在线用户集合（map[string]bool 形式）。
func (c *MessageCache) GetOnlineUsers(ctx context.Context) (map[string]bool, error) {
	members, err := c.client.SMembers(ctx, c.onlineUsersKey()).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, m := range members {
		result[m] = true
	}
	return result, nil
}

// SetOnlineUser 将用户加入在线集合。
func (c *MessageCache) SetOnlineUser(ctx context.Context, userID string) error {
	return c.client.SAdd(ctx, c.onlineUsersKey(), userID).Err()
}

// RemoveOnlineUser 从在线集合移除用户。
func (c *MessageCache) RemoveOnlineUser(ctx context.Context, userID string) error {
	return c.client.SRem(ctx, c.onlineUsersKey(), userID).Err()
}

// GetUnreadCount 返回用户未读计数（直接读取 key）。
func (c *MessageCache) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return c.client.Get(ctx, c.unreadCountKey(userID)).Int64()
}

// IncrUnreadCount 未读计数自增。
func (c *MessageCache) IncrUnreadCount(ctx context.Context, userID string) error {
	return c.client.Incr(ctx, c.unreadCountKey(userID)).Err()
}

// DecrUnreadCount 未读计数自减。
func (c *MessageCache) DecrUnreadCount(ctx context.Context, userID string) error {
	return c.client.Decr(ctx, c.unreadCountKey(userID)).Err()
}
