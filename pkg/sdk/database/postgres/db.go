package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
)

const (
	// errCodeUniqueViolation is derived from https://www.postgresql.org/docs/11/errcodes-appendix.html
	errCodeUniqueViolation = "23505"
)

var (
	postgresConnFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
)

var (
	// ErrAlreadyExist is returned when the new record violates unique key restriction.
	ErrAlreadyExist = errors.New("record already exist")
	// ErrNullDB is returned when nil is sent to mandatory db parameter.
	ErrNullDB = errors.New("db instance is null")
)

// Config holds configuration for PostgreSQL.
type Config struct {
	Host     string `env:"POSTGRES_HOST,default=localhost"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Name     string `env:"POSTGRES_NAME,required"`
	SSLMode  string `env:"POSTGRES_SSL_MODE,default=disable"`
	Port     int    `env:"POSTGRES_PORT,default=5432"`
}

// NewTxGetter creates a new transaction getter.
func NewTxGetter() *trmpgx.CtxGetter {
	return trmpgx.DefaultCtxGetter
}

// NewTxManager creates a new transaction manager using pgx as driver.
func NewTxManager(pool *pgxpool.Pool) (*manager.Manager, error) {
	// see https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2#NewDefaultFactory
	// to understand the usage of trmpgx.NewDefaultFactory.
	return manager.New(trmpgx.NewDefaultFactory(pool))
}

// NewPgxPool creates a new pgx pool.
func NewPgxPool(cfg Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(postgresConnFormat,
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)
	connCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	connCfg.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	return pgxpool.NewWithConfig(context.Background(), connCfg)
}

// NewDBWithPgx creates a bun.DB using pgx as driver.
func NewDBWithPgx(cfg Config) (*bun.DB, error) {
	pool, err := NewPgxPool(cfg)
	if err != nil {
		return nil, err
	}

	s := stdlib.OpenDBFromPool(pool)
	return bun.NewDB(s, pgdialect.New()), nil
}

// BunDB wraps uptrace/bun to comply with internal use.
type BunDB struct {
	db bun.IDB
}

// NewBunDB creates an instance of Bun.
func NewBunDB(db bun.IDB) (*BunDB, error) {
	if db == nil {
		return nil, ErrNullDB
	}
	return &BunDB{db: db}, nil
}

// Begin begins the transaction.
func (b *BunDB) Begin(ctx context.Context) (uow.Tx, error) {
	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return nil, fmt.Errorf("[BunDB] begin error: %v", err)
	}
	return &BunTx{tx: tx}, nil
}

// Exec executes the given query.
func (b *BunDB) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := b.db.ExecContext(ctx, query, args...)
	if isUniqueViolationErr(err) {
		return 0, ErrAlreadyExist
	}
	if err != nil {
		return 0, fmt.Errorf("[BunDB] exec error: %v", err)
	}
	return res.RowsAffected()
}

// Query queries the given query.
func (b *BunDB) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return b.db.NewRaw(query, args...).Scan(ctx, dest)
}

// BunTx wraps uptrace/bun tx to comply with internal use.
type BunTx struct {
	tx bun.Tx
}

// Begin begins the transaction
func (t *BunTx) Begin(_ context.Context) (uow.Tx, error) {
	return t, nil
}

// Commit commits the transaction.
func (t *BunTx) Commit(_ context.Context) error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction.
func (t *BunTx) Rollback(_ context.Context) error {
	return t.tx.Rollback()
}

// Exec executes the given query.
func (t *BunTx) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := t.tx.ExecContext(ctx, query, args...)
	if isUniqueViolationErr(err) {
		return 0, ErrAlreadyExist
	}
	if err != nil {
		return 0, fmt.Errorf("[BunTx] exec error: %v", err)
	}
	return res.RowsAffected()
}

// Query queries the given query.
func (t *BunTx) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.NewRaw(query, args...).Scan(ctx, dest)
}

func isUniqueViolationErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), errCodeUniqueViolation)
}

// TxDB implements DB with transaction.
type TxDB struct {
	db       *pgxpool.Pool
	txGetter *trmpgx.CtxGetter
}

// NewTxDB creates an instance of TxDB.
func NewTxDB(pool *pgxpool.Pool, txGetter *trmpgx.CtxGetter) *TxDB {
	return &TxDB{
		db:       pool,
		txGetter: txGetter,
	}
}

// Exec executes the given query using transaction.
func (d *TxDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	tx := d.txGetter.DefaultTrOrDB(ctx, d.db)
	return tx.Exec(ctx, sql, args...)
}

// Query executes the given query using transaction.
func (d *TxDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	tx := d.txGetter.DefaultTrOrDB(ctx, d.db)
	return tx.Query(ctx, sql, args...)
}

// QueryRow executes the given query using transaction.
func (d *TxDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	tx := d.txGetter.DefaultTrOrDB(ctx, d.db)
	return tx.QueryRow(ctx, sql, args...)
}
