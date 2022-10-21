// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gross_measures "bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// GetMagnitudeser is an autogenerated mock type for the GetMagnitudeser type
type GetMagnitudeser struct {
	mock.Mock
}

// GetMagnitudes provides a mock function with given fields: ctx, equipmentID, date
func (_m *GetMagnitudeser) GetMagnitudes(ctx context.Context, equipmentID string, date time.Time) ([]gross_measures.MeasureCurveWrite, error) {
	ret := _m.Called(ctx, equipmentID, date)

	var r0 []gross_measures.MeasureCurveWrite
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) []gross_measures.MeasureCurveWrite); ok {
		r0 = rf(ctx, equipmentID, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]gross_measures.MeasureCurveWrite)
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

type mockConstructorTestingTNewGetMagnitudeser interface {
	mock.TestingT
	Cleanup(func())
}

// NewGetMagnitudeser creates a new instance of GetMagnitudeser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGetMagnitudeser(t mockConstructorTestingTNewGetMagnitudeser) *GetMagnitudeser {
	mock := &GetMagnitudeser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
