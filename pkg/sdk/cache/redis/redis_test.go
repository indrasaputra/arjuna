package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/pkg/sdk/cache/redis"
)

var (
	testCtx = context.Background()
)

type IdempotencySuite struct {
	idp  *redis.Idempotency
	mock redismock.ClientMock
}

func TestNewIdempotency(t *testing.T) {
	t.Run("successfully create an instance of Idempotency", func(t *testing.T) {
		st := createIdempotencySuite()
		assert.NotNil(t, st.idp)
	})
}

func TestIdempotency_Set(t *testing.T) {
	key := "set-key"
	expectedKey := "idempotency:" + key

	t.Run("set returns error", func(t *testing.T) {
		st := createIdempotencySuite()
		st.mock.ExpectSet(expectedKey, []byte("1"), time.Hour).SetErr(assert.AnError)

		err := st.idp.Set(testCtx, key, []byte("1"), 0)

		assert.Error(t, err)
	})

	t.Run("set returns success", func(t *testing.T) {
		st := createIdempotencySuite()
		st.mock.ExpectSet(expectedKey, []byte("1"), time.Hour).SetVal("1")

		err := st.idp.Set(testCtx, key, []byte("1"), 0)

		assert.NoError(t, err)
	})
}

func TestIdempotency_Get(t *testing.T) {
	key := "set-key"
	expectedKey := "idempotency:" + key

	t.Run("get returns error", func(t *testing.T) {
		st := createIdempotencySuite()
		st.mock.ExpectGet(expectedKey).SetErr(assert.AnError)

		_, err := st.idp.Get(testCtx, key)

		assert.Error(t, err)
	})

	t.Run("get returns goredis nil", func(t *testing.T) {
		st := createIdempotencySuite()
		st.mock.ExpectGet(expectedKey).RedisNil()

		_, err := st.idp.Get(testCtx, key)

		assert.NoError(t, err)
	})

	t.Run("get returns success", func(t *testing.T) {
		st := createIdempotencySuite()
		st.mock.ExpectGet(expectedKey).SetVal("1")

		val, err := st.idp.Get(testCtx, key)

		assert.NoError(t, err)
		assert.Equal(t, []byte("1"), val)
	})
}

func createIdempotencySuite() *IdempotencySuite {
	c, m := redismock.NewClientMock()
	i := redis.NewIdempotency(c, time.Hour)
	return &IdempotencySuite{
		idp:  i,
		mock: m,
	}
}
