// Server main program.
package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/builder"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/server"
)

func main() {
	cfg, err := config.NewConfig(".env")
	checkError(err)

	exp, err := trace.NewJaegerExporter(cfg.Tracer)
	checkError(err)
	_ = trace.NewProvider(cfg.Tracer, exp)

	app.Logger = sdklog.NewLogger(cfg.AppEnv)

	keycloakClient := builder.BuildKeycloakClient(cfg.Keycloak)

	dep := &builder.Dependency{
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
	command, err := builder.BuildAuthHandler(dep)
	if err != nil {
		log.Fatalf("fail build auth command handler: %v", err)
	}
	health := handler.NewHealth()

	grpcServer.AttachService(func(server *grpc.Server) {
		apiv1.RegisterAuthServiceServer(server, command)
		grpc_health_v1.RegisterHealthServer(server, health)
	})
	// end of register all module's gRPC handlers
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
