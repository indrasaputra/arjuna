// Worker main program
package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	orcact "github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/activity"
	orcwork "github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/workflow"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)

	app.Logger = sdklog.NewLogger(cfg.AppEnv)

	_, err = trace.NewProvider(ctx, cfg.Tracer)
	checkError(err)

	temporalClient, err := builder.BuildTemporalClient(cfg.Temporal.Address)
	checkError(err)
	defer temporalClient.Close()
	keycloakClient := builder.BuildKeycloakClient(cfg.Keycloak)
	bunDB, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)

	dep := &builder.Dependency{
		KeycloakClient: keycloakClient,
		Config:         cfg,
	}

	kcConfig := &keycloak.Config{
		Client:        dep.KeycloakClient,
		Realm:         dep.Config.Keycloak.Realm,
		AdminUsername: dep.Config.Keycloak.AdminUser,
		AdminPassword: dep.Config.Keycloak.AdminPassword,
	}
	kc, _ := keycloak.NewUser(kcConfig)
	db := postgres.NewUser(bunDB)

	act := orcact.NewRegisterUserActivity(kc, db)

	w := worker.New(temporalClient, orcwork.TaskQueueRegisterUser, worker.Options{
		DisableRegistrationAliasing: true,
	})
	w.RegisterWorkflow(orcwork.RegisterUser)
	w.RegisterActivityWithOptions(act, activity.RegisterOptions{Name: "RegisterUserActivity", SkipInvalidStructFunctions: true})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Panic("Unable to start worker", err)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
