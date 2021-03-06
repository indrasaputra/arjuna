package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/indrasaputra/arjuna/gateway/config"
	"github.com/indrasaputra/arjuna/gateway/server"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

func main() {
	cfg, err := config.NewConfig(".env")
	checkError(err)

	gatewayServer := server.NewGrpcGateway(cfg.Port)
	registerGrpcGatewayService(context.Background(), gatewayServer, cfg, grpc.WithTransportCredentials(insecure.NewCredentials()))

	_ = gatewayServer.Serve()
}

func registerGrpcGatewayService(ctx context.Context, gatewayServer *server.GrpcGateway, cfg *config.Config, options ...grpc.DialOption) {
	gatewayServer.AttachService(func(server *runtime.ServeMux) error {
		if err := apiv1.RegisterUserCommandServiceHandlerFromEndpoint(ctx, server, cfg.UserServiceAddress, options); err != nil {
			return err
		}
		return nil
	})
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
