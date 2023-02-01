package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/authzed/authzed-go/proto/authzed/api/v0"
	_ "github.com/jackc/pgx/v4/stdlib"
	newrelic "github.com/newrelic/go-agent"

	"github.com/odpf/shield/config"
	"github.com/odpf/shield/core/action"
	"github.com/odpf/shield/core/group"
	"github.com/odpf/shield/core/namespace"
	"github.com/odpf/shield/core/organization"
	"github.com/odpf/shield/core/policy"
	"github.com/odpf/shield/core/project"
	"github.com/odpf/shield/core/relation"
	"github.com/odpf/shield/core/resource"
	"github.com/odpf/shield/core/role"
	"github.com/odpf/shield/core/user"
	"github.com/odpf/shield/internal/api"
	"github.com/odpf/shield/internal/schema"
	"github.com/odpf/shield/internal/server"
	"github.com/odpf/shield/internal/store/blob"
	"github.com/odpf/shield/internal/store/postgres"
	"github.com/odpf/shield/internal/store/spicedb"
	"github.com/odpf/shield/pkg/db"

	"github.com/odpf/salt/log"
	salt_server "github.com/odpf/salt/server"
	"github.com/pkg/profile"
	"google.golang.org/grpc/codes"
)

var (
	ruleCacheRefreshDelay = time.Minute * 2
)

func StartServer(logger *log.Zap, cfg *config.Shield) error {
	if profiling := os.Getenv("SHIELD_PROFILE"); profiling == "true" || profiling == "1" {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	}

	// @TODO: need to inject custom logger wrapper over zap into ctx to use it internally
	ctx, cancelFunc := context.WithCancel(salt_server.HandleSignals(context.Background()))
	defer cancelFunc()

	dbClient, err := setupDB(cfg.DB, logger)
	if err != nil {
		return err
	}
	defer func() {
		logger.Info("cleaning up db")
		dbClient.Close()
	}()

	// load resource config
	if cfg.App.ResourcesConfigPath == "" {
		return errors.New("resource config path cannot be left empty")
	}

	resourceBlobFS, err := blob.NewStore(ctx, cfg.App.ResourcesConfigPath, cfg.App.ResourcesConfigPathSecret)
	if err != nil {
		return err
	}
	resourceBlobRepository := blob.NewResourcesRepository(logger, resourceBlobFS)
	if err := resourceBlobRepository.InitCache(ctx, ruleCacheRefreshDelay); err != nil {
		return err
	}
	defer func() {
		logger.Info("cleaning up resource blob")
		defer resourceBlobRepository.Close()
	}()

	spiceDBClient, err := spicedb.New(cfg.SpiceDB, logger)
	if err != nil {
		return err
	}

	nrApp, err := setupNewRelic(cfg.NewRelic, logger)
	if err != nil {
		return err
	}

	//
	actionRepository := postgres.NewActionRepository(dbClient)
	actionService := action.NewService(actionRepository)

	roleRepository := postgres.NewRoleRepository(dbClient)
	roleService := role.NewService(roleRepository)

	policyPGRepository := postgres.NewPolicyRepository(dbClient)
	policySpiceRepository := spicedb.NewPolicyRepository(spiceDBClient)
	policyService := policy.NewService(policyPGRepository)

	namespaceRepository := postgres.NewNamespaceRepository(dbClient)
	namespaceService := namespace.NewService(namespaceRepository)

	s := schema.NewSchemaMigrationService(
		blob.NewSchemaConfigRepository(resourceBlobFS),
		namespaceService,
		roleService,
		actionService,
		policyService,
		policySpiceRepository,
	)

	err = s.RunMigrations(ctx)
	if err != nil {
		return err
	}

	deps, err := buildAPIDependencies(ctx, logger, resourceBlobRepository, dbClient, spiceDBClient)
	if err != nil {
		return err
	}

	// serving proxies
	cbs, cps, err := serveProxies(ctx, logger, cfg.App.IdentityProxyHeader, cfg.App.UserIDHeader, cfg.Proxy, deps.ResourceService, deps.RelationService, deps.UserService, deps.ProjectService)
	if err != nil {
		return err
	}
	defer func() {
		// clean up stage
		logger.Info("cleaning up rules proxy blob")
		for _, f := range cbs {
			if err := f(); err != nil {
				logger.Warn("error occurred during shutdown rules proxy blob storages", "err", err)
			}
		}

		logger.Info("cleaning up proxies")
		for _, f := range cps {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*20)
			if err := f(shutdownCtx); err != nil {
				shutdownCancel()
				logger.Warn("error occurred during shutdown proxies", "err", err)
				continue
			}
			shutdownCancel()
		}
	}()

	keystrokeTermChan := make(chan os.Signal, 1)
	// we'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(keystrokeTermChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	// serving server
	muxServer, err := server.Serve(ctx, logger, cfg.App, nrApp, deps)
	if err != nil {
		return err
	}
	defer func() {
		logger.Info("cleaning up server")
		server.Cleanup(ctx, muxServer)
	}()

	// wait for termination
	select {
	case <-ctx.Done():
		fmt.Printf("process: ctx done bye\n")
		break
	case <-keystrokeTermChan:
		fmt.Printf("process: kill signal received. bye \n")
		break
	}

	return nil
}

