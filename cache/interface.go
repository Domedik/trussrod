package cache

import (
	"context"
	"time"
)

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Ping(ctx context.Context) error
	Del(ctx context.Context, key string) error
	Close() error
}
