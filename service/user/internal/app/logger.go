package app

import (
	"go.uber.org/zap"
)

// Logger provides logging functionality.
type Logger struct {
	*zap.SugaredLogger
}
