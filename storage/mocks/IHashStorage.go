// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	model "github.com/zhayt/cert-tz/model"
)

// IHashStorage is an autogenerated mock type for the IHashStorage type
type IHashStorage struct {
	mock.Mock
}

// CreateHash provides a mock function with given fields: hash
func (_m *IHashStorage) CreateHash(hash model.CertHash) (uint64, error) {
	ret := _m.Called(hash)

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(model.CertHash) (uint64, error)); ok {
		return rf(hash)
	}
	if rf, ok := ret.Get(0).(func(model.CertHash) uint64); ok {
		r0 = rf(hash)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(model.CertHash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHash provides a mock function with given fields: hashID
func (_m *IHashStorage) GetHash(hashID uint64) (model.CertHash, error) {
	ret := _m.Called(hashID)

	var r0 model.CertHash
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (model.CertHash, error)); ok {
		return rf(hashID)
	}
	if rf, ok := ret.Get(0).(func(uint64) model.CertHash); ok {
		r0 = rf(hashID)
	} else {
		r0 = ret.Get(0).(model.CertHash)
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(hashID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateHash provides a mock function with given fields: hash
func (_m *IHashStorage) UpdateHash(hash model.CertHash) error {
	ret := _m.Called(hash)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.CertHash) error); ok {
		r0 = rf(hash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIHashStorage creates a new instance of IHashStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIHashStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *IHashStorage {
	mock := &IHashStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
