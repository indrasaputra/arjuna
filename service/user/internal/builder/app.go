package builder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.temporal.io/sdk/client"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	"github.com/indrasaputra/arjuna/service/user/internal/workflow/temporal"
)

var (
	postgresConnFormat = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s sslmode=%s"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config         *config.Config
	PgxPool        *pgxpool.Pool
	KeycloakClient kcsdk.Keycloak
	TemporalClient client.Client
}

// BuildUserCommandHandler builds user command handler including all of its dependencies.
func BuildUserCommandHandler(dep *Dependency) (*handler.UserCommand, error) {
	// kcConfig := &keycloak.Config{
	// 	Client:        dep.KeycloakClient,
	// 	Realm:         dep.Config.Keycloak.Realm,
	// 	AdminUsername: dep.Config.Keycloak.AdminUser,
	// 	AdminPassword: dep.Config.Keycloak.AdminPassword,
	// }
	// kc, err := keycloak.NewUser(kcConfig)
	// if err != nil {
	// 	return nil, err
	// }
	// pg := postgres.NewUser(dep.PgxPool)
	tp := temporal.NewRegisterUserWorkflow(dep.TemporalClient)
	// regRepo := repository.NewUserRegistrar(kc, pg)
	rg := service.NewUserRegistrar(tp)
	return handler.NewUserCommand(rg), nil
}

// BuildUserCommandInternalHandler builds user command handler including all of its dependencies.
func BuildUserCommandInternalHandler(dep *Dependency) (*handler.UserCommandInternal, error) {
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
	tx := postgres.NewDatabaseTransaction(dep.PgxPool)
	deleter := service.NewUserDeleter(pg, kc, tx)
	return handler.NewUserCommandInternal(deleter), nil
}

// BuildUserQueryHandler builds user query handler including all of its dependencies.
func BuildUserQueryHandler(dep *Dependency) *handler.UserQuery {
	pg := postgres.NewUser(dep.PgxPool)
	getter := service.NewUserGetter(pg)
	return handler.NewUserQuery(getter)
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

// BuildTemporalClient builds temporal client.
func BuildTemporalClient() client.Client {
	c, _ := client.Dial(client.Options{})
	return c
}
