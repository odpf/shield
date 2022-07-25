package v1beta1

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/odpf/shield/core/namespace"
	shieldv1beta1 "github.com/odpf/shield/proto/v1beta1"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var testNsMap = map[string]namespace.Namespace{
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
		req              *shieldv1beta1.ListNamespacesRequest
		want             *shieldv1beta1.ListNamespacesResponse
		err              error
	}{
		{
			title: "error in Namespace Service",
			mockNamespaceSrv: mockNamespaceSrv{ListNamespacesFunc: func(ctx context.Context) (namespaces []namespace.Namespace, err error) {
				return []namespace.Namespace{}, errors.New("some error")
			}},
			want: nil,
			err:  status.Errorf(codes.Internal, internalServerError.Error()),
		},
		{
			title: "success",
			mockNamespaceSrv: mockNamespaceSrv{ListNamespacesFunc: func(ctx context.Context) ([]namespace.Namespace, error) {
				var testNSList []namespace.Namespace
				for _, ns := range testNsMap {
					testNSList = append(testNSList, ns)
				}
				sort.Slice(testNSList[:], func(i, j int) bool {
					return strings.Compare(testNSList[i].Id, testNSList[j].Id) < 1
				})
				return testNSList, nil
			}},
			want: &shieldv1beta1.ListNamespacesResponse{Namespaces: []*shieldv1beta1.Namespace{
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
				{
					Id:        "team",
					Name:      "Team",
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
		req              *shieldv1beta1.CreateNamespaceRequest
		want             *shieldv1beta1.CreateNamespaceResponse
		err              error
	}{
		{
			title: "error in creating namespace",
			mockNamespaceSrv: mockNamespaceSrv{CreateNamespaceFunc: func(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
				return namespace.Namespace{}, errors.New("some error")
			}},
			req: &shieldv1beta1.CreateNamespaceRequest{Body: &shieldv1beta1.NamespaceRequestBody{
				Id:   "team",
				Name: "Team",
			}},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "success",
			mockNamespaceSrv: mockNamespaceSrv{CreateNamespaceFunc: func(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
				return namespace.Namespace{
					Id:   "team",
					Name: "Team",
				}, nil
			}},
			req: &shieldv1beta1.CreateNamespaceRequest{Body: &shieldv1beta1.NamespaceRequestBody{
				Id:   "team",
				Name: "Team",
			}},
			want: &shieldv1beta1.CreateNamespaceResponse{Namespace: &shieldv1beta1.Namespace{
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
	GetNamespaceFunc    func(ctx context.Context, id string) (namespace.Namespace, error)
	CreateNamespaceFunc func(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error)
	ListNamespacesFunc  func(ctx context.Context) ([]namespace.Namespace, error)
	UpdateNamespaceFunc func(ctx context.Context, id string, ns namespace.Namespace) (namespace.Namespace, error)
}

func (m mockNamespaceSrv) GetNamespace(ctx context.Context, id string) (namespace.Namespace, error) {
	return m.GetNamespaceFunc(ctx, id)
}

func (m mockNamespaceSrv) ListNamespaces(ctx context.Context) ([]namespace.Namespace, error) {
	return m.ListNamespacesFunc(ctx)
}

func (m mockNamespaceSrv) CreateNamespace(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
	return m.CreateNamespaceFunc(ctx, ns)
}

func (m mockNamespaceSrv) UpdateNamespace(ctx context.Context, id string, ns namespace.Namespace) (namespace.Namespace, error) {
	return m.UpdateNamespaceFunc(ctx, id, ns)
}
