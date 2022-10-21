// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	billing_measures "bitbucket.org/sercide/data-ingestion/internal/billing_measures"

	mock "github.com/stretchr/testify/mock"
)

// BillingSelfConsumptionRepository is an autogenerated mock type for the BillingSelfConsumptionRepository type
type BillingSelfConsumptionRepository struct {
	mock.Mock
}

// GetBySelfConsumptionBetweenDates provides a mock function with given fields: ctx, query
func (_m *BillingSelfConsumptionRepository) GetBySelfConsumptionBetweenDates(ctx context.Context, query billing_measures.QueryGetBillingSelfConsumption) (billing_measures.BillingSelfConsumption, error) {
	ret := _m.Called(ctx, query)

	var r0 billing_measures.BillingSelfConsumption
	if rf, ok := ret.Get(0).(func(context.Context, billing_measures.QueryGetBillingSelfConsumption) billing_measures.BillingSelfConsumption); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(billing_measures.BillingSelfConsumption)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, billing_measures.QueryGetBillingSelfConsumption) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSelfConsumptionByCau provides a mock function with given fields: ctx, query
func (_m *BillingSelfConsumptionRepository) GetSelfConsumptionByCau(ctx context.Context, query billing_measures.QueryGetBillingSelfConsumptionByCau) ([]billing_measures.BillingSelfConsumption, error) {
	ret := _m.Called(ctx, query)

	var r0 []billing_measures.BillingSelfConsumption
	if rf, ok := ret.Get(0).(func(context.Context, billing_measures.QueryGetBillingSelfConsumptionByCau) []billing_measures.BillingSelfConsumption); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]billing_measures.BillingSelfConsumption)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, billing_measures.QueryGetBillingSelfConsumptionByCau) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, b
func (_m *BillingSelfConsumptionRepository) Save(ctx context.Context, b billing_measures.BillingSelfConsumption) error {
	ret := _m.Called(ctx, b)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, billing_measures.BillingSelfConsumption) error); ok {
		r0 = rf(ctx, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewBillingSelfConsumptionRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewBillingSelfConsumptionRepository creates a new instance of BillingSelfConsumptionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBillingSelfConsumptionRepository(t mockConstructorTestingTNewBillingSelfConsumptionRepository) *BillingSelfConsumptionRepository {
	mock := &BillingSelfConsumptionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
