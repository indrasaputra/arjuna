package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	defaultTTL           = 1 * time.Hour
	idempotencyKeyPrefix = "idempotency:"
)

// Idempotency is responsible to connect idempotency flow with redis.
type Idempotency struct {
	client goredis.Cmdable
	ttl    time.Duration
}

// Config holds configuration for Redis.
type Config struct {
	Address string        `env:"REDIS_ADDRESS,default=localhost:6379"`
	TTL     time.Duration `env:"REDIS_TTL,default=1h"`
}

// NewRedisClient creates an instance of Redis client.
func NewRedisClient(cfg Config) (*goredis.Client, error) {
	opt := &goredis.Options{
		Addr: cfg.Address,
	}
	client := goredis.NewClient(opt)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewIdempotency creates an instance of Idempotency.
func NewIdempotency(client goredis.Cmdable, ttl time.Duration) *Idempotency {
	if ttl <= 0 {
		ttl = defaultTTL
	}
	return &Idempotency{client: client, ttl: ttl}
}

// Get retrieves a response from Redis by key.
func (i *Idempotency) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := i.client.Get(ctx, idempotencyKeyPrefix+key).Bytes()
	if err == goredis.Nil {
		return nil, nil
	}
	return val, err
}

// Set stores a response in Redis with the given key and TTL.
func (i *Idempotency) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = i.ttl
	}
	return i.client.Set(ctx, idempotencyKeyPrefix+key, value, ttl).Err()
}
