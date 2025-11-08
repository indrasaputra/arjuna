package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/interceptor"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
)

const (
	connProtocol                = "tcp"
	prometheusReadHeaderTimeout = 5 * time.Second
)

// Server is responsible to act as gRPC server.
// It composes grpc.Server.
type Server struct {
	listener    net.Listener
	server      *grpc.Server
	httpServer  *http.Server
	name        string
	port        string
	serviceFunc []func(*grpc.Server)
}

// Config represents server's config.
type Config struct {
	IdempotencyStore          interceptor.IdempotencyStore
	Name                      string
	Port                      string
	Username                  string
	Password                  string
	AppliedBearerAuthMethods  []string
	AppliedBasicAuthMethods   []string
	AppliedIdempotencyMethods []string
	Secret                    []byte
}

// newGrpc creates an instance of Server.
func newServer(name, port string, options ...grpc.ServerOption) *Server {
	return &Server{
		server: grpc.NewServer(options...),
		name:   name,
		port:   port,
	}
}

// NewServer creates an instance of Server for used in development environment.
//
// These are list of interceptors that are attached (from innermost to outermost):
//   - Metrics, using Prometheus.
//   - Logging, using log/slog.
//   - Recoverer, using grpcrecovery.
func NewServer(cfg *Config) *Server {
	logger := sdklog.NewSlogLogger(cfg.Name)

	grpc_prometheus.EnableHandlingTimeHistogram()

	unary := defaultUnaryServerInterceptors(logger, cfg)
	stream := defaultStreamServerInterceptors(logger, cfg)
	unaryMdw := grpc.ChainUnaryInterceptor(unary...)
	streamMdw := grpc.ChainStreamInterceptor(stream...)
	trace := grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(otel.GetTracerProvider())))

	srv := newServer(cfg.Name, cfg.Port, trace, unaryMdw, streamMdw)
	grpc_prometheus.Register(srv.server)
	return srv
}

// Name returns server's name.
func (gs *Server) Name() string {
	return gs.name
}

// Port returns server's port.
func (gs *Server) Port() string {
	return gs.port
}

// AttachService attaches service to gRPC server.
// It will be called before serve.
func (gs *Server) AttachService(fn func(*grpc.Server)) {
	gs.serviceFunc = append(gs.serviceFunc, fn)
}

// EnablePrometheus registers prometheus metrics.
func (gs *Server) EnablePrometheus(port string) {
	grpc_prometheus.Register(gs.server)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: prometheusReadHeaderTimeout,
	}
	http.Handle("/metrics", promhttp.Handler())
	gs.httpServer = srv
}

// Serve runs the server.
// It basically runs grpc.Server.Serve and is a blocking.
func (gs *Server) Serve() error {
	for _, service := range gs.serviceFunc {
		service(gs.server)
	}

	var err error
	lc := net.ListenConfig{}
	gs.listener, err = lc.Listen(context.Background(), connProtocol, fmt.Sprintf(":%s", gs.port))
	if err != nil {
		return err
	}
	if gs.httpServer != nil {
		go func() {
			_ = gs.httpServer.ListenAndServe()
		}()
	}
	go func() {
		_ = gs.server.Serve(gs.listener)
	}()
	return nil
}

// GracefulStop blocks the server and wait for termination signal.
// The termination signal must be one of SIGINT or SIGTERM.
// Once it receives one of those signals, the gRPC server will perform graceful stop and close the listener.
//
// It receives Closer and will perform all closers before closing itself.
func (gs *Server) GracefulStop() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	gs.server.GracefulStop()
	if gs.listener != nil {
		_ = gs.listener.Close()
	}
}

// Stop immediately stops the gRPC server.
// Currently, this method exists just for the sake of testing.
// For production purpose, use GracefulStop().
func (gs *Server) Stop() {
	gs.server.Stop()
}

func defaultUnaryServerInterceptors(logger *slog.Logger, cfg *Config) []grpc.UnaryServerInterceptor {
	opts := []logging.Option{logging.WithLogOnEvents(logging.StartCall, logging.FinishCall)}

	interceptors := []grpc.UnaryServerInterceptor{
		grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(recoveryHandler)),
		logging.UnaryServerInterceptor(interceptor.SlogLogger(logger), opts...),
		grpc_prometheus.UnaryServerInterceptor,
		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(interceptor.AuthBasic(cfg.Username, cfg.Password)), selector.MatchFunc(interceptor.ApplyMethod(cfg.AppliedBasicAuthMethods...))),
		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(interceptor.AuthBearer(cfg.Secret)), selector.MatchFunc(interceptor.ApplyMethod(cfg.AppliedBearerAuthMethods...))),
	}

	if cfg.IdempotencyStore != nil && len(cfg.AppliedIdempotencyMethods) > 0 {
		idempotencyInterceptor := selector.UnaryServerInterceptor(
			interceptor.IdempotencyUnaryServerInterceptor(cfg.IdempotencyStore),
			selector.MatchFunc(interceptor.ApplyMethod(cfg.AppliedIdempotencyMethods...)),
		)
		interceptors = append(interceptors, idempotencyInterceptor)
	}

	return interceptors
}

func defaultStreamServerInterceptors(logger *slog.Logger, cfg *Config) []grpc.StreamServerInterceptor {
	opts := []logging.Option{logging.WithLogOnEvents(logging.StartCall, logging.FinishCall)}

	// Note: Idempotency interceptor is typically only used for unary requests,
	// not streaming requests, as streaming doesn't fit the idempotency pattern well
	return []grpc.StreamServerInterceptor{
		grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandler(recoveryHandler)),
		logging.StreamServerInterceptor(interceptor.SlogLogger(logger), opts...),
		grpc_prometheus.StreamServerInterceptor,
		selector.StreamServerInterceptor(auth.StreamServerInterceptor(interceptor.AuthBasic(cfg.Username, cfg.Password)), selector.MatchFunc(interceptor.ApplyMethod(cfg.AppliedBasicAuthMethods...))),
		selector.StreamServerInterceptor(auth.StreamServerInterceptor(interceptor.AuthBearer(cfg.Secret)), selector.MatchFunc(interceptor.ApplyMethod(cfg.AppliedBearerAuthMethods...))),
	}
}

func recoveryHandler(p any) error {
	return status.Errorf(codes.Unknown, "%v", p)
}
