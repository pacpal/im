package cache

import (
	"IM/services/user/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewUserCache(client *redis.Client) *UserCache {
	return &UserCache{
		client: client,
		ttl:    time.Hour,
	}
}

func (c *UserCache) userKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

func (c *UserCache) Get(ctx context.Context, userID string) (*entity.User, error) {
	key := c.userKey(userID)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *UserCache) Set(ctx context.Context, user *entity.User) error {
	key := c.userKey(user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *UserCache) Delete(ctx context.Context, userID string) error {
	key := c.userKey(userID)
	return c.client.Del(ctx, key).Err()
}

func (c *UserCache) Exists(ctx context.Context, userID string) (bool, error) {
	key := c.userKey(userID)
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
