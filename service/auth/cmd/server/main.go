// Server main program.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/server"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/builder"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)

	_, err = trace.NewProvider(ctx, cfg.Tracer)
	checkError(err)

	app.Logger = sdklog.NewLogger(cfg.AppEnv)

	bunDB, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)

	dep := &builder.Dependency{
		Config:             cfg,
		DB:                 bunDB,
		SigningKey:         cfg.Token.SecretKey,
		ExpiryTimeInMinute: cfg.Token.ExpiryTimeInMinutes,
	}

	c := &server.Config{
		Name:           cfg.ServiceName,
		Port:           cfg.Port,
		Secret:         []byte(cfg.Token.SecretKey),
		SkippedMethods: strings.Split(cfg.SkippedAuth, ","),
	}
	srv := server.NewServer(c)
	registerGrpcService(srv, dep)
	srv.EnablePrometheus(cfg.PrometheusPort)

	_ = srv.Serve()
	fmt.Println("server start.. waiting signal")
	srv.GracefulStop()
}

func registerGrpcService(srv *server.Server, dep *builder.Dependency) {
	// start register all module's gRPC handlers
	command, err := builder.BuildAuthHandler(dep)
	if err != nil {
		log.Fatalf("fail build auth command handler: %v", err)
	}
	health := handler.NewHealth()

	srv.AttachService(func(server *grpc.Server) {
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
