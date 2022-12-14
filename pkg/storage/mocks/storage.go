// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Storage) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Copy provides a mock function with given fields: ctx, originPath, dstPath
func (_m *Storage) Copy(ctx context.Context, originPath string, dstPath string) error {
	ret := _m.Called(ctx, originPath, dstPath)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, originPath, dstPath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, path
func (_m *Storage) Delete(ctx context.Context, path string) error {
	ret := _m.Called(ctx, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewReader provides a mock function with given fields: ctx, path
func (_m *Storage) NewReader(ctx context.Context, path string) (io.Reader, error) {
	ret := _m.Called(ctx, path)

	var r0 io.Reader
	if rf, ok := ret.Get(0).(func(context.Context, string) io.Reader); ok {
		r0 = rf(ctx, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Reader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadAll provides a mock function with given fields: ctx, path
func (_m *Storage) ReadAll(ctx context.Context, path string) ([]byte, error) {
	ret := _m.Called(ctx, path)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage(t mockConstructorTestingTNewStorage) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
