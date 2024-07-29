// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	bytes "bytes"

	mock "github.com/stretchr/testify/mock"
	form "github.com/stripe/stripe-go/v79/form"

	stripe "github.com/stripe/stripe-go/v79"
)

// Backend is an autogenerated mock type for the Backend type
type Backend struct {
	mock.Mock
}

type Backend_Expecter struct {
	mock *mock.Mock
}

func (_m *Backend) EXPECT() *Backend_Expecter {
	return &Backend_Expecter{mock: &_m.Mock}
}

// Call provides a mock function with given fields: method, path, key, params, v
func (_m *Backend) Call(method string, path string, key string, params stripe.ParamsContainer, v stripe.LastResponseSetter) error {
	ret := _m.Called(method, path, key, params, v)

	if len(ret) == 0 {
		panic("no return value specified for Call")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, stripe.ParamsContainer, stripe.LastResponseSetter) error); ok {
		r0 = rf(method, path, key, params, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Backend_Call_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Call'
type Backend_Call_Call struct {
	*mock.Call
}

// Call is a helper method to define mock.On call
//   - method string
//   - path string
//   - key string
//   - params stripe.ParamsContainer
//   - v stripe.LastResponseSetter
func (_e *Backend_Expecter) Call(method interface{}, path interface{}, key interface{}, params interface{}, v interface{}) *Backend_Call_Call {
	return &Backend_Call_Call{Call: _e.mock.On("Call", method, path, key, params, v)}
}

func (_c *Backend_Call_Call) Run(run func(method string, path string, key string, params stripe.ParamsContainer, v stripe.LastResponseSetter)) *Backend_Call_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(stripe.ParamsContainer), args[4].(stripe.LastResponseSetter))
	})
	return _c
}

func (_c *Backend_Call_Call) Return(_a0 error) *Backend_Call_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Backend_Call_Call) RunAndReturn(run func(string, string, string, stripe.ParamsContainer, stripe.LastResponseSetter) error) *Backend_Call_Call {
	_c.Call.Return(run)
	return _c
}

// CallMultipart provides a mock function with given fields: method, path, key, boundary, body, params, v
func (_m *Backend) CallMultipart(method string, path string, key string, boundary string, body *bytes.Buffer, params *stripe.Params, v stripe.LastResponseSetter) error {
	ret := _m.Called(method, path, key, boundary, body, params, v)

	if len(ret) == 0 {
		panic("no return value specified for CallMultipart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, *bytes.Buffer, *stripe.Params, stripe.LastResponseSetter) error); ok {
		r0 = rf(method, path, key, boundary, body, params, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Backend_CallMultipart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CallMultipart'
type Backend_CallMultipart_Call struct {
	*mock.Call
}

// CallMultipart is a helper method to define mock.On call
//   - method string
//   - path string
//   - key string
//   - boundary string
//   - body *bytes.Buffer
//   - params *stripe.Params
//   - v stripe.LastResponseSetter
func (_e *Backend_Expecter) CallMultipart(method interface{}, path interface{}, key interface{}, boundary interface{}, body interface{}, params interface{}, v interface{}) *Backend_CallMultipart_Call {
	return &Backend_CallMultipart_Call{Call: _e.mock.On("CallMultipart", method, path, key, boundary, body, params, v)}
}

func (_c *Backend_CallMultipart_Call) Run(run func(method string, path string, key string, boundary string, body *bytes.Buffer, params *stripe.Params, v stripe.LastResponseSetter)) *Backend_CallMultipart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(string), args[4].(*bytes.Buffer), args[5].(*stripe.Params), args[6].(stripe.LastResponseSetter))
	})
	return _c
}

func (_c *Backend_CallMultipart_Call) Return(_a0 error) *Backend_CallMultipart_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Backend_CallMultipart_Call) RunAndReturn(run func(string, string, string, string, *bytes.Buffer, *stripe.Params, stripe.LastResponseSetter) error) *Backend_CallMultipart_Call {
	_c.Call.Return(run)
	return _c
}

