package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

// PgxPool defines a little interface for pgxpool functionality.
// Since in the real implementation we can use pgxpool.Pool,
// this interface exists mostly for testing purpose.
type PgxPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// PgxTx defines a little interface for transaction functionality.
// Since in the real implementation we can use pgx.Tx,
// this interface exists mostly for testing purpose.
type PgxTx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	PgxPool
}

// TxContext represents transaction context.
type TxContext string

var (
	// TxContextKey sets as tx-key.
	TxContextKey = TxContext("tx-key")
)

// DatabaseTransaction is responsible to provide atomic transaction in PostgreSQL.
type DatabaseTransaction struct {
	pool PgxPool
}

// NewDatabaseTransaction creates an instance of DatabaseTransaction.
func NewDatabaseTransaction(pool PgxPool) *DatabaseTransaction {
	return &DatabaseTransaction{pool: pool}
}

// WithinTransaction executes given function inside a transaction.
func (dt *DatabaseTransaction) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := dt.pool.Begin(ctx)
	if err != nil {
		return entity.ErrInternal(err.Error())
	}

	err = fn(injectContextWithTx(ctx, tx))

	if err != nil {
		if errRbck := tx.Rollback(ctx); errRbck != nil {
			log.Println("[DatabaseTransaction-WithinTransaction] error rollback: ", errRbck)
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		log.Println("[DatabaseTransaction-WithinTransaction] error commit: ", err)
	}
	return nil
}

func injectContextWithTx(ctx context.Context, tx PgxTx) context.Context {
	return context.WithValue(ctx, TxContextKey, tx)
}

func extractTxFromContext(ctx context.Context) PgxTx {
	if tx, ok := ctx.Value(TxContextKey).(PgxTx); ok {
		return tx
	}
	return nil
}
