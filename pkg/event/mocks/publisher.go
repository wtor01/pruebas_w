// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Publisher is an autogenerated mock type for the Publisher type
type Publisher struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Publisher) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publish provides a mock function with given fields: ctx, topic, message, attributes
func (_m *Publisher) Publish(ctx context.Context, topic string, message []byte, attributes map[string]string) error {
	ret := _m.Called(ctx, topic, message, attributes)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte, map[string]string) error); ok {
		r0 = rf(ctx, topic, message, attributes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPublisher creates a new instance of Publisher. It also registers a cleanup function to assert the mocks expectations.
func NewPublisher(t testing.TB) *Publisher {
	mock := &Publisher{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
