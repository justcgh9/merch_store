//go:build !coverage
// +build !coverage

// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Sender is an autogenerated mock type for the Sender type
type Sender struct {
	mock.Mock
}

// Send provides a mock function with given fields: from, to, amount
func (_m *Sender) Send(from string, to string, amount int) error {
	ret := _m.Called(from, to, amount)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int) error); ok {
		r0 = rf(from, to, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewSender creates a new instance of Sender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *Sender {
	mock := &Sender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
