// Server main program.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/arjuna/pkg/sdk/cache/redis"
	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/server"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
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

	logger := sdklog.NewSlogLogger(cfg.ServiceName)
	slog.SetDefault(logger)

	_, err = trace.NewProvider(ctx, cfg.Tracer)
	checkError(err)

	temporalClient, err := builder.BuildTemporalClient(cfg.Temporal.Address)
	checkError(err)
	defer temporalClient.Close()
	pool, err := sdkpostgres.NewPgxPool(cfg.Postgres)
	checkError(err)
	defer pool.Close()
	txm, err := uow.NewTxManager(pool)
	checkError(err)
	redisClient, err := redis.NewRedisClient(cfg.Redis)
	checkError(err)
	idempotencyStore := redis.NewIdempotency(redisClient, cfg.Redis.TTL)

	queries := builder.BuildQueries(pool, uow.NewTxGetter())

	dep := &builder.Dependency{
		TemporalClient: temporalClient,
		Config:         cfg,
		TxManager:      txm,
		Queries:        queries,
	}

	c := &server.Config{
		Name:                      cfg.ServiceName,
		Port:                      cfg.Port,
		Secret:                    []byte(cfg.SecretKey),
		Username:                  cfg.Username,
		Password:                  cfg.Password,
		AppliedBearerAuthMethods:  strings.Split(cfg.AppliedAuthBearer, ","),
		AppliedBasicAuthMethods:   strings.Split(cfg.AppliedAuthBasic, ","),
		AppliedIdempotencyMethods: strings.Split(cfg.AppliedIdempotency, ","),
		IdempotencyStore:          idempotencyStore,
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

	logger := sdklog.NewSlogLogger(cfg.ServiceName)
	slog.SetDefault(logger)

	_, err = trace.NewProvider(ctx, cfg.Tracer)
	checkError(err)

	temporalClient, err := builder.BuildTemporalClient(cfg.Temporal.Address)
	checkError(err)
	defer temporalClient.Close()
	authClient, err := builder.BuildAuthClient(cfg.AuthServiceHost, cfg.AuthServiceUsername, cfg.AuthServicePassword)
	checkError(err)
	walletClient, err := builder.BuildWalletClient(cfg.WalletServiceHost, cfg.WalletServiceUsername, cfg.WalletServicePassword)
	checkError(err)
	pool, err := sdkpostgres.NewPgxPool(cfg.Postgres)
	checkError(err)
	defer pool.Close()
	queries := builder.BuildQueries(pool, uow.NewTxGetter())

	ac := connauth.NewAuth(authClient)
	wc := connwallet.NewWallet(walletClient)
	db := postgres.NewUser(queries)

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

	logger := sdklog.NewSlogLogger(cfg.ServiceName)
	slog.SetDefault(logger)

	_, err = trace.NewProvider(ctx, cfg.Tracer)
	checkError(err)

	temporalClient, err := builder.BuildTemporalClient(cfg.Temporal.Address)
	checkError(err)
	defer temporalClient.Close()

	pool, err := sdkpostgres.NewPgxPool(cfg.Postgres)
	checkError(err)
	defer pool.Close()
	txm, err := uow.NewTxManager(pool)
	checkError(err)
	queries := builder.BuildQueries(pool, uow.NewTxGetter())

	p := postgres.NewUserOutbox(queries)
	w := orcwork.NewRegisterUserWorkflow(temporalClient)

	svc := service.NewUserRelayRegistrar(p, w, txm)

	for {
		slog.InfoContext(ctx, "running user registration relayer", "time", time.Now())
		if err := svc.Register(ctx); err != nil {
			slog.ErrorContext(ctx, "error during relay register", "error", err)
		}
		time.Sleep(time.Duration(cfg.RelayerSleepTimeMillisecond) * time.Millisecond)
	}
}

// Seed is the entry point for running the seeder.
func Seed(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)
	pool, err := sdkpostgres.NewPgxPool(cfg.Postgres)
	checkError(err)
	defer pool.Close()

	val := openJSON("test/fixture/users.json")

	insertUsers(ctx, pool, val)
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

func insertUsers(ctx context.Context, db uow.Tr, val []byte) {
	var users []*entity.User
	_ = json.Unmarshal(val, &users)

	query := `INSERT INTO users (id, name, created_at, updated_at, created_by, updated_by)
				VALUES ($1, $2, NOW(), NOW(), $3, $4)
				ON CONFLICT (id) DO NOTHING;`
	for _, user := range users {
		_, err := db.Exec(ctx, query, user.ID, user.Name, user.ID, user.ID)
		checkError(err)
	}
	log.Printf("Successfully insert %d users\n", len(users))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
