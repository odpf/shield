// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	authenticate "github.com/raystack/shield/core/authenticate"

	mock "github.com/stretchr/testify/mock"

	organization "github.com/raystack/shield/core/organization"

	user "github.com/raystack/shield/core/user"
)

// RegistrationService is an autogenerated mock type for the RegistrationService type
type RegistrationService struct {
	mock.Mock
}

type RegistrationService_Expecter struct {
	mock *mock.Mock
}

func (_m *RegistrationService) EXPECT() *RegistrationService_Expecter {
	return &RegistrationService_Expecter{mock: &_m.Mock}
}

// Finish provides a mock function with given fields: ctx, request
func (_m *RegistrationService) Finish(ctx context.Context, request authenticate.RegistrationFinishRequest) (*authenticate.RegistrationFinishResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *authenticate.RegistrationFinishResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, authenticate.RegistrationFinishRequest) (*authenticate.RegistrationFinishResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, authenticate.RegistrationFinishRequest) *authenticate.RegistrationFinishResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authenticate.RegistrationFinishResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, authenticate.RegistrationFinishRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegistrationService_Finish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Finish'
type RegistrationService_Finish_Call struct {
	*mock.Call
}

// Finish is a helper method to define mock.On call
//   - ctx context.Context
//   - request authenticate.RegistrationFinishRequest
func (_e *RegistrationService_Expecter) Finish(ctx interface{}, request interface{}) *RegistrationService_Finish_Call {
	return &RegistrationService_Finish_Call{Call: _e.mock.On("Finish", ctx, request)}
}

func (_c *RegistrationService_Finish_Call) Run(run func(ctx context.Context, request authenticate.RegistrationFinishRequest)) *RegistrationService_Finish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(authenticate.RegistrationFinishRequest))
	})
	return _c
}

func (_c *RegistrationService_Finish_Call) Return(_a0 *authenticate.RegistrationFinishResponse, _a1 error) *RegistrationService_Finish_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RegistrationService_Finish_Call) RunAndReturn(run func(context.Context, authenticate.RegistrationFinishRequest) (*authenticate.RegistrationFinishResponse, error)) *RegistrationService_Finish_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: ctx, request
func (_m *RegistrationService) Start(ctx context.Context, request authenticate.RegistrationStartRequest) (*authenticate.RegistrationStartResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *authenticate.RegistrationStartResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, authenticate.RegistrationStartRequest) (*authenticate.RegistrationStartResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, authenticate.RegistrationStartRequest) *authenticate.RegistrationStartResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authenticate.RegistrationStartResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, authenticate.RegistrationStartRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegistrationService_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type RegistrationService_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - ctx context.Context
//   - request authenticate.RegistrationStartRequest
func (_e *RegistrationService_Expecter) Start(ctx interface{}, request interface{}) *RegistrationService_Start_Call {
	return &RegistrationService_Start_Call{Call: _e.mock.On("Start", ctx, request)}
}

func (_c *RegistrationService_Start_Call) Run(run func(ctx context.Context, request authenticate.RegistrationStartRequest)) *RegistrationService_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(authenticate.RegistrationStartRequest))
	})
	return _c
}

func (_c *RegistrationService_Start_Call) Return(_a0 *authenticate.RegistrationStartResponse, _a1 error) *RegistrationService_Start_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RegistrationService_Start_Call) RunAndReturn(run func(context.Context, authenticate.RegistrationStartRequest) (*authenticate.RegistrationStartResponse, error)) *RegistrationService_Start_Call {
	_c.Call.Return(run)
	return _c
}

// SupportedStrategies provides a mock function with given fields:
func (_m *RegistrationService) SupportedStrategies() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// RegistrationService_SupportedStrategies_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SupportedStrategies'
type RegistrationService_SupportedStrategies_Call struct {
	*mock.Call
}

// SupportedStrategies is a helper method to define mock.On call
func (_e *RegistrationService_Expecter) SupportedStrategies() *RegistrationService_SupportedStrategies_Call {
	return &RegistrationService_SupportedStrategies_Call{Call: _e.mock.On("SupportedStrategies")}
}

func (_c *RegistrationService_SupportedStrategies_Call) Run(run func()) *RegistrationService_SupportedStrategies_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *RegistrationService_SupportedStrategies_Call) Return(_a0 []string) *RegistrationService_SupportedStrategies_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RegistrationService_SupportedStrategies_Call) RunAndReturn(run func() []string) *RegistrationService_SupportedStrategies_Call {
	_c.Call.Return(run)
	return _c
}

// Token provides a mock function with given fields: _a0, orgs
func (_m *RegistrationService) Token(_a0 user.User, orgs []organization.Organization) ([]byte, error) {
	ret := _m.Called(_a0, orgs)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(user.User, []organization.Organization) ([]byte, error)); ok {
		return rf(_a0, orgs)
	}
	if rf, ok := ret.Get(0).(func(user.User, []organization.Organization) []byte); ok {
		r0 = rf(_a0, orgs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(user.User, []organization.Organization) error); ok {
		r1 = rf(_a0, orgs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegistrationService_Token_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Token'
type RegistrationService_Token_Call struct {
	*mock.Call
}

// Token is a helper method to define mock.On call
//   - _a0 user.User
//   - orgs []organization.Organization
func (_e *RegistrationService_Expecter) Token(_a0 interface{}, orgs interface{}) *RegistrationService_Token_Call {
	return &RegistrationService_Token_Call{Call: _e.mock.On("Token", _a0, orgs)}
}

func (_c *RegistrationService_Token_Call) Run(run func(_a0 user.User, orgs []organization.Organization)) *RegistrationService_Token_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(user.User), args[1].([]organization.Organization))
	})
	return _c
}

func (_c *RegistrationService_Token_Call) Return(_a0 []byte, _a1 error) *RegistrationService_Token_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RegistrationService_Token_Call) RunAndReturn(run func(user.User, []organization.Organization) ([]byte, error)) *RegistrationService_Token_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewRegistrationService interface {
	mock.TestingT
	Cleanup(func())
}

// NewRegistrationService creates a new instance of RegistrationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRegistrationService(t mockConstructorTestingTNewRegistrationService) *RegistrationService {
	mock := &RegistrationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
