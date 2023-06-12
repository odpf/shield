// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	resource "github.com/raystack/shield/core/resource"
	mock "github.com/stretchr/testify/mock"
)

// ResourceService is an autogenerated mock type for the ResourceService type
type ResourceService struct {
	mock.Mock
}

type ResourceService_Expecter struct {
	mock *mock.Mock
}

func (_m *ResourceService) EXPECT() *ResourceService_Expecter {
	return &ResourceService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *ResourceService) Create(ctx context.Context, _a1 resource.Resource) (resource.Resource, error) {
	ret := _m.Called(ctx, _a1)

	var r0 resource.Resource
	if rf, ok := ret.Get(0).(func(context.Context, resource.Resource) resource.Resource); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(resource.Resource)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, resource.Resource) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type ResourceService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//  - ctx context.Context
//  - _a1 resource.Resource
func (_e *ResourceService_Expecter) Create(ctx interface{}, _a1 interface{}) *ResourceService_Create_Call {
	return &ResourceService_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *ResourceService_Create_Call) Run(run func(ctx context.Context, _a1 resource.Resource)) *ResourceService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(resource.Resource))
	})
	return _c
}

func (_c *ResourceService_Create_Call) Return(_a0 resource.Resource, _a1 error) *ResourceService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewResourceService interface {
	mock.TestingT
	Cleanup(func())
}

// NewResourceService creates a new instance of ResourceService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewResourceService(t mockConstructorTestingTNewResourceService) *ResourceService {
	mock := &ResourceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
