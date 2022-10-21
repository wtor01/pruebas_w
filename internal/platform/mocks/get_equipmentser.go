// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	smarkia "bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	mock "github.com/stretchr/testify/mock"
)

// GetEquipmentser is an autogenerated mock type for the GetEquipmentser type
type GetEquipmentser struct {
	mock.Mock
}

// GetEquipments provides a mock function with given fields: ctx, query
func (_m *GetEquipmentser) GetEquipments(ctx context.Context, query smarkia.GetEquipmentsQuery) ([]smarkia.Equipment, error) {
	ret := _m.Called(ctx, query)

	var r0 []smarkia.Equipment
	if rf, ok := ret.Get(0).(func(context.Context, smarkia.GetEquipmentsQuery) []smarkia.Equipment); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]smarkia.Equipment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, smarkia.GetEquipmentsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGetEquipmentser interface {
	mock.TestingT
	Cleanup(func())
}

// NewGetEquipmentser creates a new instance of GetEquipmentser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGetEquipmentser(t mockConstructorTestingTNewGetEquipmentser) *GetEquipmentser {
	mock := &GetEquipmentser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}