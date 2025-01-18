// Server main program.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/server"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/builder"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
)

func main() {
	command := &cobra.Command{Use: "auth", Short: "Start the service."}

	command.AddCommand(&cobra.Command{
		Use:   "api",
		Short: "Run the API server.",
		Run:   API,
	})
	command.AddCommand(&cobra.Command{
		Use:   "seed",
		Short: "Run the seeder.",
		Run:   Seed,
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

	_, err = trace.NewProvider(ctx, cfg.Tracer)
	checkError(err)

	app.Logger = sdklog.NewLogger(cfg.AppEnv)

	pool, err := postgres.NewPgxPool(cfg.Postgres)
	checkError(err)
	defer pool.Close()
	queries := builder.BuildQueries(pool, postgres.NewTxGetter())

	dep := &builder.Dependency{
		Config:             cfg,
		SigningKey:         cfg.Token.SecretKey,
		ExpiryTimeInMinute: cfg.Token.ExpiryTimeInMinutes,
		Queries:            queries,
	}

	c := &server.Config{
		Name:                     cfg.ServiceName,
		Port:                     cfg.Port,
		Secret:                   []byte(cfg.Token.SecretKey),
		Username:                 cfg.Username,
		Password:                 cfg.Password,
		AppliedBearerAuthMethods: strings.Split(cfg.AppliedAuthBearer, ","),
		AppliedBasicAuthMethods:  strings.Split(cfg.AppliedAuthBasic, ","),
	}
	srv := server.NewServer(c)
	registerGrpcService(srv, dep)
	srv.EnablePrometheus(cfg.PrometheusPort)

	_ = srv.Serve()
	fmt.Println("server start.. waiting signal")
	srv.GracefulStop()
}

// Seed is the entry point for running the seeder.
func Seed(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)
	pool, err := postgres.NewPgxPool(cfg.Postgres)
	checkError(err)
	defer pool.Close()

	val := openJSON("test/fixture/accounts.json")

	insertAccounts(ctx, pool, val)
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

func insertAccounts(ctx context.Context, db uow.Tr, val []byte) {
	var accounts []*entity.Account
	_ = json.Unmarshal(val, &accounts)

	query := "INSERT INTO accounts (id, user_id, email, password, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, NOW(), NOW(), $5, $6)"
	for _, account := range accounts {
		password, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.MinCost)
		_, err := db.Exec(ctx, query, account.ID, account.UserID, account.Email, string(password), account.ID, account.ID)
		checkInsertError(err)
	}
	log.Printf("Successfully insert %d accounts\n", len(accounts))
}

func checkInsertError(err error) {
	if uow.IsUniqueViolationError(err) {
		return
	}
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
