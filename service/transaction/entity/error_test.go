package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/transaction/entity"
)

func TestErrInternal(t *testing.T) {
	t.Run("success get internal error", func(t *testing.T) {
		err := entity.ErrInternal("")

		assert.Contains(t, err.Error(), "rpc error: code = Internal")
	})
}

func TestErrAlreadyExists(t *testing.T) {
	t.Run("success get already exists error", func(t *testing.T) {
		err := entity.ErrAlreadyExists()

		assert.Contains(t, err.Error(), "rpc error: code = AlreadyExists")
	})
}

func TestErrEmptyTransaction(t *testing.T) {
	t.Run("success get empty transaction error", func(t *testing.T) {
		err := entity.ErrEmptyTransaction()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidSender(t *testing.T) {
	t.Run("success get invalid sender error", func(t *testing.T) {
		err := entity.ErrInvalidSender()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidReceiver(t *testing.T) {
	t.Run("success get invalid receiver error", func(t *testing.T) {
		err := entity.ErrInvalidReceiver()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidAmount(t *testing.T) {
	t.Run("success get invalid amount error", func(t *testing.T) {
		err := entity.ErrInvalidAmount()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrMissingIdempotencyKey(t *testing.T) {
	t.Run("success get missing idempotency key error", func(t *testing.T) {
		err := entity.ErrMissingIdempotencyKey()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidWallet(t *testing.T) {
	t.Run("success get invalid wallet error", func(t *testing.T) {
		err := entity.ErrInvalidWallet()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}
