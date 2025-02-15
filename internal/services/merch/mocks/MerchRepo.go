//go:build !coverage
// +build !coverage

// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	inventory "github.com/justcgh9/merch_store/internal/models/inventory"

	mock "github.com/stretchr/testify/mock"

	transaction "github.com/justcgh9/merch_store/internal/models/transaction"
)

// MerchRepo is an autogenerated mock type for the MerchRepo type
type MerchRepo struct {
	mock.Mock
}

// BuyStuff provides a mock function with given fields: username, item, cost
func (_m *MerchRepo) BuyStuff(username string, item string, cost int) error {
	ret := _m.Called(username, item, cost)

	if len(ret) == 0 {
		panic("no return value specified for BuyStuff")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int) error); ok {
		r0 = rf(username, item, cost)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBalance provides a mock function with given fields: username
func (_m *MerchRepo) GetBalance(username string) (int, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistory provides a mock function with given fields: username
func (_m *MerchRepo) GetHistory(username string) (transaction.TransactionHistory, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for GetHistory")
	}

	var r0 transaction.TransactionHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (transaction.TransactionHistory, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) transaction.TransactionHistory); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(transaction.TransactionHistory)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInventory provides a mock function with given fields: username
func (_m *MerchRepo) GetInventory(username string) ([]inventory.Item, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for GetInventory")
	}

	var r0 []inventory.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]inventory.Item, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) []inventory.Item); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]inventory.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMerchRepo creates a new instance of MerchRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMerchRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MerchRepo {
	mock := &MerchRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
