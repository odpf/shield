package v1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/odpf/shield/model"
	shieldv1 "go.buf.build/odpf/gwv/odpf/proton/odpf/shield/v1"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var testNsMap = map[string]model.Namespace{
	"team": {
		Id:        "team",
		Name:      "Team",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	"org": {
		Id:        "org",
		Name:      "Org",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	"project": {
		Id:        "project",
		Name:      "Project",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
}

func TestListNamespaces(t *testing.T) {
	t.Parallel()
	table := []struct {
		title            string
		mockNamespaceSrv mockNamespaceSrv
		req              *shieldv1.ListNamespacesRequest
		want             *shieldv1.ListNamespacesResponse
		err              error
	}{
		{
			title: "error in Namespace Service",
			mockNamespaceSrv: mockNamespaceSrv{ListNamespacesFunc: func(ctx context.Context) (namespaces []model.Namespace, err error) {
				return []model.Namespace{}, errors.New("some error")
			}},
			want: nil,
			err:  status.Errorf(codes.Internal, internalServerError.Error()),
		},
		{
			title: "success",
			mockNamespaceSrv: mockNamespaceSrv{ListNamespacesFunc: func(ctx context.Context) (namespaces []model.Namespace, err error) {
				var testNSList []model.Namespace
				for _, ns := range testNsMap {
					testNSList = append(testNSList, ns)
				}
				return testNSList, nil
			}},
			want: &shieldv1.ListNamespacesResponse{Namespaces: []*shieldv1.Namespace{
				{
					Id:        "team",
					Name:      "Team",
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
				{
					Id:        "org",
					Name:      "Org",
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
				{
					Id:        "project",
					Name:      "Project",
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			t.Parallel()

			mockDep := Dep{NamespaceService: tt.mockNamespaceSrv}
			resp, err := mockDep.ListNamespaces(context.Background(), tt.req)

			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

func TestCreateNamespace(t *testing.T) {
	t.Parallel()

	table := []struct {
		title            string
		mockNamespaceSrv mockNamespaceSrv
		req              *shieldv1.CreateNamespaceRequest
		want             *shieldv1.CreateNamespaceResponse
		err              error
	}{
		{
			title: "error in creating namespace",
			mockNamespaceSrv: mockNamespaceSrv{CreateNamespaceFunc: func(ctx context.Context, ns model.Namespace) (model.Namespace, error) {
				return model.Namespace{}, errors.New("some error")
			}},
			req: &shieldv1.CreateNamespaceRequest{Body: &shieldv1.NamespaceRequestBody{
				Id:   "team",
				Name: "Team",
			}},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "success",
			mockNamespaceSrv: mockNamespaceSrv{CreateNamespaceFunc: func(ctx context.Context, ns model.Namespace) (model.Namespace, error) {
				return model.Namespace{
					Id:   "team",
					Name: "Team",
				}, nil
			}},
			req: &shieldv1.CreateNamespaceRequest{Body: &shieldv1.NamespaceRequestBody{
				Id:   "team",
				Name: "Team",
			}},
			want: &shieldv1.CreateNamespaceResponse{Namespace: &shieldv1.Namespace{
				Id:        "team",
				Name:      "Team",
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			t.Parallel()

			mockDep := Dep{NamespaceService: tt.mockNamespaceSrv}
			resp, err := mockDep.CreateNamespace(context.Background(), tt.req)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

type mockNamespaceSrv struct {
	GetNamespaceFunc    func(ctx context.Context, id string) (model.Namespace, error)
	CreateNamespaceFunc func(ctx context.Context, ns model.Namespace) (model.Namespace, error)
	ListNamespacesFunc  func(ctx context.Context) ([]model.Namespace, error)
	UpdateNamespaceFunc func(ctx context.Context, id string, ns model.Namespace) (model.Namespace, error)
}

func (m mockNamespaceSrv) GetNamespace(ctx context.Context, id string) (model.Namespace, error) {
	return m.GetNamespaceFunc(ctx, id)
}

func (m mockNamespaceSrv) ListNamespaces(ctx context.Context) ([]model.Namespace, error) {
	return m.ListNamespacesFunc(ctx)
}

func (m mockNamespaceSrv) CreateNamespace(ctx context.Context, ns model.Namespace) (model.Namespace, error) {
	return m.CreateNamespaceFunc(ctx, ns)
}

func (m mockNamespaceSrv) UpdateNamespace(ctx context.Context, id string, ns model.Namespace) (model.Namespace, error) {
	return m.UpdateNamespaceFunc(ctx, id, ns)
}
