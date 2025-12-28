package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	conn *redis.Client
}

func NewRedisClient(host, port, password, db string) (*RedisClient, error) {
	uri := fmt.Sprintf("%s:%s", host, port)
	dbInt, err := strconv.Atoi(db)
	if err != nil {
		return nil, err
	}
	client := &RedisClient{
		conn: redis.NewClient(&redis.Options{
			Addr:     uri,
			Password: password,
			DB:       dbInt,
		}),
	}
	if err := client.conn.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.conn.Get(ctx, key).Result()
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.conn.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Close() error {
	return c.conn.Close()
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	return c.conn.Del(ctx, key).Err()
}

func (c *RedisClient) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx).Err()
}
