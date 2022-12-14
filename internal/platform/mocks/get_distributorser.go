// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	smarkia "bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	mock "github.com/stretchr/testify/mock"
)

// GetDistributorser is an autogenerated mock type for the GetDistributorser type
type GetDistributorser struct {
	mock.Mock
}

// GetDistributors provides a mock function with given fields: ctx
func (_m *GetDistributorser) GetDistributors(ctx context.Context) ([]smarkia.Distributor, error) {
	ret := _m.Called(ctx)

	var r0 []smarkia.Distributor
	if rf, ok := ret.Get(0).(func(context.Context) []smarkia.Distributor); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]smarkia.Distributor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGetDistributorser interface {
	mock.TestingT
	Cleanup(func())
}

// NewGetDistributorser creates a new instance of GetDistributorser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGetDistributorser(t mockConstructorTestingTNewGetDistributorser) *GetDistributorser {
	mock := &GetDistributorser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
