// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Authenticator is an autogenerated mock type for the Authenticator type
type Authenticator struct {
	mock.Mock
}

// Authorize provides a mock function with given fields: username, password
func (_m *Authenticator) Authorize(username string, password string) (string, error) {
	ret := _m.Called(username, password)

	if len(ret) == 0 {
		panic("no return value specified for Authorize")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(username, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuthenticator creates a new instance of Authenticator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthenticator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Authenticator {
	mock := &Authenticator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