// CallRaw provides a mock function with given fields: method, path, key, body, params, v
func (_m *Backend) CallRaw(method string, path string, key string, body *form.Values, params *stripe.Params, v stripe.LastResponseSetter) error {
	ret := _m.Called(method, path, key, body, params, v)

	if len(ret) == 0 {
		panic("no return value specified for CallRaw")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, *form.Values, *stripe.Params, stripe.LastResponseSetter) error); ok {
		r0 = rf(method, path, key, body, params, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Backend_CallRaw_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CallRaw'
type Backend_CallRaw_Call struct {
	*mock.Call
}

// CallRaw is a helper method to define mock.On call
//   - method string
//   - path string
//   - key string
//   - body *form.Values
//   - params *stripe.Params
//   - v stripe.LastResponseSetter
func (_e *Backend_Expecter) CallRaw(method interface{}, path interface{}, key interface{}, body interface{}, params interface{}, v interface{}) *Backend_CallRaw_Call {
	return &Backend_CallRaw_Call{Call: _e.mock.On("CallRaw", method, path, key, body, params, v)}
}

func (_c *Backend_CallRaw_Call) Run(run func(method string, path string, key string, body *form.Values, params *stripe.Params, v stripe.LastResponseSetter)) *Backend_CallRaw_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(*form.Values), args[4].(*stripe.Params), args[5].(stripe.LastResponseSetter))
	})
	return _c
}

func (_c *Backend_CallRaw_Call) Return(_a0 error) *Backend_CallRaw_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Backend_CallRaw_Call) RunAndReturn(run func(string, string, string, *form.Values, *stripe.Params, stripe.LastResponseSetter) error) *Backend_CallRaw_Call {
	_c.Call.Return(run)
	return _c
}

// CallStreaming provides a mock function with given fields: method, path, key, params, v
func (_m *Backend) CallStreaming(method string, path string, key string, params stripe.ParamsContainer, v stripe.StreamingLastResponseSetter) error {
	ret := _m.Called(method, path, key, params, v)

	if len(ret) == 0 {
		panic("no return value specified for CallStreaming")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, stripe.ParamsContainer, stripe.StreamingLastResponseSetter) error); ok {
		r0 = rf(method, path, key, params, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Backend_CallStreaming_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CallStreaming'
type Backend_CallStreaming_Call struct {
	*mock.Call
}

// CallStreaming is a helper method to define mock.On call
//   - method string
//   - path string
//   - key string
//   - params stripe.ParamsContainer
//   - v stripe.StreamingLastResponseSetter
func (_e *Backend_Expecter) CallStreaming(method interface{}, path interface{}, key interface{}, params interface{}, v interface{}) *Backend_CallStreaming_Call {
	return &Backend_CallStreaming_Call{Call: _e.mock.On("CallStreaming", method, path, key, params, v)}
}

func (_c *Backend_CallStreaming_Call) Run(run func(method string, path string, key string, params stripe.ParamsContainer, v stripe.StreamingLastResponseSetter)) *Backend_CallStreaming_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(stripe.ParamsContainer), args[4].(stripe.StreamingLastResponseSetter))
	})
	return _c
}

func (_c *Backend_CallStreaming_Call) Return(_a0 error) *Backend_CallStreaming_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Backend_CallStreaming_Call) RunAndReturn(run func(string, string, string, stripe.ParamsContainer, stripe.StreamingLastResponseSetter) error) *Backend_CallStreaming_Call {
	_c.Call.Return(run)
	return _c
}

// SetMaxNetworkRetries provides a mock function with given fields: maxNetworkRetries
func (_m *Backend) SetMaxNetworkRetries(maxNetworkRetries int64) {
	_m.Called(maxNetworkRetries)
}

// Backend_SetMaxNetworkRetries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetMaxNetworkRetries'
type Backend_SetMaxNetworkRetries_Call struct {
	*mock.Call
}

// SetMaxNetworkRetries is a helper method to define mock.On call
//   - maxNetworkRetries int64
func (_e *Backend_Expecter) SetMaxNetworkRetries(maxNetworkRetries interface{}) *Backend_SetMaxNetworkRetries_Call {
	return &Backend_SetMaxNetworkRetries_Call{Call: _e.mock.On("SetMaxNetworkRetries", maxNetworkRetries)}
}

func (_c *Backend_SetMaxNetworkRetries_Call) Run(run func(maxNetworkRetries int64)) *Backend_SetMaxNetworkRetries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *Backend_SetMaxNetworkRetries_Call) Return() *Backend_SetMaxNetworkRetries_Call {
	_c.Call.Return()
	return _c
}

func (_c *Backend_SetMaxNetworkRetries_Call) RunAndReturn(run func(int64)) *Backend_SetMaxNetworkRetries_Call {
	_c.Call.Return(run)
	return _c
}

// NewBackend creates a new instance of Backend. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBackend(t interface {
	mock.TestingT
	Cleanup(func())
}) *Backend {
	mock := &Backend{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
