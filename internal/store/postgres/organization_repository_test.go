package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/odpf/salt/log"
	"github.com/odpf/shield/core/organization"
	"github.com/odpf/shield/core/relation"
	"github.com/odpf/shield/core/user"
	"github.com/odpf/shield/internal/schema"
	"github.com/odpf/shield/internal/store/postgres"
	"github.com/odpf/shield/pkg/db"
	"github.com/odpf/shield/pkg/metadata"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
)

type OrganizationRepositoryTestSuite struct {
	suite.Suite
	ctx                 context.Context
	client              *db.Client
	pool                *dockertest.Pool
	resource            *dockertest.Resource
	repository          *postgres.OrganizationRepository
	relationRepository  *postgres.RelationRepository
	namespaceRepository *postgres.NamespaceRepository
	roleRepository      *postgres.RoleRepository
	orgs                []organization.Organization
	users               []user.User
}

func (s *OrganizationRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewOrganizationRepository(s.client)

	s.users, err = bootstrapUser(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	s.relationRepository = postgres.NewRelationRepository(s.client)
	s.namespaceRepository = postgres.NewNamespaceRepository(s.client)
	s.roleRepository = postgres.NewRoleRepository(s.client)
}

func (s *OrganizationRepositoryTestSuite) SetupTest() {
	var err error
	s.orgs, err = bootstrapOrganization(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapAction(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapRole(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.relationRepository.Create(context.Background(), relation.RelationV2{
		Subject: relation.Subject{
			ID:        s.users[0].ID,
			Namespace: schema.UserPrincipal,
			RoleID:    schema.OwnerRole,
		},
		Object: relation.Object{
			ID:        s.orgs[0].ID,
			Namespace: schema.OrganizationNamespace,
		},
	})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *OrganizationRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *OrganizationRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *OrganizationRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_ORGANIZATIONS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_RELATIONS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_ROLES),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_NAMESPACES),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *OrganizationRepositoryTestSuite) TestGetByID() {
	type testCase struct {
		Description          string
		SelectedID           string
		ExpectedOrganization organization.Organization
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should get an organization",
			SelectedID:  s.orgs[0].ID,
			ExpectedOrganization: organization.Organization{
				Name:  "org1",
				Slug:  "org-1",
				State: organization.Enabled,
			},
		},
		{
			Description: "should return error no exist if can't found organization",
			SelectedID:  uuid.NewString(),
			ErrString:   organization.ErrNotExist.Error(),
		},
		{
			Description: "should return error if id empty",
			ErrString:   organization.ErrInvalidID.Error(),
		},
		{
			Description: "should return error if id is not uuid",
			SelectedID:  "10000",
			ErrString:   organization.ErrInvalidUUID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetByID(s.ctx, tc.SelectedID)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedOrganization, cmpopts.IgnoreFields(organization.Organization{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedOrganization)
			}
		})
	}
}

func (s *OrganizationRepositoryTestSuite) TestGetBySlug() {
	type testCase struct {
		Description          string
		SelectedSlug         string
		ExpectedOrganization organization.Organization
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description:  "should get an organization",
			SelectedSlug: "org-1",
			ExpectedOrganization: organization.Organization{
				Name:  "org1",
				Slug:  "org-1",
				State: organization.Enabled,
			},
		},
		{
			Description:  "should return error no exist if can't found organization",
			SelectedSlug: "randomslug",
			ErrString:    organization.ErrNotExist.Error(),
		},
		{
			Description: "should return error if slug empty",
			ErrString:   organization.ErrInvalidID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetBySlug(s.ctx, tc.SelectedSlug)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedOrganization, cmpopts.IgnoreFields(organization.Organization{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedOrganization)
			}
		})
	}
}

func (s *OrganizationRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description          string
		OrganizationToCreate organization.Organization
		ExpectedOrganization organization.Organization
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should create an organization",
			OrganizationToCreate: organization.Organization{
				Name:     "new-org",
				Slug:     "new-org-slug",
				Metadata: metadata.Metadata{},
			},
			ExpectedOrganization: organization.Organization{
				Name:     "new-org",
				Slug:     "new-org-slug",
				State:    organization.Enabled,
				Metadata: metadata.Metadata{},
			},
		},
		{
			Description: "should return error if organization slug already exist",
			OrganizationToCreate: organization.Organization{
				Name:     "newslug",
				Slug:     "org-1",
				Metadata: metadata.Metadata{},
			},
			ErrString: organization.ErrConflict.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Create(s.ctx, tc.OrganizationToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedOrganization, cmpopts.IgnoreFields(organization.Organization{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedOrganization)
			}
		})
	}
}

