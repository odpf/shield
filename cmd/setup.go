package cmd

import (
	"fmt"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgrpc"
	"github.com/odpf/salt/log"
	"github.com/odpf/shield/config"
	"github.com/odpf/shield/grpc_interceptors"
	"github.com/odpf/shield/pkg/sql"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setupNewRelic(cfg config.NewRelic, logger log.Logger) newrelic.Application {
	nrCfg := newrelic.NewConfig(cfg.AppName, cfg.License)
	nrCfg.Enabled = cfg.Enabled

	nrApp, err := newrelic.NewApplication(nrCfg)
	if err != nil {
		logger.Fatal("failed to load Newrelic Application")
	}
	return nrApp
}

// REVISIT: passing config.Shield as reference
func getGRPCMiddleware(cfg *config.Shield, logger log.Logger) grpc.ServerOption {
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "internal server error")
	}
	opts := []grpcRecovery.Option{
		grpcRecovery.WithRecoveryHandler(customFunc),
	}

	return grpc.UnaryInterceptor(
		grpcMiddleware.ChainUnaryServer(
			grpc_interceptors.EnrichCtxWithIdentity(cfg.App.IdentityProxyHeader),
			grpczap.UnaryServerInterceptor(zap.NewExample()),
			grpcRecovery.UnaryServerInterceptor(opts...),
			grpcctxtags.UnaryServerInterceptor(),
			nrgrpc.UnaryServerInterceptor(setupNewRelic(cfg.NewRelic, logger)),
		))
}

func setupDB(cfg config.DBConfig, logger log.Logger) (*sql.SQL, func()) {
	db, err := sql.New(sql.Config{
		Driver:          cfg.Driver,
		URL:             cfg.URL,
		MaxIdleConns:    cfg.MaxIdleConns,
		MaxOpenConns:    cfg.MaxOpenConns,
		ConnMaxLifeTime: cfg.ConnMaxLifeTime,
	})

	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to setup db: %s", err.Error()))
	}

	return db, func() { db.Close() }
}
