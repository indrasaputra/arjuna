// Server main program.
package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/server"
)

func main() {
	ctx := context.Background()

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

	grpcServer := server.NewGrpcServer(cfg.Port)
	registerGrpcService(ctx, grpcServer, dep)

	_ = grpcServer.Serve()
	fmt.Println("server start.. waiting signal")
	grpcServer.GracefulStop()
}

func registerGrpcService(ctx context.Context, grpcServer *server.GrpcServer, dep *builder.Dependency) {
	// start register all module's gRPC handlers
	command := builder.BuildUserCommandHandler(dep)
	commandInternal, err := builder.BuildUserCommandInternalHandler(ctx, dep)
	if err != nil {
		log.Fatalf("fail build user command internal handler: %v", err)
	}
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
