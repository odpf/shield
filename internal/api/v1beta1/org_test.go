// Package v1beta1 provides v1beta1  
package v1beta1

import (
	"context"
	"testing"
	"time"

	"github.com/raystack/shield/core/authenticate"

	"github.com/raystack/shield/pkg/utils"

	"github.com/raystack/shield/internal/bootstrap/schema"

	"github.com/raystack/shield/core/organization"
	"github.com/raystack/shield/core/user"
	"github.com/raystack/shield/internal/api/v1beta1/mocks"
	"github.com/raystack/shield/pkg/errors"
	"github.com/raystack/shield/pkg/metadata"
	shieldv1beta1 "github.com/raystack/shield/proto/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	testOrgID  = "9f256f86-31a3-11ec-8d3d-0242ac130003"
	testOrgMap = map[string]organization.Organization{
		"9f256f86-31a3-11ec-8d3d-0242ac130003": {
			ID:   "9f256f86-31a3-11ec-8d3d-0242ac130003",
			Name: "org-1",
			Metadata: metadata.Metadata{
				"email":  "org1@org1.com",
				"age":    21,
				"intern": true,
			},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}
)

func TestListOrganizations(t *testing.T) {
	table := []struct {
		title string
		setup func(os *mocks.OrganizationService)
		want  *shieldv1beta1.ListOrganizationsResponse
		err   error
	}{
		{
			title: "should return internal error if org service return some error",
			setup: func(os *mocks.OrganizationService) {
				os.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), organization.Filter{}).Return([]organization.Organization{}, errors.New("some error"))
			},
			want: nil,
			err:  status.Errorf(codes.Internal, ErrInternalServer.Error()),
		}, {
			title: "should return success if org service return nil",
			setup: func(os *mocks.OrganizationService) {
				var testOrgList []organization.Organization
				for _, o := range testOrgMap {
					testOrgList = append(testOrgList, o)
				}
				os.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), organization.Filter{}).Return(testOrgList, nil)
			},
			want: &shieldv1beta1.ListOrganizationsResponse{Organizations: []*shieldv1beta1.Organization{
				{
					Id:   "9f256f86-31a3-11ec-8d3d-0242ac130003",
					Name: "org-1",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email":  structpb.NewStringValue("org1@org1.com"),
							"age":    structpb.NewNumberValue(21),
							"intern": structpb.NewBoolValue(true),
						},
					},
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockOrgSrv := new(mocks.OrganizationService)
			mockDep := Handler{orgService: mockOrgSrv}
			if tt.setup != nil {
				tt.setup(mockOrgSrv)
			}
			resp, err := mockDep.ListOrganizations(context.Background(), nil)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestCreateOrganization(t *testing.T) {
	email := "user@raystack.org"
	table := []struct {
		title string
		setup func(ctx context.Context, os *mocks.OrganizationService, ms *mocks.MetaSchemaService) context.Context
		req   *shieldv1beta1.CreateOrganizationRequest
		want  *shieldv1beta1.CreateOrganizationResponse
		err   error
	}{
		{
			title: "should return forbidden error if auth email in context is empty and org service return invalid user email",
			setup: func(ctx context.Context, os *mocks.OrganizationService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), organization.Organization{
					Name:     "some-org",
					Metadata: metadata.Metadata{},
				}).Return(organization.Organization{}, user.ErrInvalidEmail)
				return ctx
			},
			req: &shieldv1beta1.CreateOrganizationRequest{Body: &shieldv1beta1.OrganizationRequestBody{
				Name:     "some-org",
				Metadata: &structpb.Struct{},
			}},
			want: nil,
			err:  grpcUnauthenticated,
		},
		{
			title: "should return internal error if org service return some error",
			setup: func(ctx context.Context, os *mocks.OrganizationService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), organization.Organization{
					Name:     "abc",
					Metadata: metadata.Metadata{},
				}).Return(organization.Organization{}, errors.New("some error"))
				return authenticate.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateOrganizationRequest{Body: &shieldv1beta1.OrganizationRequestBody{
				Name:     "abc",
				Metadata: &structpb.Struct{},
			}},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "should return bad request error if name is empty",
			setup: func(ctx context.Context, os *mocks.OrganizationService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), organization.Organization{
					Name:     "abc",
					Metadata: metadata.Metadata{},
				}).Return(organization.Organization{}, organization.ErrInvalidDetail)
				return authenticate.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateOrganizationRequest{Body: &shieldv1beta1.OrganizationRequestBody{
				Name:     "abc",
				Metadata: &structpb.Struct{},
			}},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return already exist error if org service return error conflict",
			setup: func(ctx context.Context, os *mocks.OrganizationService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), organization.Organization{
					Name:     "abc",
					Metadata: metadata.Metadata{},
				}).Return(organization.Organization{}, organization.ErrConflict)
				return authenticate.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateOrganizationRequest{Body: &shieldv1beta1.OrganizationRequestBody{
				Name:     "abc",
				Metadata: &structpb.Struct{},
			}},
			want: nil,
			err:  grpcConflictError,
		},
		{
			title: "should return bad request error if metadata is not parsable",
			req: &shieldv1beta1.CreateOrganizationRequest{Body: &shieldv1beta1.OrganizationRequestBody{
				Title: "some-org",
				Name:  "abc",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"count": structpb.NewNullValue(),
					},
				},
			}},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return success if org service return nil error",
			setup: func(ctx context.Context, os *mocks.OrganizationService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), organization.Organization{
					Name: "some-org",
					Metadata: metadata.Metadata{
						"email": "a",
					},
				}).Return(organization.Organization{
					ID:   "new-abc",
					Name: "some-org",
					Metadata: metadata.Metadata{
						"email": "a",
					},
				}, nil)
				return authenticate.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateOrganizationRequest{Body: &shieldv1beta1.OrganizationRequestBody{
				Name: "some-org",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"email": structpb.NewStringValue("a"),
					},
				},
			}},
			want: &shieldv1beta1.CreateOrganizationResponse{Organization: &shieldv1beta1.Organization{
				Id:   "new-abc",
				Name: "some-org",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"email": structpb.NewStringValue("a"),
					}},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockOrgSrv := new(mocks.OrganizationService)
			mockMetaSchemaSvc := new(mocks.MetaSchemaService)
			ctx := context.Background()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockOrgSrv, mockMetaSchemaSvc)
			}
			mockDep := Handler{orgService: mockOrgSrv, metaSchemaService: mockMetaSchemaSvc}
			resp, err := mockDep.CreateOrganization(ctx, tt.req)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

