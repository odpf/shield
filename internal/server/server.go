package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgrpc"
	"github.com/odpf/salt/log"
	"github.com/odpf/salt/mux"
	"github.com/odpf/salt/server"
	"github.com/odpf/shield/internal/api"
	"github.com/odpf/shield/internal/api/v1beta1"
	"github.com/odpf/shield/internal/server/grpc_interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

/*
	func registerHandler(ctx context.Context, s *server.MuxServer, gw *server.GRPCGateway, deps api.Deps) {
		s.RegisterHandler("/admin*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "pong")
		}))

		s.RegisterHandler("/admin/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "pong")
		}))

		// grpc gateway api will have version endpoints
		s.SetGateway("/admin", gw)
		v1beta1.Register(ctx, s, gw, deps)
	}
*/
func Serve(
	ctx context.Context,
	logger log.Logger,
	cfg Config,
	nrApp newrelic.Application,
	deps api.Deps,
) error {
	httpMux := http.NewServeMux()
	httpMux.Handle("/admin*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}))

	httpMux.Handle("/admin/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}))

	grpcGateway := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(customHeaderMatcherFunc(map[string]bool{cfg.IdentityProxyHeader: true})))
	httpMux.Handle("/admin/", http.StripPrefix("/admin", grpcGateway))
	grpcServer := grpc.NewServer(getGRPCMiddleware(cfg, logger, nrApp))
	reflection.Register(grpcServer)

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPCPort)
	v1beta1.Register(ctx, address, grpcServer, grpcGateway, deps)

	logger.Info("[shield] api server starting", "http-port", cfg.HTTPPort, "grpc-port", cfg.GRPCPort)

	if err := mux.Serve(
		ctx,
		mux.WithHTTPTarget(fmt.Sprintf(":%d", cfg.HTTPPort), &http.Server{
			Handler:        httpMux,
			ReadTimeout:    120 * time.Second,
			WriteTimeout:   120 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}),
		mux.WithGRPCTarget(fmt.Sprintf(":%d", cfg.GRPCPort), grpcServer),
		mux.WithGracePeriod(5*time.Second),
	); !errors.Is(err, context.Canceled) {
		logger.Error("mux serve error", "err", err)
	}

	return nil
}

func Cleanup(ctx context.Context, s *server.MuxServer) {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	s.Shutdown(shutdownCtx)
}

// REVISIT: passing config.Shield as reference
func getGRPCMiddleware(cfg Config, logger log.Logger, nrApp newrelic.Application) grpc.ServerOption {
	recoveryFunc := func(p interface{}) (err error) {
		fmt.Println("-----------------------------")
		return status.Errorf(codes.Internal, "internal server error")
	}

	grpcRecoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	grpcZapLogger := zap.NewExample().Sugar()
	loggerZap, ok := logger.(*log.Zap)
	if ok {
		grpcZapLogger = loggerZap.GetInternalZapLogger()
	}
	return grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_interceptors.EnrichCtxWithIdentity(cfg.IdentityProxyHeader),
			grpc_zap.UnaryServerInterceptor(grpcZapLogger.Desugar()),
			grpc_recovery.UnaryServerInterceptor(grpcRecoveryOpts...),
			grpc_ctxtags.UnaryServerInterceptor(),
			nrgrpc.UnaryServerInterceptor(nrApp),
		))
}

func customHeaderMatcherFunc(headerKeys map[string]bool) func(key string) (string, bool) {
	return func(key string) (string, bool) {
		if _, ok := headerKeys[key]; ok {
			return key, true
		}
		return runtime.DefaultHeaderMatcher(key)
	}
}
