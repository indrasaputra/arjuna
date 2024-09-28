package uow

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trm "github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DB defines an interface for database functionality.
// This interface exists for the sake of Unit of Work (UoW) pattern.
// The moment this project is being made, Go doesn't have proper UoW library
// that can abstract the process.
//
// I will agree to the opinion that says this is not proper/correct way to achieve
// UoW or transaction as a business.
// I read many articles and this is the conclusion I made using all knowledges I have.
type DB interface {
	// Begin begins the transaction.
	Begin(ctx context.Context) (Tx, error)
	// Exec executes the query without getting the data from the row.
	Exec(ctx context.Context, query string, args ...interface{}) (int64, error)
	// Query runs the query given.
	Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Tx defines an interface for business transaction functionality.
// This interface exists for the sake of Unit of Work (UoW) pattern.
// See more in the comment on DB interface.
type Tx interface {
	DB
	// Commit commits the transaction.
	Commit(ctx context.Context) error
	// Rollback rolls back the transaction.
	Rollback(ctx context.Context) error
}

// UnitWorker represents a unit of work.
// It actually just to begin and finish the transaction.
type UnitWorker struct {
	db DB
}

// UnitOfWork defines the high-level use case of Unit of Work pattern.
type UnitOfWork interface {
	// Begin begins the transaction.
	Begin(ctx context.Context) (Tx, error)
	// Finish finishes the transaction.
	Finish(ctx context.Context, tx Tx, err error) error
}

// NewUnitWorker creates an instance of UnitWorker.
func NewUnitWorker(db DB) *UnitWorker {
	return &UnitWorker{db: db}
}

// Begin begins the transaction.
func (u *UnitWorker) Begin(ctx context.Context) (Tx, error) {
	return u.db.Begin(ctx)
}

// Finish finishes the transaction.
func (u *UnitWorker) Finish(ctx context.Context, tx Tx, err error) error {
	if err != nil {
		return tx.Rollback(ctx)
	}
	return tx.Commit(ctx)
}

// TxManager defines the interface to manage transaction.
// Copied from https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/trm/v2/manager#Manager.
// It is used mostly for testing.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
	DoWithSettings(ctx context.Context, s trm.Settings, fn func(ctx context.Context) error) (err error)
	Init(ctx context.Context, s trm.Settings) (context.Context, manager.Closer, error)
}

// Tr defines an interface for transaction.
// Copied from https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2#Tr
// It is used mostly for testing.
type Tr interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// TxGetter defines an interface to get transaction from context or use db from param.
// Copied from https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2#CtxGetter.
// It is used mostly for testing.
type TxGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmpgx.Tr) trmpgx.Tr
}
