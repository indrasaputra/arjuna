package builder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/repository"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

var (
	postgresConnFormat = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s sslmode=%s"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config         *config.Config
	PgxPool        *pgxpool.Pool
	KeycloakClient kcsdk.Keycloak
}

// BuildUserCommandHandler builds toggle command handler including all of its dependencies.
func BuildUserCommandHandler(dep *Dependency) (*handler.UserCommand, error) {
	kcConfig := &keycloak.Config{
		Client:        dep.KeycloakClient,
		Realm:         dep.Config.Keycloak.Realm,
		AdminUsername: dep.Config.Keycloak.AdminUser,
		AdminPassword: dep.Config.Keycloak.AdminPassword,
	}
	kc, err := keycloak.NewUser(kcConfig)
	if err != nil {
		return nil, err
	}
	pg := postgres.NewUser(dep.PgxPool)
	regRepo := repository.NewUserRegistrator(kc, pg)
	registrator := service.NewUserRegistrator(regRepo)
	return handler.NewUserCommand(registrator), nil
}

// BuildPostgrePgxPool builds a pool of pgx client.
func BuildPostgrePgxPool(cfg config.Postgres) (*pgxpool.Pool, error) {
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

// BuildKeycloakClient builds a keycloak client.
func BuildKeycloakClient(cfg config.Keycloak) kcsdk.Keycloak {
	hc := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	client := kcsdk.NewClient(hc, cfg.Address)
	return client
}
