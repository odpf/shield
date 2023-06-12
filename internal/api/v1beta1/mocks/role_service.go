// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	role "github.com/raystack/shield/core/role"
	mock "github.com/stretchr/testify/mock"
)

// RoleService is an autogenerated mock type for the RoleService type
type RoleService struct {
	mock.Mock
}

type RoleService_Expecter struct {
	mock *mock.Mock
}

func (_m *RoleService) EXPECT() *RoleService_Expecter {
	return &RoleService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, toCreate
func (_m *RoleService) Create(ctx context.Context, toCreate role.Role) (role.Role, error) {
	ret := _m.Called(ctx, toCreate)

	var r0 role.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, role.Role) (role.Role, error)); ok {
		return rf(ctx, toCreate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, role.Role) role.Role); ok {
		r0 = rf(ctx, toCreate)
	} else {
		r0 = ret.Get(0).(role.Role)
	}

	if rf, ok := ret.Get(1).(func(context.Context, role.Role) error); ok {
		r1 = rf(ctx, toCreate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type RoleService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - toCreate role.Role
func (_e *RoleService_Expecter) Create(ctx interface{}, toCreate interface{}) *RoleService_Create_Call {
	return &RoleService_Create_Call{Call: _e.mock.On("Create", ctx, toCreate)}
}

func (_c *RoleService_Create_Call) Run(run func(ctx context.Context, toCreate role.Role)) *RoleService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(role.Role))
	})
	return _c
}

func (_c *RoleService_Create_Call) Return(_a0 role.Role, _a1 error) *RoleService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RoleService_Create_Call) RunAndReturn(run func(context.Context, role.Role) (role.Role, error)) *RoleService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *RoleService) Get(ctx context.Context, id string) (role.Role, error) {
	ret := _m.Called(ctx, id)

	var r0 role.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (role.Role, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) role.Role); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(role.Role)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type RoleService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *RoleService_Expecter) Get(ctx interface{}, id interface{}) *RoleService_Get_Call {
	return &RoleService_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *RoleService_Get_Call) Run(run func(ctx context.Context, id string)) *RoleService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *RoleService_Get_Call) Return(_a0 role.Role, _a1 error) *RoleService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RoleService_Get_Call) RunAndReturn(run func(context.Context, string) (role.Role, error)) *RoleService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx
func (_m *RoleService) List(ctx context.Context) ([]role.Role, error) {
	ret := _m.Called(ctx)

	var r0 []role.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]role.Role, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []role.Role); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]role.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type RoleService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RoleService_Expecter) List(ctx interface{}) *RoleService_List_Call {
	return &RoleService_List_Call{Call: _e.mock.On("List", ctx)}
}

func (_c *RoleService_List_Call) Run(run func(ctx context.Context)) *RoleService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RoleService_List_Call) Return(_a0 []role.Role, _a1 error) *RoleService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RoleService_List_Call) RunAndReturn(run func(context.Context) ([]role.Role, error)) *RoleService_List_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, toUpdate
func (_m *RoleService) Update(ctx context.Context, toUpdate role.Role) (role.Role, error) {
	ret := _m.Called(ctx, toUpdate)

	var r0 role.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, role.Role) (role.Role, error)); ok {
		return rf(ctx, toUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, role.Role) role.Role); ok {
		r0 = rf(ctx, toUpdate)
	} else {
		r0 = ret.Get(0).(role.Role)
	}

	if rf, ok := ret.Get(1).(func(context.Context, role.Role) error); ok {
		r1 = rf(ctx, toUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type RoleService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - toUpdate role.Role
func (_e *RoleService_Expecter) Update(ctx interface{}, toUpdate interface{}) *RoleService_Update_Call {
	return &RoleService_Update_Call{Call: _e.mock.On("Update", ctx, toUpdate)}
}

func (_c *RoleService_Update_Call) Run(run func(ctx context.Context, toUpdate role.Role)) *RoleService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(role.Role))
	})
	return _c
}

func (_c *RoleService_Update_Call) Return(_a0 role.Role, _a1 error) *RoleService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RoleService_Update_Call) RunAndReturn(run func(context.Context, role.Role) (role.Role, error)) *RoleService_Update_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewRoleService interface {
	mock.TestingT
	Cleanup(func())
}

// NewRoleService creates a new instance of RoleService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRoleService(t mockConstructorTestingTNewRoleService) *RoleService {
	mock := &RoleService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
