// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"
	auth "go-api/src/models/auth"

	mock "github.com/stretchr/testify/mock"
)

// AuthService is an autogenerated mock type for the AuthService type
type AuthService struct {
	mock.Mock
}

type AuthService_Expecter struct {
	mock *mock.Mock
}

func (_m *AuthService) EXPECT() *AuthService_Expecter {
	return &AuthService_Expecter{mock: &_m.Mock}
}

// CreateSession provides a mock function with given fields: ctx, request
func (_m *AuthService) CreateSession(ctx context.Context, request auth.CreateSessionRequest) (*auth.SessionInfo, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for CreateSession")
	}

	var r0 *auth.SessionInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.CreateSessionRequest) (*auth.SessionInfo, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, auth.CreateSessionRequest) *auth.SessionInfo); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.SessionInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, auth.CreateSessionRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthService_CreateSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateSession'
type AuthService_CreateSession_Call struct {
	*mock.Call
}

// CreateSession is a helper method to define mock.On call
//   - ctx context.Context
//   - request auth.CreateSessionRequest
func (_e *AuthService_Expecter) CreateSession(ctx interface{}, request interface{}) *AuthService_CreateSession_Call {
	return &AuthService_CreateSession_Call{Call: _e.mock.On("CreateSession", ctx, request)}
}

func (_c *AuthService_CreateSession_Call) Run(run func(ctx context.Context, request auth.CreateSessionRequest)) *AuthService_CreateSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(auth.CreateSessionRequest))
	})
	return _c
}

func (_c *AuthService_CreateSession_Call) Return(_a0 *auth.SessionInfo, _a1 error) *AuthService_CreateSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthService_CreateSession_Call) RunAndReturn(run func(context.Context, auth.CreateSessionRequest) (*auth.SessionInfo, error)) *AuthService_CreateSession_Call {
	_c.Call.Return(run)
	return _c
}

// FinishSession provides a mock function with given fields: ctx, request
func (_m *AuthService) FinishSession(ctx context.Context, request auth.FinishSessionRequest) (*auth.FinishSessionResponse, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for FinishSession")
	}

	var r0 *auth.FinishSessionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.FinishSessionRequest) (*auth.FinishSessionResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, auth.FinishSessionRequest) *auth.FinishSessionResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.FinishSessionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, auth.FinishSessionRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthService_FinishSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FinishSession'
type AuthService_FinishSession_Call struct {
	*mock.Call
}

// FinishSession is a helper method to define mock.On call
//   - ctx context.Context
//   - request auth.FinishSessionRequest
func (_e *AuthService_Expecter) FinishSession(ctx interface{}, request interface{}) *AuthService_FinishSession_Call {
	return &AuthService_FinishSession_Call{Call: _e.mock.On("FinishSession", ctx, request)}
}

func (_c *AuthService_FinishSession_Call) Run(run func(ctx context.Context, request auth.FinishSessionRequest)) *AuthService_FinishSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(auth.FinishSessionRequest))
	})
	return _c
}

func (_c *AuthService_FinishSession_Call) Return(_a0 *auth.FinishSessionResponse, _a1 error) *AuthService_FinishSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthService_FinishSession_Call) RunAndReturn(run func(context.Context, auth.FinishSessionRequest) (*auth.FinishSessionResponse, error)) *AuthService_FinishSession_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserInfo provides a mock function with given fields: ctx, request
func (_m *AuthService) GetUserInfo(ctx context.Context, request auth.VerifySessionRequest) (*auth.UserInfo, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for GetUserInfo")
	}

	var r0 *auth.UserInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.VerifySessionRequest) (*auth.UserInfo, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, auth.VerifySessionRequest) *auth.UserInfo); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.UserInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, auth.VerifySessionRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthService_GetUserInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserInfo'
type AuthService_GetUserInfo_Call struct {
	*mock.Call
}

// GetUserInfo is a helper method to define mock.On call
//   - ctx context.Context
//   - request auth.VerifySessionRequest
func (_e *AuthService_Expecter) GetUserInfo(ctx interface{}, request interface{}) *AuthService_GetUserInfo_Call {
	return &AuthService_GetUserInfo_Call{Call: _e.mock.On("GetUserInfo", ctx, request)}
}

func (_c *AuthService_GetUserInfo_Call) Run(run func(ctx context.Context, request auth.VerifySessionRequest)) *AuthService_GetUserInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(auth.VerifySessionRequest))
	})
	return _c
}

func (_c *AuthService_GetUserInfo_Call) Return(_a0 *auth.UserInfo, _a1 error) *AuthService_GetUserInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthService_GetUserInfo_Call) RunAndReturn(run func(context.Context, auth.VerifySessionRequest) (*auth.UserInfo, error)) *AuthService_GetUserInfo_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSession provides a mock function with given fields: ctx, request
func (_m *AuthService) UpdateSession(ctx context.Context, request auth.UpdateSessionRequest) (*auth.SessionInfo, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSession")
	}

	var r0 *auth.SessionInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.UpdateSessionRequest) (*auth.SessionInfo, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, auth.UpdateSessionRequest) *auth.SessionInfo); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.SessionInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, auth.UpdateSessionRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthService_UpdateSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSession'
type AuthService_UpdateSession_Call struct {
	*mock.Call
}

// UpdateSession is a helper method to define mock.On call
//   - ctx context.Context
//   - request auth.UpdateSessionRequest
func (_e *AuthService_Expecter) UpdateSession(ctx interface{}, request interface{}) *AuthService_UpdateSession_Call {
	return &AuthService_UpdateSession_Call{Call: _e.mock.On("UpdateSession", ctx, request)}
}

func (_c *AuthService_UpdateSession_Call) Run(run func(ctx context.Context, request auth.UpdateSessionRequest)) *AuthService_UpdateSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(auth.UpdateSessionRequest))
	})
	return _c
}

func (_c *AuthService_UpdateSession_Call) Return(_a0 *auth.SessionInfo, _a1 error) *AuthService_UpdateSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthService_UpdateSession_Call) RunAndReturn(run func(context.Context, auth.UpdateSessionRequest) (*auth.SessionInfo, error)) *AuthService_UpdateSession_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuthService creates a new instance of AuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthService {
	mock := &AuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
