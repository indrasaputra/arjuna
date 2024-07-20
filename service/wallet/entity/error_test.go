package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/wallet/entity"
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

func TestErrEmptyWallet(t *testing.T) {
	t.Run("success get empty wallet error", func(t *testing.T) {
		err := entity.ErrEmptyWallet()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidBalance(t *testing.T) {
	t.Run("success get invalid balance error", func(t *testing.T) {
		err := entity.ErrInvalidBalance()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidUser(t *testing.T) {
	t.Run("success get invalid user error", func(t *testing.T) {
		err := entity.ErrInvalidUser()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrMissingIdempotencyKey(t *testing.T) {
	t.Run("success get missing idempotency key error", func(t *testing.T) {
		err := entity.ErrMissingIdempotencyKey()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidAmount(t *testing.T) {
	t.Run("success get invalid amount error", func(t *testing.T) {
		err := entity.ErrInvalidAmount()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrSameAccount(t *testing.T) {
	t.Run("success get same account error", func(t *testing.T) {
		err := entity.ErrSameAccount()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInsufficientBalance(t *testing.T) {
	t.Run("success get insufficient balance error", func(t *testing.T) {
		err := entity.ErrInsufficientBalance()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidTransfer(t *testing.T) {
	t.Run("success get invalid transfer error", func(t *testing.T) {
		err := entity.ErrInvalidTransfer()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}
