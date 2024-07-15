// Server main program.
package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/server"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
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
	redisClient, err := builder.BuildRedisClient(&cfg.Redis)
	checkError(err)

	dep := &builder.Dependency{
		TemporalClient: temporalClient,
		Config:         cfg,
		DB:             bunDB,
		RedisClient:    redisClient,
	}

	grpcServer := server.NewGrpcServer(cfg.ServiceName, cfg.Port)
	registerGrpcService(grpcServer, dep)
	grpcServer.EnablePrometheus(cfg.PrometheusPort)

	_ = grpcServer.Serve()
	grpcServer.GracefulStop()
}

func registerGrpcService(grpcServer *server.GrpcServer, dep *builder.Dependency) {
	// start register all module's gRPC handlers
	command := builder.BuildUserCommandHandler(dep)
	commandInternal := builder.BuildUserCommandInternalHandler(dep)
	query := builder.BuildUserQueryHandler(dep)
	health := handler.NewHealth()

	grpcServer.AttachService(func(server *grpc.Server) {
		apiv1.RegisterUserCommandServiceServer(server, command)
		apiv1.RegisterUserCommandInternalServiceServer(server, commandInternal)
		apiv1.RegisterUserQueryServiceServer(server, query)
		grpc_health_v1.RegisterHealthServer(server, health)
	})
	// end of register all module's gRPC handlers
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
