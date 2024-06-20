// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	subscription "github.com/raystack/frontier/billing/subscription"
	mock "github.com/stretchr/testify/mock"
)

// SubscriptionService is an autogenerated mock type for the SubscriptionService type
type SubscriptionService struct {
	mock.Mock
}

type SubscriptionService_Expecter struct {
	mock *mock.Mock
}

func (_m *SubscriptionService) EXPECT() *SubscriptionService_Expecter {
	return &SubscriptionService_Expecter{mock: &_m.Mock}
}

// Cancel provides a mock function with given fields: ctx, id, immediate
func (_m *SubscriptionService) Cancel(ctx context.Context, id string, immediate bool) (subscription.Subscription, error) {
	ret := _m.Called(ctx, id, immediate)

	if len(ret) == 0 {
		panic("no return value specified for Cancel")
	}

	var r0 subscription.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) (subscription.Subscription, error)); ok {
		return rf(ctx, id, immediate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) subscription.Subscription); ok {
		r0 = rf(ctx, id, immediate)
	} else {
		r0 = ret.Get(0).(subscription.Subscription)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, bool) error); ok {
		r1 = rf(ctx, id, immediate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionService_Cancel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Cancel'
type SubscriptionService_Cancel_Call struct {
	*mock.Call
}

// Cancel is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - immediate bool
func (_e *SubscriptionService_Expecter) Cancel(ctx interface{}, id interface{}, immediate interface{}) *SubscriptionService_Cancel_Call {
	return &SubscriptionService_Cancel_Call{Call: _e.mock.On("Cancel", ctx, id, immediate)}
}

func (_c *SubscriptionService_Cancel_Call) Run(run func(ctx context.Context, id string, immediate bool)) *SubscriptionService_Cancel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *SubscriptionService_Cancel_Call) Return(_a0 subscription.Subscription, _a1 error) *SubscriptionService_Cancel_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubscriptionService_Cancel_Call) RunAndReturn(run func(context.Context, string, bool) (subscription.Subscription, error)) *SubscriptionService_Cancel_Call {
	_c.Call.Return(run)
	return _c
}

// ChangePlan provides a mock function with given fields: ctx, id, change
func (_m *SubscriptionService) ChangePlan(ctx context.Context, id string, change subscription.ChangeRequest) (subscription.Phase, error) {
	ret := _m.Called(ctx, id, change)

	if len(ret) == 0 {
		panic("no return value specified for ChangePlan")
	}

	var r0 subscription.Phase
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, subscription.ChangeRequest) (subscription.Phase, error)); ok {
		return rf(ctx, id, change)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, subscription.ChangeRequest) subscription.Phase); ok {
		r0 = rf(ctx, id, change)
	} else {
		r0 = ret.Get(0).(subscription.Phase)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, subscription.ChangeRequest) error); ok {
		r1 = rf(ctx, id, change)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionService_ChangePlan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChangePlan'
type SubscriptionService_ChangePlan_Call struct {
	*mock.Call
}

// ChangePlan is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - change subscription.ChangeRequest
func (_e *SubscriptionService_Expecter) ChangePlan(ctx interface{}, id interface{}, change interface{}) *SubscriptionService_ChangePlan_Call {
	return &SubscriptionService_ChangePlan_Call{Call: _e.mock.On("ChangePlan", ctx, id, change)}
}

func (_c *SubscriptionService_ChangePlan_Call) Run(run func(ctx context.Context, id string, change subscription.ChangeRequest)) *SubscriptionService_ChangePlan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(subscription.ChangeRequest))
	})
	return _c
}

func (_c *SubscriptionService_ChangePlan_Call) Return(_a0 subscription.Phase, _a1 error) *SubscriptionService_ChangePlan_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubscriptionService_ChangePlan_Call) RunAndReturn(run func(context.Context, string, subscription.ChangeRequest) (subscription.Phase, error)) *SubscriptionService_ChangePlan_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *SubscriptionService) GetByID(ctx context.Context, id string) (subscription.Subscription, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 subscription.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (subscription.Subscription, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) subscription.Subscription); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(subscription.Subscription)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionService_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type SubscriptionService_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *SubscriptionService_Expecter) GetByID(ctx interface{}, id interface{}) *SubscriptionService_GetByID_Call {
	return &SubscriptionService_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *SubscriptionService_GetByID_Call) Run(run func(ctx context.Context, id string)) *SubscriptionService_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SubscriptionService_GetByID_Call) Return(_a0 subscription.Subscription, _a1 error) *SubscriptionService_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubscriptionService_GetByID_Call) RunAndReturn(run func(context.Context, string) (subscription.Subscription, error)) *SubscriptionService_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// HasUserSubscribedBefore provides a mock function with given fields: ctx, customerID, planID
func (_m *SubscriptionService) HasUserSubscribedBefore(ctx context.Context, customerID string, planID string) (bool, error) {
	ret := _m.Called(ctx, customerID, planID)

	if len(ret) == 0 {
		panic("no return value specified for HasUserSubscribedBefore")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (bool, error)); ok {
		return rf(ctx, customerID, planID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, customerID, planID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, customerID, planID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionService_HasUserSubscribedBefore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasUserSubscribedBefore'
type SubscriptionService_HasUserSubscribedBefore_Call struct {
	*mock.Call
}

// HasUserSubscribedBefore is a helper method to define mock.On call
//   - ctx context.Context
//   - customerID string
//   - planID string
func (_e *SubscriptionService_Expecter) HasUserSubscribedBefore(ctx interface{}, customerID interface{}, planID interface{}) *SubscriptionService_HasUserSubscribedBefore_Call {
	return &SubscriptionService_HasUserSubscribedBefore_Call{Call: _e.mock.On("HasUserSubscribedBefore", ctx, customerID, planID)}
}

func (_c *SubscriptionService_HasUserSubscribedBefore_Call) Run(run func(ctx context.Context, customerID string, planID string)) *SubscriptionService_HasUserSubscribedBefore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *SubscriptionService_HasUserSubscribedBefore_Call) Return(_a0 bool, _a1 error) *SubscriptionService_HasUserSubscribedBefore_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubscriptionService_HasUserSubscribedBefore_Call) RunAndReturn(run func(context.Context, string, string) (bool, error)) *SubscriptionService_HasUserSubscribedBefore_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, filter
func (_m *SubscriptionService) List(ctx context.Context, filter subscription.Filter) ([]subscription.Subscription, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []subscription.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, subscription.Filter) ([]subscription.Subscription, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, subscription.Filter) []subscription.Subscription); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]subscription.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, subscription.Filter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type SubscriptionService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - filter subscription.Filter
func (_e *SubscriptionService_Expecter) List(ctx interface{}, filter interface{}) *SubscriptionService_List_Call {
	return &SubscriptionService_List_Call{Call: _e.mock.On("List", ctx, filter)}
}

func (_c *SubscriptionService_List_Call) Run(run func(ctx context.Context, filter subscription.Filter)) *SubscriptionService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(subscription.Filter))
	})
	return _c
}

func (_c *SubscriptionService_List_Call) Return(_a0 []subscription.Subscription, _a1 error) *SubscriptionService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubscriptionService_List_Call) RunAndReturn(run func(context.Context, subscription.Filter) ([]subscription.Subscription, error)) *SubscriptionService_List_Call {
	_c.Call.Return(run)
	return _c
}

// NewSubscriptionService creates a new instance of SubscriptionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSubscriptionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *SubscriptionService {
	mock := &SubscriptionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
