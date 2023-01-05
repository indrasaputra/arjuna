// Worker main program
package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/workflow/temporal"
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

	w := worker.New(c, temporal.TaskQueueRegisterUser, worker.Options{
		DisableRegistrationAliasing: true,
	})
	w.RegisterWorkflow(temporal.RegisterUser)
	w.RegisterActivityWithOptions(kc, activity.RegisterOptions{Name: "Keycloak", SkipInvalidStructFunctions: true})
	w.RegisterActivityWithOptions(db, activity.RegisterOptions{Name: "Postgres", SkipInvalidStructFunctions: true})

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
