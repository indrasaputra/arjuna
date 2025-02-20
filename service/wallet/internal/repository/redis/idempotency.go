package redis

import (
	"context"
	"log/slog"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/indrasaputra/arjuna/service/wallet/entity"
)

const (
	redisNotFound = "redis: nil"
	ttl           = 24 * time.Hour
)

// IdempotencyKey is responsible to connect idempotency flow with redis.
type IdempotencyKey struct {
	client goredis.Cmdable
}

// NewIdempotencyKey creates an instance of IdempotencyKey.
func NewIdempotencyKey(client goredis.Cmdable) *IdempotencyKey {
	return &IdempotencyKey{client: client}
}

// Exists check if the key exists in redis.
func (i *IdempotencyKey) Exists(ctx context.Context, key string) (bool, error) {
	args := goredis.SetArgs{Get: true, TTL: ttl}
	_, err := i.client.SetArgs(ctx, key, 1, args).Result()
	if err != nil && err.Error() == redisNotFound {
		return false, nil
	}
	if err != nil {
		slog.ErrorContext(ctx, "[IdempotencyKeyRedis-Exists] internal error", "error", err)
		return false, entity.ErrInternal("fail check in redis")
	}
	return true, nil
}
