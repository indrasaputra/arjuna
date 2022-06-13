package builder

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

var (
	postgresConnFormat = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s sslmode=%s"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config  *config.Config
	PgxPool *pgxpool.Pool
}

// BuildUserCommandHandler builds toggle command handler including all of its dependencies.
func BuildUserCommandHandler(dep *Dependency) *handler.UserCommand {
	repo := postgres.NewUser(dep.PgxPool)
	registrator := service.NewUserRegistrator(repo)
	return handler.NewUserCommand(registrator)
}

// BuildPostgrePgxPool builds a pool of pgx client.
func BuildPostgrePgxPool(cfg *config.Postgres) (*pgxpool.Pool, error) {
	connCfg := fmt.Sprintf(postgresConnFormat,
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.MaxOpenConns,
		cfg.MaxConnLifetime,
		cfg.MaxIdleLifetime,
		cfg.SSLMode,
	)
	return pgxpool.Connect(context.Background(), connCfg)
}
