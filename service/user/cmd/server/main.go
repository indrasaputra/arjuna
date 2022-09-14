package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/server"
)

func main() {
	cfg, err := config.NewConfig(".env")
	checkError(err)

	keycloakClient := builder.BuildKeycloakClient(cfg.Keycloak)
	postgresPool, err := builder.BuildPostgrePgxPool(cfg.Postgres)
	checkError(err)

	dep := &builder.Dependency{
		PgxPool:        postgresPool,
		KeycloakClient: keycloakClient,
		Config:         cfg,
	}

	grpcServer := server.NewGrpcServer(cfg.Port)
	registerGrpcService(grpcServer, dep)

	_ = grpcServer.Serve()
	fmt.Println("server start.. waiting signal")
	grpcServer.GracefulStop()
}

func registerGrpcService(grpcServer *server.GrpcServer, dep *builder.Dependency) {
	// start register all module's gRPC handlers
	command, err := builder.BuildUserCommandHandler(dep)
	if err != nil {
		log.Fatalf("fail build user command handler: %v", err)
	}
	query := builder.BuildUserQueryHandler(dep)
	health := handler.NewHealth()

	grpcServer.AttachService(func(server *grpc.Server) {
		apiv1.RegisterUserCommandServiceServer(server, command)
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
