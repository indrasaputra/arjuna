// Server main program.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	skdpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/server"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	connauth "github.com/indrasaputra/arjuna/service/user/internal/connection/auth"
	connwallet "github.com/indrasaputra/arjuna/service/user/internal/connection/wallet"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	orcact "github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/activity"
	orcwork "github.com/indrasaputra/arjuna/service/user/internal/orchestration/temporal/workflow"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

func main() {
	command := &cobra.Command{Use: "user", Short: "Start the service."}

	command.AddCommand(&cobra.Command{
		Use:   "api",
		Short: "Run the API server.",
		Run:   API,
	})
	command.AddCommand(&cobra.Command{
		Use:   "worker",
		Short: "Run the worker.",
		Run:   Worker,
	})
	command.AddCommand(&cobra.Command{
		Use:   "relayer",
		Short: "Run the relayer.",
		Run:   Relayer,
	})
	command.AddCommand(&cobra.Command{
		Use:   "seed",
		Short: "Run the seeder.",
		Run:   Seed,
	})

	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}

// API is the entry point for running the API server.
func API(_ *cobra.Command, _ []string) {
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
	redisClient, err := builder.BuildRedisClient(&cfg.Redis)
	checkError(err)

	dep := &builder.Dependency{
		TemporalClient: temporalClient,
		Config:         cfg,
		DB:             bunDB,
		RedisClient:    redisClient,
	}

	c := &server.Config{
		Name:                     cfg.ServiceName,
		Port:                     cfg.Port,
		Secret:                   []byte(cfg.SecretKey),
		Username:                 cfg.Username,
		Password:                 cfg.Password,
		AppliedBearerAuthMethods: strings.Split(cfg.AppliedAuthBearer, ","),
		AppliedBasicAuthMethods:  strings.Split(cfg.AppliedAuthBasic, ","),
	}
	srv := server.NewServer(c)
	registerGrpcService(srv, dep)
	srv.EnablePrometheus(cfg.PrometheusPort)

	_ = srv.Serve()
	srv.GracefulStop()
}

// Worker is the entry point for running the worker server.
func Worker(_ *cobra.Command, _ []string) {
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
	authClient, err := builder.BuildAuthClient(cfg.AuthServiceHost, cfg.AuthServiceUsername, cfg.AuthServicePassword)
	checkError(err)
	walletClient, err := builder.BuildWalletClient(cfg.WalletServiceHost, cfg.WalletServiceUsername, cfg.WalletServicePassword)
	checkError(err)

	ac := connauth.NewAuth(authClient)
	wc := connwallet.NewWallet(walletClient)
	db := postgres.NewUser(bunDB)

	act := orcact.NewRegisterUserActivity(ac, wc, db)

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

// Relayer is the entry point for running the Relayer server.
func Relayer(_ *cobra.Command, _ []string) {
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

// Seed is the entry point for running the seeder.
func Seed(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)
	db, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)

	val := openJSON("test/fixture/users.json")

	insertUsers(ctx, db, val)
}

func registerGrpcService(srv *server.Server, dep *builder.Dependency) {
	// start register all module's gRPC handlers
	command := builder.BuildUserCommandHandler(dep)
	commandInternal := builder.BuildUserCommandInternalHandler(dep)
	query := builder.BuildUserQueryHandler(dep)
	health := handler.NewHealth()

	srv.AttachService(func(server *grpc.Server) {
		apiv1.RegisterUserCommandServiceServer(server, command)
		apiv1.RegisterUserCommandInternalServiceServer(server, commandInternal)
		apiv1.RegisterUserQueryServiceServer(server, query)
		grpc_health_v1.RegisterHealthServer(server, health)
	})
	// end of register all module's gRPC handlers
}

func openJSON(file string) []byte {
	jsonFile, err := os.Open(filepath.Clean(file))
	checkError(err)
	defer func() {
		_ = jsonFile.Close()
	}()

	val, err := io.ReadAll(jsonFile)
	checkError(err)

	return val
}

func insertUsers(ctx context.Context, db *skdpostgres.BunDB, val []byte) {
	var users []*entity.User
	_ = json.Unmarshal(val, &users)

	query := "INSERT INTO users (id, name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
	for _, user := range users {
		_, err := db.Exec(ctx, query, user.ID, user.Name)
		checkError(err)
	}
	log.Printf("Successfully insert %d users\n", len(users))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
