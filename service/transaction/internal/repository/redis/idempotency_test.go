package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/redis"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type IdempotencySuite struct {
	idp  *redis.IdempotencyKey
	mock redismock.ClientMock
}

func TestNewIdempotency(t *testing.T) {
	t.Run("successfully create an instance of Idempotency", func(t *testing.T) {
		st := createIdempotencySuite()
		assert.NotNil(t, st.idp)
	})
}

func TestIdempotency_Exists(t *testing.T) {
	app.Logger = sdklog.NewLogger(testEnv)
	key := "idempotency"

	t.Run("set returns nil", func(t *testing.T) {
		args := goredis.SetArgs{Get: true, TTL: 24 * time.Hour}
		st := createIdempotencySuite()
		st.mock.ExpectSetArgs(key, 1, args).RedisNil()

		res, err := st.idp.Exists(testCtx, key)

		assert.NoError(t, err)
		assert.False(t, res)
	})

	t.Run("set returns error", func(t *testing.T) {
		args := goredis.SetArgs{Get: true, TTL: 24 * time.Hour}
		st := createIdempotencySuite()
		st.mock.ExpectSetArgs(key, 1, args).SetErr(assert.AnError)

		res, err := st.idp.Exists(testCtx, key)

		assert.Error(t, err)
		assert.False(t, res)
	})

	t.Run("set returns no error", func(t *testing.T) {
		args := goredis.SetArgs{Get: true, TTL: 24 * time.Hour}
		st := createIdempotencySuite()
		st.mock.ExpectSetArgs(key, 1, args).SetVal("1")

		res, err := st.idp.Exists(testCtx, key)

		assert.NoError(t, err)
		assert.True(t, res)
	})
}

func createIdempotencySuite() *IdempotencySuite {
	c, m := redismock.NewClientMock()
	i := redis.NewIdempotencyKey(c)
	return &IdempotencySuite{
		idp:  i,
		mock: m,
	}
}
