// Server main program.
package main

import (
	"context"
	"log"

	grpclogsettable "github.com/grpc-ecosystem/go-grpc-middleware/logging/settable"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/indrasaputra/arjuna/gateway/config"
	"github.com/indrasaputra/arjuna/gateway/server"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)

	_, err = trace.NewProvider(ctx, cfg.Tracer)
	checkError(err)

	gatewayServer := server.NewGrpcGateway(cfg.Port)
	options := defaultGrpcServerOptions()
	registerGrpcGatewayService(context.Background(), gatewayServer, cfg, options...)

	log.Println("running grpc gateway server...")
	_ = gatewayServer.Serve()
}

func registerGrpcGatewayService(ctx context.Context, gatewayServer *server.GrpcGateway, cfg *config.Config, options ...grpc.DialOption) {
	gatewayServer.AttachService(func(server *runtime.ServeMux) error {
		if err := apiv1.RegisterUserCommandServiceHandlerFromEndpoint(ctx, server, cfg.UserServiceAddress, options); err != nil {
			return err
		}
		if err := apiv1.RegisterUserQueryServiceHandlerFromEndpoint(ctx, server, cfg.UserServiceAddress, options); err != nil {
			return err
		}
		return apiv1.RegisterAuthServiceHandlerFromEndpoint(ctx, server, cfg.AuthServiceAddress, options)
	})
}

func defaultGrpcServerOptions() []grpc.DialOption {
	logger, _ := zap.NewProduction() // error is impossible, hence ignored.
	grpczap.SetGrpcLoggerV2(grpclogsettable.ReplaceGrpcLoggerV2(), logger)

	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(otel.GetTracerProvider())),
			grpczap.UnaryClientInterceptor(logger),
			grpc_prometheus.UnaryClientInterceptor,
		),
		grpc.WithChainStreamInterceptor(
			otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(otel.GetTracerProvider())),
			grpczap.StreamClientInterceptor(logger),
			grpc_prometheus.StreamClientInterceptor,
		),
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