func (s *OrganizationRepositoryTestSuite) TestList() {
	type testCase struct {
		Description           string
		ExpectedOrganizations []organization.Organization
		ErrString             string
	}

	var testCases = []testCase{
		{
			Description: "should get all organizations",
			ExpectedOrganizations: []organization.Organization{
				{
					Name:     "org1",
					Slug:     "org-1",
					State:    organization.Enabled,
					Metadata: metadata.Metadata{},
				},
				{
					Name:     "org2",
					Slug:     "org-2",
					State:    organization.Enabled,
					Metadata: metadata.Metadata{},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.List(s.ctx, organization.Filter{})
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedOrganizations, cmpopts.IgnoreFields(organization.Organization{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedOrganizations)
			}
		})
	}
}

func (s *OrganizationRepositoryTestSuite) TestUpdateByID() {
	type testCase struct {
		Description          string
		OrganizationToUpdate organization.Organization
		ExpectedOrganization organization.Organization
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should update a organization",
			OrganizationToUpdate: organization.Organization{
				ID:       s.orgs[0].ID,
				Name:     "new org update",
				Slug:     "new-org-update",
				Metadata: metadata.Metadata{},
			},
			ExpectedOrganization: organization.Organization{
				Name:     "new org update",
				Slug:     "new-org-update",
				State:    organization.Enabled,
				Metadata: metadata.Metadata{},
			},
		},
		{
			Description: "should return error if organization slug already exist",
			OrganizationToUpdate: organization.Organization{
				ID:       s.orgs[0].ID,
				Name:     "new-org-2",
				Slug:     "org-2",
				Metadata: metadata.Metadata{},
			},
			ErrString: organization.ErrConflict.Error(),
		},
		{
			Description: "should return error if organization not found",
			OrganizationToUpdate: organization.Organization{
				ID:       uuid.NewString(),
				Name:     "not-exist",
				Slug:     "some-slug",
				Metadata: metadata.Metadata{},
			},
			ErrString: organization.ErrNotExist.Error(),
		},
		{
			Description: "should return error if organization id is not uuid",
			OrganizationToUpdate: organization.Organization{
				ID:       "12345",
				Name:     "not-exist",
				Slug:     "some-slug",
				Metadata: metadata.Metadata{},
			},
			ErrString: organization.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if organization id is empty",
			ErrString:   organization.ErrInvalidID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.UpdateByID(s.ctx, tc.OrganizationToUpdate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedOrganization, cmpopts.IgnoreFields(organization.Organization{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedOrganization)
			}
		})
	}
}

func (s *OrganizationRepositoryTestSuite) TestUpdateBySlug() {
	type testCase struct {
		Description          string
		OrganizationToUpdate organization.Organization
		ExpectedOrganization organization.Organization
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should update a organization",
			OrganizationToUpdate: organization.Organization{
				Slug:     "org-1",
				Name:     "new org update",
				Metadata: metadata.Metadata{},
			},
			ExpectedOrganization: organization.Organization{
				Name:     "new org update",
				Slug:     "org-1",
				State:    organization.Enabled,
				Metadata: metadata.Metadata{},
			},
		},
		{
			Description: "should return error if organization not found",
			OrganizationToUpdate: organization.Organization{
				Slug:     "slug",
				Name:     "not-exist",
				Metadata: metadata.Metadata{},
			},
			ErrString: organization.ErrNotExist.Error(),
		},
		{
			Description: "should return error if organization slug is empty",
			ErrString:   organization.ErrInvalidID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.UpdateBySlug(s.ctx, tc.OrganizationToUpdate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedOrganization, cmpopts.IgnoreFields(organization.Organization{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedOrganization)
			}
		})
	}
}

func TestOrganizationRepository(t *testing.T) {
	suite.Run(t, new(OrganizationRepositoryTestSuite))
}
