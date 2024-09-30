// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	context "context"

	authenticate "github.com/raystack/frontier/core/authenticate"

	mock "github.com/stretchr/testify/mock"

	organization "github.com/raystack/frontier/core/organization"
)

// OrganizationService is an autogenerated mock type for the OrganizationService type
type OrganizationService struct {
	mock.Mock
}

type OrganizationService_Expecter struct {
	mock *mock.Mock
}

func (_m *OrganizationService) EXPECT() *OrganizationService_Expecter {
	return &OrganizationService_Expecter{mock: &_m.Mock}
}

// AddMember provides a mock function with given fields: ctx, orgID, relationName, principal
func (_m *OrganizationService) AddMember(ctx context.Context, orgID string, relationName string, principal authenticate.Principal) error {
	ret := _m.Called(ctx, orgID, relationName, principal)

	if len(ret) == 0 {
		panic("no return value specified for AddMember")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, authenticate.Principal) error); ok {
		r0 = rf(ctx, orgID, relationName, principal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrganizationService_AddMember_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddMember'
type OrganizationService_AddMember_Call struct {
	*mock.Call
}

// AddMember is a helper method to define mock.On call
//   - ctx context.Context
//   - orgID string
//   - relationName string
//   - principal authenticate.Principal
func (_e *OrganizationService_Expecter) AddMember(ctx interface{}, orgID interface{}, relationName interface{}, principal interface{}) *OrganizationService_AddMember_Call {
	return &OrganizationService_AddMember_Call{Call: _e.mock.On("AddMember", ctx, orgID, relationName, principal)}
}

func (_c *OrganizationService_AddMember_Call) Run(run func(ctx context.Context, orgID string, relationName string, principal authenticate.Principal)) *OrganizationService_AddMember_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(authenticate.Principal))
	})
	return _c
}

func (_c *OrganizationService_AddMember_Call) Return(_a0 error) *OrganizationService_AddMember_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrganizationService_AddMember_Call) RunAndReturn(run func(context.Context, string, string, authenticate.Principal) error) *OrganizationService_AddMember_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *OrganizationService) Get(ctx context.Context, id string) (organization.Organization, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 organization.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (organization.Organization, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) organization.Organization); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(organization.Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrganizationService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type OrganizationService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *OrganizationService_Expecter) Get(ctx interface{}, id interface{}) *OrganizationService_Get_Call {
	return &OrganizationService_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *OrganizationService_Get_Call) Run(run func(ctx context.Context, id string)) *OrganizationService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *OrganizationService_Get_Call) Return(_a0 organization.Organization, _a1 error) *OrganizationService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrganizationService_Get_Call) RunAndReturn(run func(context.Context, string) (organization.Organization, error)) *OrganizationService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// ListByUser provides a mock function with given fields: ctx, userID, f
func (_m *OrganizationService) ListByUser(ctx context.Context, userID string, f organization.Filter) ([]organization.Organization, error) {
	ret := _m.Called(ctx, userID, f)

	if len(ret) == 0 {
		panic("no return value specified for ListByUser")
	}

	var r0 []organization.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, organization.Filter) ([]organization.Organization, error)); ok {
		return rf(ctx, userID, f)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, organization.Filter) []organization.Organization); ok {
		r0 = rf(ctx, userID, f)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]organization.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, organization.Filter) error); ok {
		r1 = rf(ctx, userID, f)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrganizationService_ListByUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListByUser'
type OrganizationService_ListByUser_Call struct {
	*mock.Call
}

// ListByUser is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - f organization.Filter
func (_e *OrganizationService_Expecter) ListByUser(ctx interface{}, userID interface{}, f interface{}) *OrganizationService_ListByUser_Call {
	return &OrganizationService_ListByUser_Call{Call: _e.mock.On("ListByUser", ctx, userID, f)}
}

func (_c *OrganizationService_ListByUser_Call) Run(run func(ctx context.Context, userID string, f organization.Filter)) *OrganizationService_ListByUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(organization.Filter))
	})
	return _c
}

func (_c *OrganizationService_ListByUser_Call) Return(_a0 []organization.Organization, _a1 error) *OrganizationService_ListByUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrganizationService_ListByUser_Call) RunAndReturn(run func(context.Context, string, organization.Filter) ([]organization.Organization, error)) *OrganizationService_ListByUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewOrganizationService creates a new instance of OrganizationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrganizationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrganizationService {
	mock := &OrganizationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
