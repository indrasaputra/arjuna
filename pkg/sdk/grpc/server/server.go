package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpclogsettable "github.com/grpc-ecosystem/go-grpc-middleware/logging/settable"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	connProtocol                = "tcp"
	prometheusReadHeaderTimeout = 5 * time.Second
)

// GrpcServer is responsible to act as gRPC server.
// It composes grpc.Server.
type GrpcServer struct {
	listener    net.Listener
	server      *grpc.Server
	httpServer  *http.Server
	name        string
	port        string
	serviceFunc []func(*grpc.Server)
}

// newGrpc creates an instance of GrpcServer.
func newGrpcServer(name, port string, options ...grpc.ServerOption) *GrpcServer {
	return &GrpcServer{
		server: grpc.NewServer(options...),
		name:   name,
		port:   port,
	}
}

// NewGrpcServer creates an instance of GrpcServer for used in development environment.
//
// These are list of interceptors that are attached (from innermost to outermost):
//   - Metrics, using Prometheus.
//   - Logging, using zap logger.
//   - Recoverer, using grpcrecovery.
func NewGrpcServer(name, port string) *GrpcServer {
	logger, _ := zap.NewProduction() // error is impossible, hence ignored.
	grpczap.SetGrpcLoggerV2(grpclogsettable.ReplaceGrpcLoggerV2(), logger)
	grpc_prometheus.EnableHandlingTimeHistogram()

	unary := defaultUnaryServerInterceptors(logger)
	stream := defaultStreamServerInterceptors(logger)
	unaryMdw := grpc.ChainUnaryInterceptor(unary...)
	streamMdw := grpc.ChainStreamInterceptor(stream...)
	trace := grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(otel.GetTracerProvider())))

	srv := newGrpcServer(name, port, trace, unaryMdw, streamMdw)
	grpc_prometheus.Register(srv.server)
	return srv
}

// Name returns server's name.
func (gs *GrpcServer) Name() string {
	return gs.name
}

// Port returns server's port.
func (gs *GrpcServer) Port() string {
	return gs.port
}

// AttachService attaches service to gRPC server.
// It will be called before serve.
func (gs *GrpcServer) AttachService(fn func(*grpc.Server)) {
	gs.serviceFunc = append(gs.serviceFunc, fn)
}

// EnablePrometheus registers prometheus metrics.
func (gs *GrpcServer) EnablePrometheus(port string) {
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
func (gs *GrpcServer) Serve() error {
	for _, service := range gs.serviceFunc {
		service(gs.server)
	}

	var err error
	gs.listener, err = net.Listen(connProtocol, fmt.Sprintf(":%s", gs.port))
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
func (gs *GrpcServer) GracefulStop() {
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
func (gs *GrpcServer) Stop() {
	gs.server.Stop()
}

func defaultUnaryServerInterceptors(logger *zap.Logger) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(recoveryHandler)),
		grpczap.UnaryServerInterceptor(logger),
		grpc_prometheus.UnaryServerInterceptor,
	}
}

func defaultStreamServerInterceptors(logger *zap.Logger) []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandler(recoveryHandler)),
		grpczap.StreamServerInterceptor(logger),
		grpc_prometheus.StreamServerInterceptor,
	}
}

func recoveryHandler(p interface{}) error {
	return status.Errorf(codes.Unknown, "%v", p)
}