func buildAPIDependencies(
	ctx context.Context,
	logger log.Logger,
	resourceBlobRepository *blob.ResourcesRepository,
	dbc *db.Client,
	sdb *spicedb.SpiceDB,
) (api.Deps, error) {
	actionRepository := postgres.NewActionRepository(dbc)
	actionService := action.NewService(actionRepository)

	namespaceRepository := postgres.NewNamespaceRepository(dbc)
	namespaceService := namespace.NewService(namespaceRepository)

	userRepository := postgres.NewUserRepository(dbc)
	userService := user.NewService(userRepository)

	roleRepository := postgres.NewRoleRepository(dbc)
	roleService := role.NewService(roleRepository)

	relationPGRepository := postgres.NewRelationRepository(dbc)
	relationSpiceRepository := spicedb.NewRelationRepository(sdb)
	relationService := relation.NewService(relationPGRepository, relationSpiceRepository, roleService, userService)

	groupRepository := postgres.NewGroupRepository(dbc)
	groupService := group.NewService(groupRepository, relationService, userService)

	organizationRepository := postgres.NewOrganizationRepository(dbc)
	organizationService := organization.NewService(organizationRepository, relationService, userService)

	projectRepository := postgres.NewProjectRepository(dbc)
	projectService := project.NewService(projectRepository, relationService, userService)

	policyPGRepository := postgres.NewPolicyRepository(dbc)
	policyService := policy.NewService(policyPGRepository)

	resourcePGRepository := postgres.NewResourceRepository(dbc)
	resourceService := resource.NewService(
		resourcePGRepository,
		resourceBlobRepository,
		relationService,
		userService)

	dependencies := api.Deps{
		OrgService:       organizationService,
		UserService:      userService,
		ProjectService:   projectService,
		GroupService:     groupService,
		RelationService:  relationService,
		ResourceService:  resourceService,
		RoleService:      roleService,
		PolicyService:    policyService,
		ActionService:    actionService,
		NamespaceService: namespaceService,
	}
	return dependencies, nil
}

func setupNewRelic(cfg config.NewRelic, logger log.Logger) (newrelic.Application, error) {
	nrCfg := newrelic.NewConfig(cfg.AppName, cfg.License)
	nrCfg.Enabled = cfg.Enabled
	nrCfg.ErrorCollector.IgnoreStatusCodes = []int{
		http.StatusNotFound,
		http.StatusUnauthorized,
		int(codes.Unauthenticated),
		int(codes.PermissionDenied),
		int(codes.InvalidArgument),
		int(codes.AlreadyExists),
	}

	if nrCfg.Enabled {
		nrApp, err := newrelic.NewApplication(nrCfg)
		if err != nil {
			return nil, errors.New("failed to load Newrelic Application")
		}
		return nrApp, nil
	}
	return nil, nil
}

func setupDB(cfg db.Config, logger log.Logger) (dbc *db.Client, err error) {
	// prefer use pgx instead of lib/pq for postgres to catch pg error
	if cfg.Driver == "postgres" {
		cfg.Driver = "pgx"
	}
	dbc, err = db.New(cfg)
	if err != nil {
		err = fmt.Errorf("failed to setup db: %w", err)
		return
	}

	return
}
