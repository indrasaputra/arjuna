package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/auth/entity"
)

func TestErrInternal(t *testing.T) {
	t.Run("success get internal error", func(t *testing.T) {
		err := entity.ErrInternal("")

		assert.Contains(t, err.Error(), "rpc error: code = Internal")
	})
}

func TestErrEmptyField(t *testing.T) {
	t.Run("success get empty field error", func(t *testing.T) {
		fields := []string{"clientId", "email", "password"}
		for _, field := range fields {
			err := entity.ErrEmptyField(field)

			assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
		}
	})
}

func TestErrUnauthorized(t *testing.T) {
	t.Run("success get unauthorized error", func(t *testing.T) {
		err := entity.ErrUnauthorized()

		assert.Contains(t, err.Error(), "rpc error: code = Unauthenticated")
	})
}

func TestErrInvalidArgument(t *testing.T) {
	t.Run("success get invalid argument error", func(t *testing.T) {
		err := entity.ErrInvalidArgument("arg")

		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	})
}
