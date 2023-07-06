// Worker main program
package main

import (
	"context"
	"time"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	orcwork "github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/workflow"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
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

	p := postgres.NewUserOutbox(bunDB)
	w := orcwork.NewRegisterUserWorkflow(temporalClient)

	svc := service.NewUserRelayRegistrar(p, w)

	for {
		app.Logger.Infof(ctx, "running user registration relayer at %v", time.Now())
		if err := svc.Register(ctx); err != nil {
			app.Logger.Errorf(ctx, "error during relay register: %v", err)
		}
		time.Sleep(time.Duration(cfg.RelayerSleepTimeMillisecond) * time.Millisecond)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
