package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
)

const (
	// errCodeUniqueViolation is derived from https://www.postgresql.org/docs/11/errcodes-appendix.html
	errCodeUniqueViolation = "23505"
)

var (
	postgresConnFormat = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"
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
	Port     string `env:"POSTGRES_PORT,default=5432"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Name     string `env:"POSTGRES_NAME,required"`
	SSLMode  string `env:"POSTGRES_SSL_MODE,default=disable"`
}

// NewDBWithPgx creates a bun.DB using pgx as driver.
func NewDBWithPgx(cfg Config) (*bun.DB, error) {
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
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), connCfg)
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
