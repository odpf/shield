package resource

import (
	"context"
	"errors"
	"strings"

	"github.com/odpf/shield/core/organization"
	"github.com/odpf/shield/core/project"
	"github.com/odpf/shield/pkg/utils"

	"github.com/odpf/shield/core/relation"
	"github.com/odpf/shield/core/user"
	"github.com/odpf/shield/internal/bootstrap/schema"
)

type RelationService interface {
	Create(ctx context.Context, rel relation.Relation) (relation.Relation, error)
	CheckPermission(ctx context.Context, subject relation.Subject, object relation.Object, permName string) (bool, error)
	Delete(ctx context.Context, rel relation.Relation) error
}

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type ProjectService interface {
	Get(ctx context.Context, idOrName string) (project.Project, error)
}

type OrgService interface {
	Get(ctx context.Context, idOrName string) (organization.Organization, error)
}

type Service struct {
	repository       Repository
	configRepository ConfigRepository
	relationService  RelationService
	userService      UserService
	projectService   ProjectService
	orgService       OrgService
}

func NewService(repository Repository, configRepository ConfigRepository,
	relationService RelationService, userService UserService,
	projectService ProjectService, orgService OrgService) *Service {
	return &Service{
		repository:       repository,
		configRepository: configRepository,
		relationService:  relationService,
		userService:      userService,
		projectService:   projectService,
		orgService:       orgService,
	}
}

func (s Service) Get(ctx context.Context, id string) (Resource, error) {
	return s.repository.GetByID(ctx, id)
}

func (s Service) Create(ctx context.Context, res Resource) (Resource, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Resource{}, err
	}

	userId := res.UserID
	if strings.TrimSpace(userId) == "" {
		userId = currentUser.ID
	}

	newResource, err := s.repository.Create(ctx, Resource{
		URN:         res.CreateURN(),
		Name:        res.Name,
		ProjectID:   res.ProjectID,
		NamespaceID: res.NamespaceID,
		UserID:      userId,
	})
	if err != nil {
		return Resource{}, err
	}

	if err = s.relationService.Delete(ctx, relation.Relation{
		Object: relation.Object{
			ID:        newResource.ID,
			Namespace: newResource.NamespaceID,
		},
	}); err != nil && !errors.Is(err, relation.ErrNotExist) {
		return Resource{}, err
	}

	if err = s.AddProjectToResource(ctx, newResource.ProjectID, newResource); err != nil {
		return Resource{}, err
	}
	if err = s.AddResourceOwner(ctx, newResource); err != nil {
		return Resource{}, err
	}

	return newResource, nil
}

func (s Service) List(ctx context.Context, flt Filter) ([]Resource, error) {
	return s.repository.List(ctx, flt)
}

func (s Service) Update(ctx context.Context, resource Resource) (Resource, error) {
	// TODO there should be an update logic like create here
	return s.repository.Update(ctx, resource)
}

func (s Service) AddProjectToResource(ctx context.Context, projectID string, res Resource) error {
	rel := relation.Relation{
		Object: relation.Object{
			ID:        res.ID,
			Namespace: res.NamespaceID,
		},
		Subject: relation.Subject{
			ID:        projectID,
			Namespace: schema.ProjectNamespace,
		},
		RelationName: schema.ProjectRelationName,
	}

	if _, err := s.relationService.Create(ctx, rel); err != nil {
		return err
	}
	return nil
}

func (s Service) AddResourceOwner(ctx context.Context, res Resource) error {
	if _, err := s.relationService.Create(ctx, relation.Relation{
		Object: relation.Object{
			ID:        res.ID,
			Namespace: res.NamespaceID,
		},
		Subject: relation.Subject{
			ID:        res.UserID,
			Namespace: schema.UserPrincipal,
		},
		RelationName: schema.OwnerRelationName,
	}); err != nil {
		return err
	}
	return nil
}

func (s Service) GetAllConfigs(ctx context.Context) ([]YAML, error) {
	return s.configRepository.GetAll(ctx)
}

func (s Service) CheckAuthz(ctx context.Context, rel relation.Object, permissionName string) (bool, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return false, err
	}

	// a user can pass object name instead of id in the request
	// we should convert name to id based on object namespace
	if !utils.IsValidUUID(rel.ID) {
		if rel.Namespace == schema.ProjectNamespace {
			// if object is project, then fetch project by name
			project, err := s.projectService.Get(ctx, rel.ID)
			if err != nil {
				return false, err
			}
			rel.ID = project.ID
		}
		if rel.Namespace == schema.OrganizationNamespace {
			// if object is org, then fetch org by name
			org, err := s.orgService.Get(ctx, rel.ID)
			if err != nil {
				return false, err
			}
			rel.ID = org.ID
		}
	}

	return s.relationService.CheckPermission(ctx, relation.Subject{
		ID: currentUser.ID,
		// TODO(kushsharma): refactor this to also support app/serviceuser
		Namespace: schema.UserPrincipal,
	}, rel, permissionName)
}

func (s Service) Delete(ctx context.Context, namespaceID, id string) error {
	if err := s.relationService.Delete(ctx, relation.Relation{
		Object: relation.Object{
			ID:        id,
			Namespace: namespaceID,
		},
	}); err != nil && !errors.Is(err, relation.ErrNotExist) {
		return err
	}
	return s.repository.Delete(ctx, id)
}
