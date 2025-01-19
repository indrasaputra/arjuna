package uow

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trm "github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

// TxManager defines the interface to manage transaction.
// Copied from https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/trm/v2/manager#Manager.
// It is used mostly for testing.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
	DoWithSettings(ctx context.Context, s trm.Settings, fn func(ctx context.Context) error) (err error)
	Init(ctx context.Context, s trm.Settings) (context.Context, manager.Closer, error)
}

// Transactional defines an interface for transactional.
type Transactional trmpgx.Transactional

// NewTxManager creates a new transaction manager using pgx as driver.
func NewTxManager(tx Transactional) (*manager.Manager, error) {
	// see https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2#NewDefaultFactory
	// to understand the usage of trmpgx.NewDefaultFactory.
	return manager.New(trmpgx.NewDefaultFactory(tx))
}

// Tr defines an interface for transaction.
type Tr trmpgx.Tr

// TxGetter defines an interface to get transaction from context or use db from param.
// Copied from https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2#CtxGetter.
// It is used mostly for testing.
type TxGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmpgx.Tr) trmpgx.Tr
}

// NewTxGetter creates a new transaction getter.
func NewTxGetter() *trmpgx.CtxGetter {
	return trmpgx.DefaultCtxGetter
}