func TestHandler_GetOrganization(t *testing.T) {
	someOrgID := utils.NewString()
	tests := []struct {
		name    string
		setup   func(os *mocks.OrganizationService)
		request *shieldv1beta1.GetOrganizationRequest
		want    *shieldv1beta1.GetOrganizationResponse
		wantErr error
	}{

		{
			name: "should return internal error if org service return some error",
			setup: func(os *mocks.OrganizationService) {
				os.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), someOrgID).Return(organization.Organization{}, errors.New("some error"))
			},
			request: &shieldv1beta1.GetOrganizationRequest{
				Id: someOrgID,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return not found error if org id is not uuid (slug) and org not exist",
			setup: func(os *mocks.OrganizationService) {
				os.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), someOrgID).Return(organization.Organization{}, organization.ErrNotExist)
			},
			request: &shieldv1beta1.GetOrganizationRequest{
				Id: someOrgID,
			},
			want:    nil,
			wantErr: grpcOrgNotFoundErr,
		},
		{
			name: "should return not found error if org id is invalid",
			setup: func(os *mocks.OrganizationService) {
				os.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), "").Return(organization.Organization{}, organization.ErrInvalidID)
			},
			request: &shieldv1beta1.GetOrganizationRequest{},
			want:    nil,
			wantErr: grpcOrgNotFoundErr,
		},
		{
			name: "should return success if org service return nil error",
			setup: func(os *mocks.OrganizationService) {
				os.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), "9f256f86-31a3-11ec-8d3d-0242ac130003").Return(testOrgMap["9f256f86-31a3-11ec-8d3d-0242ac130003"], nil)
			},
			request: &shieldv1beta1.GetOrganizationRequest{
				Id: "9f256f86-31a3-11ec-8d3d-0242ac130003",
			},
			want: &shieldv1beta1.GetOrganizationResponse{
				Organization: &shieldv1beta1.Organization{
					Id:   "9f256f86-31a3-11ec-8d3d-0242ac130003",
					Name: "org-1",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email":  structpb.NewStringValue("org1@org1.com"),
							"age":    structpb.NewNumberValue(21),
							"intern": structpb.NewBoolValue(true),
						},
					},
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrgSrv := new(mocks.OrganizationService)
			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(mockOrgSrv)
			}
			mockDep := Handler{orgService: mockOrgSrv}
			got, err := mockDep.GetOrganization(ctx, tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_UpdateOrganization(t *testing.T) {
	someOrgID := utils.NewString()
	tests := []struct {
		name    string
		setup   func(os *mocks.OrganizationService, ms *mocks.MetaSchemaService)
		request *shieldv1beta1.UpdateOrganizationRequest
		want    *shieldv1beta1.UpdateOrganizationResponse
		wantErr error
	}{
		{
			name: "should return internal error if org service return some error",
			setup: func(os *mocks.OrganizationService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), organization.Organization{
					ID: someOrgID,
					Metadata: metadata.Metadata{
						"email": "org1@org1.com",
						"age":   float64(21),
						"valid": true,
					},
					Name: "new-org",
				}).Return(organization.Organization{}, errors.New("some error"))
			},
			request: &shieldv1beta1.UpdateOrganizationRequest{
				Id: someOrgID,
				Body: &shieldv1beta1.OrganizationRequestBody{
					Name: "new-org",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
							"age":   structpb.NewNumberValue(21),
							"valid": structpb.NewBoolValue(true),
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return not found error if org id is not uuid (slug) and not exist",
			setup: func(os *mocks.OrganizationService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), organization.Organization{
					ID: someOrgID,
					Metadata: metadata.Metadata{
						"email": "org1@org1.com",
						"age":   float64(21),
						"valid": true,
					},
					Name: "new-org",
				}).Return(organization.Organization{}, organization.ErrNotExist)
			},
			request: &shieldv1beta1.UpdateOrganizationRequest{
				Id: someOrgID,
				Body: &shieldv1beta1.OrganizationRequestBody{
					Name: "new-org",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
							"age":   structpb.NewNumberValue(21),
							"valid": structpb.NewBoolValue(true),
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcOrgNotFoundErr,
		},
		{
			name: "should return not found error if org id is empty",
			setup: func(os *mocks.OrganizationService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), organization.Organization{
					Name: "new-org",
					Metadata: metadata.Metadata{
						"email": "org1@org1.com",
						"age":   float64(21),
						"valid": true,
					},
				}).Return(organization.Organization{}, organization.ErrInvalidID)
			},
			request: &shieldv1beta1.UpdateOrganizationRequest{
				Body: &shieldv1beta1.OrganizationRequestBody{
					Name: "new-org",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
							"age":   structpb.NewNumberValue(21),
							"valid": structpb.NewBoolValue(true),
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcOrgNotFoundErr,
		},
		{
			name: "should return already exist error if org service return err conflict",
			setup: func(os *mocks.OrganizationService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), organization.Organization{
					ID: someOrgID,
					Metadata: metadata.Metadata{
						"email": "org1@org1.com",
						"age":   float64(21),
						"valid": true,
					},
					Name: "new-org",
				}).Return(organization.Organization{}, organization.ErrConflict)
			},
			request: &shieldv1beta1.UpdateOrganizationRequest{
				Id: someOrgID,
				Body: &shieldv1beta1.OrganizationRequestBody{
					Name: "new-org",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
							"age":   structpb.NewNumberValue(21),
							"valid": structpb.NewBoolValue(true),
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcConflictError,
		},
		{
			name: "should return success if org service is updated by id and return nil error",
			setup: func(os *mocks.OrganizationService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), organization.Organization{
					ID: someOrgID,
					Metadata: metadata.Metadata{
						"email": "org1@org1.com",
						"age":   float64(21),
						"valid": true,
					},
					Name: "new-org",
				}).Return(organization.Organization{
					ID: someOrgID,
					Metadata: metadata.Metadata{
						"email": "org1@org1.com",
						"age":   float64(21),
						"valid": true,
					},
					Name: "new-org",
				}, nil)
			},
			request: &shieldv1beta1.UpdateOrganizationRequest{
				Id: someOrgID,
				Body: &shieldv1beta1.OrganizationRequestBody{
					Name: "new-org",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
							"age":   structpb.NewNumberValue(21),
							"valid": structpb.NewBoolValue(true),
						},
					},
				},
			},
			want: &shieldv1beta1.UpdateOrganizationResponse{
				Organization: &shieldv1beta1.Organization{
					Id:   someOrgID,
					Name: "new-org",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
							"age":   structpb.NewNumberValue(21),
							"valid": structpb.NewBoolValue(true),
						},
					},
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
			},
			wantErr: nil,
		},
		{
			name: "should return success if org service is updated by name and return nil error",
			setup: func(os *mocks.OrganizationService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), orgMetaSchema).Return(nil)
				os.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), organization.Organization{
					Name: "new-org",
					Metadata: metadata.Metadata{
						"email": "org1@org1.com",
						"age":   float64(21),
						"valid": true,
					},
				}).Return(organization.Organization{
					ID:   someOrgID,
					Name: "new-org",
					Metadata: metadata.Metadata{
						"email": "org1@org1.com",
						"age":   float64(21),
						"valid": true,
					},
				}, nil)
			},
			request: &shieldv1beta1.UpdateOrganizationRequest{
				Body: &shieldv1beta1.OrganizationRequestBody{
					Name: "new-org",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
							"age":   structpb.NewNumberValue(21),
							"valid": structpb.NewBoolValue(true),
						},
					},
				},
			},
			want: &shieldv1beta1.UpdateOrganizationResponse{
				Organization: &shieldv1beta1.Organization{
					Id:   someOrgID,
					Name: "new-org",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
							"age":   structpb.NewNumberValue(21),
							"valid": structpb.NewBoolValue(true),
						},
					},
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrgSrv := new(mocks.OrganizationService)
			mockMetaSchemaSvc := new(mocks.MetaSchemaService)
			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(mockOrgSrv, mockMetaSchemaSvc)
			}
			mockDep := Handler{orgService: mockOrgSrv, metaSchemaService: mockMetaSchemaSvc}
			got, err := mockDep.UpdateOrganization(ctx, tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_ListOrganizationAdmins(t *testing.T) {
	someOrgID := utils.NewString()
	tests := []struct {
		name    string
		setup   func(us *mocks.UserService)
		request *shieldv1beta1.ListOrganizationAdminsRequest
		want    *shieldv1beta1.ListOrganizationAdminsResponse
		wantErr error
	}{
		{
			name: "should return internal error if org service return some error",
			setup: func(us *mocks.UserService) {
				us.EXPECT().ListByOrg(mock.AnythingOfType("*context.emptyCtx"), someOrgID, schema.UpdatePermission).Return([]user.User{}, errors.New("some error"))
			},
			request: &shieldv1beta1.ListOrganizationAdminsRequest{
				Id: someOrgID,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return empty list of users if org id is not exist",
			setup: func(us *mocks.UserService) {
				us.EXPECT().ListByOrg(mock.AnythingOfType("*context.emptyCtx"), someOrgID, schema.UpdatePermission).Return([]user.User{}, nil)
			},
			request: &shieldv1beta1.ListOrganizationAdminsRequest{
				Id: someOrgID,
			},
			want:    &shieldv1beta1.ListOrganizationAdminsResponse{},
			wantErr: nil,
		},
		{
			name: "should return success if org service return nil error",
			setup: func(us *mocks.UserService) {
				var testUserList []user.User
				for _, u := range testUserMap {
					testUserList = append(testUserList, u)
				}
				us.EXPECT().ListByOrg(mock.AnythingOfType("*context.emptyCtx"), someOrgID, schema.UpdatePermission).Return(testUserList, nil)
			},
			request: &shieldv1beta1.ListOrganizationAdminsRequest{
				Id: someOrgID,
			},
			want: &shieldv1beta1.ListOrganizationAdminsResponse{
				Users: []*shieldv1beta1.User{
					{
						Id:    "9f256f86-31a3-11ec-8d3d-0242ac130003",
						Title: "User 1",
						Name:  "user1",
						Email: "test@test.com",
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"foo":    structpb.NewStringValue("bar"),
								"age":    structpb.NewNumberValue(21),
								"intern": structpb.NewBoolValue(true),
							},
						},
						CreatedAt: timestamppb.New(time.Time{}),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(mocks.UserService)
			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(mockUserService)
			}
			mockDep := Handler{userService: mockUserService}
			got, err := mockDep.ListOrganizationAdmins(ctx, tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}
