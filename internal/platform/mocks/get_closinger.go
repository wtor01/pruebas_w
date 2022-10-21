// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gross_measures "bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// GetClosinger is an autogenerated mock type for the GetClosinger type
type GetClosinger struct {
	mock.Mock
}

// GetClosinger provides a mock function with given fields: ctx, equipmentID, date
func (_m *GetClosinger) GetClosinger(ctx context.Context, equipmentID string, date time.Time) ([]gross_measures.MeasureCloseWrite, error) {
	ret := _m.Called(ctx, equipmentID, date)

	var r0 []gross_measures.MeasureCloseWrite
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) []gross_measures.MeasureCloseWrite); ok {
		r0 = rf(ctx, equipmentID, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]gross_measures.MeasureCloseWrite)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time) error); ok {
		r1 = rf(ctx, equipmentID, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGetClosinger interface {
	mock.TestingT
	Cleanup(func())
}

// NewGetClosinger creates a new instance of GetClosinger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGetClosinger(t mockConstructorTestingTNewGetClosinger) *GetClosinger {
	mock := &GetClosinger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
