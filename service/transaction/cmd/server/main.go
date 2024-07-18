// Server main program.
package main

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/server"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/builder"
	"github.com/indrasaputra/arjuna/service/transaction/internal/config"
	"github.com/indrasaputra/arjuna/service/transaction/internal/grpc/handler"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)

	app.Logger = sdklog.NewLogger(cfg.AppEnv)

	_, err = trace.NewProvider(ctx, cfg.Tracer)
	checkError(err)

	bunDB, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)
	redisClient, err := builder.BuildRedisClient(&cfg.Redis)
	checkError(err)

	dep := &builder.Dependency{
		Config:      cfg,
		DB:          bunDB,
		RedisClient: redisClient,
	}

	c := &server.Config{
		Name:           cfg.ServiceName,
		Port:           cfg.Port,
		Secret:         []byte(cfg.SecretKey),
		SkippedMethods: strings.Split(cfg.SkippedAuth, ","),
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