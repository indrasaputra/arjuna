// Server main program.
package main

import (
	"context"
	"log"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/server"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/transaction/internal/builder"
	"github.com/indrasaputra/arjuna/service/transaction/internal/config"
	"github.com/indrasaputra/arjuna/service/transaction/internal/grpc/handler"
)

func main() {
	command := &cobra.Command{Use: "transaction", Short: "Start the service."}

	command.AddCommand(&cobra.Command{
		Use:   "api",
		Short: "Run the API server.",
		Run:   API,
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

	pool, err := postgres.NewPgxPool(cfg.Postgres)
	checkError(err)
	defer pool.Close()
	queries := builder.BuildQueries(pool, uow.NewTxGetter())
	redisClient, err := builder.BuildRedisClient(&cfg.Redis)
	checkError(err)

	dep := &builder.Dependency{
		Config:      cfg,
		RedisClient: redisClient,
		Queries:     queries,
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

func registerGrpcService(srv *server.Server, dep *builder.Dependency) {
	// start register all module's gRPC handlers
	command := builder.BuildTransactionCommandHandler(dep)
	health := handler.NewHealth()

	srv.AttachService(func(server *grpc.Server) {
		apiv1.RegisterTransactionCommandServiceServer(server, command)
		grpc_health_v1.RegisterHealthServer(server, health)
	})
	// end of register all module's gRPC handlers
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
