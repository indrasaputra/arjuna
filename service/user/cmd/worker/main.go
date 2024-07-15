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
	connauth "github.com/indrasaputra/arjuna/service/user/internal/connection/auth"
	orcact "github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/activity"
	orcwork "github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/workflow"
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
	bunDB, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)
	authClient, err := builder.BuildAuthClient(cfg.AuthServiceHost)
	checkError(err)

	co := connauth.NewAuth(authClient)
	db := postgres.NewUser(bunDB)

	act := orcact.NewRegisterUserActivity(co, db)

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
