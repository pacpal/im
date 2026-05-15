package persistence

import (
	"context"

	"github.com/redis/go-redis/v9"
)

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

func (c *MessageCache) SetOnlineUser(ctx context.Context, userID string) error {
	return c.client.SAdd(ctx, c.onlineUsersKey(), userID).Err()
}

func (c *MessageCache) RemoveOnlineUser(ctx context.Context, userID string) error {
	return c.client.SRem(ctx, c.onlineUsersKey(), userID).Err()
}

func (c *MessageCache) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return c.client.Get(ctx, c.unreadCountKey(userID)).Int64()
}

func (c *MessageCache) IncrUnreadCount(ctx context.Context, userID string) error {
	return c.client.Incr(ctx, c.unreadCountKey(userID)).Err()
}

func (c *MessageCache) DecrUnreadCount(ctx context.Context, userID string) error {
	return c.client.Decr(ctx, c.unreadCountKey(userID)).Err()
}
