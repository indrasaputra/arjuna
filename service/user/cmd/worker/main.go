// Worker main program
package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
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

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	cfg, err := config.NewConfig(".env")
	checkError(err)

	exp, err := trace.NewJaegerExporter(cfg.Tracer)
	checkError(err)
	_ = trace.NewProvider(cfg.Tracer, exp)

	app.Logger = sdklog.NewLogger(cfg.AppEnv)

	keycloakClient := builder.BuildKeycloakClient(cfg.Keycloak)
	temporalClient, err := builder.BuildTemporalClient()
	checkError(err)
	bunDB, err := builder.BuildBunDB(ctx, cfg.Postgres)
	checkError(err)

	dep := &builder.Dependency{
		KeycloakClient: keycloakClient,
		TemporalClient: temporalClient,
		Config:         cfg,
		DB:             bunDB,
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

	w := worker.New(c, orcwork.TaskQueueRegisterUser, worker.Options{
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
