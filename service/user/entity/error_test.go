package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

func TestErrInternal(t *testing.T) {
	t.Run("success get internal error", func(t *testing.T) {
		err := entity.ErrInternal("")

		assert.Contains(t, err.Error(), "rpc error: code = Internal")
	})
}

func TestErrEmptyUser(t *testing.T) {
	t.Run("success get empty user error", func(t *testing.T) {
		err := entity.ErrEmptyUser()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrAlreadyExists(t *testing.T) {
	t.Run("success get already exists error", func(t *testing.T) {
		err := entity.ErrAlreadyExists()

		assert.Contains(t, err.Error(), "rpc error: code = AlreadyExists")
	})
}

func TestErrInvalidName(t *testing.T) {
	t.Run("success get invalid name error", func(t *testing.T) {
		err := entity.ErrInvalidName()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrInvalidEmail(t *testing.T) {
	t.Run("success get invalid email error", func(t *testing.T) {
		err := entity.ErrInvalidEmail()

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}

func TestErrNotFound(t *testing.T) {
	t.Run("success get user not found error", func(t *testing.T) {
		err := entity.ErrNotFound()

		assert.Contains(t, err.Error(), "rpc error: code = NotFound")
	})
}
