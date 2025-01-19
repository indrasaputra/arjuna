package postgres

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
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
