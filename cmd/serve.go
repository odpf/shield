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

	"github.com/odpf/shield/pkg/server"

	"github.com/odpf/shield/pkg/server/consts"

	"github.com/odpf/shield/core/invitation"

	"github.com/odpf/shield/pkg/mailer"

	"github.com/odpf/shield/core/permission"
	"github.com/odpf/shield/internal/bootstrap"

	"github.com/odpf/shield/core/deleter"

	_ "github.com/authzed/authzed-go/proto/authzed/api/v0"
	_ "github.com/jackc/pgx/v4/stdlib"
	newrelic "github.com/newrelic/go-agent"
	"github.com/odpf/shield/core/authenticate"
	"github.com/odpf/shield/core/authenticate/session"
	"github.com/odpf/shield/core/metaschema"

	"github.com/odpf/shield/config"
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
	"github.com/odpf/shield/internal/store/blob"
	"github.com/odpf/shield/internal/store/postgres"
	"github.com/odpf/shield/internal/store/spicedb"
	"github.com/odpf/shield/pkg/db"

	"github.com/odpf/salt/log"
	"github.com/pkg/profile"
	"google.golang.org/grpc/codes"
)

var (
	ruleCacheRefreshDelay = time.Minute * 2
)

func StartServer(logger *log.Zap, cfg *config.Shield) error {
	logger.Info("shield starting", "version", Version)
	if profiling := os.Getenv("SHIELD_PROFILE"); profiling == "true" || profiling == "1" {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	}

	// @TODO: need to inject custom logger wrapper over zap into ctx to use it internally
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
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

	nrApp, err := setupNewRelic(cfg.NewRelic, logger)
	if err != nil {
		return err
	}

	spiceDBClient, err := spicedb.New(cfg.SpiceDB, logger)
	if err != nil {
		return err
	}

	deps, err := buildAPIDependencies(logger, cfg, resourceBlobRepository, dbClient, spiceDBClient)
	if err != nil {
		return err
	}
	// load metadata schema in memory from db
	if schemas, err := deps.MetaSchemaService.List(context.Background()); err != nil {
		logger.Warn("metaschemas initialization failed", "err", err)
	} else {
		logger.Info("metaschemas loaded", "count", len(schemas))
	}

	// apply schema
	if err = deps.BootstrapService.MigrateSchema(ctx); err != nil {
		return err
	}

	// apply roles over nil org id
	// nil org is the default org of platform
	if err = deps.BootstrapService.MigrateRoles(ctx); err != nil {
		return err
	}
	// promote normal users to superusers
	if err = deps.BootstrapService.MakeSuperUsers(ctx); err != nil {
		return err
	}

	// session service initialization and cleanup
	if err := deps.SessionService.InitSessions(context.Background()); err != nil {
		logger.Warn("sessions database cleanup failed", "err", err)
	}
	defer func() {
		logger.Debug("cleaning up cron jobs")
		deps.SessionService.Close()
	}()

	if err := deps.RegistrationService.InitFlows(context.Background()); err != nil {
		logger.Warn("flows database cleanup failed", "err", err)
	}
	defer func() {
		deps.RegistrationService.Close()
	}()

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

	// serving server
	return server.Serve(ctx, logger, cfg.App, nrApp, deps)
}

func buildAPIDependencies(
	logger log.Logger,
	cfg *config.Shield,
	resourceBlobRepository *blob.ResourcesRepository,
	dbc *db.Client,
	sdb *spicedb.SpiceDB,
) (api.Deps, error) {
	sessionService := session.NewService(logger, postgres.NewSessionRepository(logger, dbc), consts.SessionValidity)

	namespaceRepository := postgres.NewNamespaceRepository(dbc)
	namespaceService := namespace.NewService(namespaceRepository)

	authzSchemaRepository := spicedb.NewSchemaRepository(logger, sdb)
	authzRelationRepository := spicedb.NewRelationRepository(sdb, cfg.SpiceDB.FullyConsistent)

	relationPGRepository := postgres.NewRelationRepository(dbc)
	relationService := relation.NewService(relationPGRepository, authzRelationRepository)

	permissionRepository := postgres.NewPermissionRepository(dbc)
	permissionService := permission.NewService(permissionRepository)

	policyPGRepository := postgres.NewPolicyRepository(dbc)
	policyService := policy.NewService(policyPGRepository, relationService)

	roleRepository := postgres.NewRoleRepository(dbc)
	roleService := role.NewService(roleRepository, relationService, permissionService)

	userRepository := postgres.NewUserRepository(dbc)
	userService := user.NewService(userRepository, sessionService, relationService)

	groupRepository := postgres.NewGroupRepository(dbc)
	groupService := group.NewService(groupRepository, relationService, userService)

	resourceSchemaRepository := blob.NewSchemaConfigRepository(resourceBlobRepository.Bucket)
	bootstrapService := bootstrap.NewBootstrapService(
		cfg.App.Admin,
		resourceSchemaRepository,
		namespaceService,
		roleService,
		permissionService,
		userService,
		authzSchemaRepository,
	)

	organizationRepository := postgres.NewOrganizationRepository(dbc)
	organizationService := organization.NewService(organizationRepository, relationService, userService)

	projectRepository := postgres.NewProjectRepository(dbc)
	projectService := project.NewService(projectRepository, relationService, userService)

	metaschemaRepository := postgres.NewMetaSchemaRepository(logger, dbc)
	metaschemaService := metaschema.NewService(metaschemaRepository)

	resourcePGRepository := postgres.NewResourceRepository(dbc)
	resourceService := resource.NewService(
		resourcePGRepository,
		resourceBlobRepository,
		relationService,
		userService,
	)

	var mailDialer mailer.Dialer = mailer.NewMockDialer()
	if cfg.App.Mailer.SMTPHost != "" && cfg.App.Mailer.SMTPHost != "smtp.example.com" {
		mailDialer = mailer.NewDialerImpl(cfg.App.Mailer.SMTPHost,
			cfg.App.Mailer.SMTPPort,
			cfg.App.Mailer.SMTPUsername,
			cfg.App.Mailer.SMTPPassword,
			cfg.App.Mailer.SMTPInsecure,
			cfg.App.Mailer.Headers,
		)
	}
	registrationService := authenticate.NewRegistrationService(logger, cfg.App.Authentication,
		postgres.NewFlowRepository(logger, dbc), userService, mailDialer)

	invitationService := invitation.NewService(mailDialer, postgres.NewInvitationRepository(logger, dbc),
		organizationService, groupService, userService, relationService)
	cascadeDeleter := deleter.NewCascadeDeleter(organizationService, projectService, resourceService,
		groupService, policyService, roleService, invitationService)

	dependencies := api.Deps{
		DisableOrgsListing:  cfg.App.DisableOrgsListing,
		DisableUsersListing: cfg.App.DisableUsersListing,
		OrgService:          organizationService,
		ProjectService:      projectService,
		GroupService:        groupService,
		RoleService:         roleService,
		PolicyService:       policyService,
		UserService:         userService,
		NamespaceService:    namespaceService,
		PermissionService:   permissionService,
		RelationService:     relationService,
		ResourceService:     resourceService,
		SessionService:      sessionService,
		RegistrationService: registrationService,
		DeleterService:      cascadeDeleter,
		MetaSchemaService:   metaschemaService,
		BootstrapService:    bootstrapService,
		InvitationService:   invitationService,
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
