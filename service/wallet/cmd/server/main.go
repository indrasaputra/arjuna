// Server main program.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/server"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/builder"
	"github.com/indrasaputra/arjuna/service/wallet/internal/config"
	"github.com/indrasaputra/arjuna/service/wallet/internal/grpc/handler"
)

func main() {
	command := &cobra.Command{Use: "wallet", Short: "Start the service."}

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

// Seed is the entry point for running the seeder.
func Seed(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)
	db, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)

	val := openJSON("test/fixture/wallets.json")

	insertWallets(ctx, db, val)
}

func registerGrpcService(srv *server.Server, dep *builder.Dependency) {
	// start register all module's gRPC handlers
	command := builder.BuildWalletCommandHandler(dep)
	health := handler.NewHealth()

	srv.AttachService(func(server *grpc.Server) {
		apiv1.RegisterWalletCommandServiceServer(server, command)
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

func insertWallets(ctx context.Context, db *postgres.BunDB, val []byte) {
	var wallets []*entity.Wallet
	_ = json.Unmarshal(val, &wallets)

	query := "INSERT INTO wallets (id, user_id, balance, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
	for _, wallet := range wallets {
		_, err := db.Exec(ctx, query, wallet.ID, wallet.UserID, wallet.Balance)
		checkError(err)
	}
	log.Printf("Successfully insert %d wallets\n", len(wallets))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
