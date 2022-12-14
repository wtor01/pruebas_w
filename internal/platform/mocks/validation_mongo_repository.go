// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"
	testing "testing"

	mock "github.com/stretchr/testify/mock"

	validations "bitbucket.org/sercide/data-ingestion/internal/validations"
)

// ValidationMongoRepository is an autogenerated mock type for the ValidationMongoRepository type
type ValidationMongoRepository struct {
	mock.Mock
}

// GetLoadCurveByQuery provides a mock function with given fields: ctx, q
func (_m *ValidationMongoRepository) GetLoadCurveByQuery(ctx context.Context, q validations.QueryCurveCupsMeasureOnDate) ([]validations.ProcessedLoadCurve, error) {
	ret := _m.Called(ctx, q)

	var r0 []validations.ProcessedLoadCurve
	if rf, ok := ret.Get(0).(func(context.Context, validations.QueryCurveCupsMeasureOnDate) []validations.ProcessedLoadCurve); ok {
		r0 = rf(ctx, q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]validations.ProcessedLoadCurve)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, validations.QueryCurveCupsMeasureOnDate) error); ok {
		r1 = rf(ctx, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMonthlyClosureByCup provides a mock function with given fields: ctx, q
func (_m *ValidationMongoRepository) GetMonthlyClosureByCup(ctx context.Context, q validations.QueryClosedCupsMeasureOnDate) (validations.ProcessedMonthlyClosure, error) {
	ret := _m.Called(ctx, q)

	var r0 validations.ProcessedMonthlyClosure
	if rf, ok := ret.Get(0).(func(context.Context, validations.QueryClosedCupsMeasureOnDate) validations.ProcessedMonthlyClosure); ok {
		r0 = rf(ctx, q)
	} else {
		r0 = ret.Get(0).(validations.ProcessedMonthlyClosure)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, validations.QueryClosedCupsMeasureOnDate) error); ok {
		r1 = rf(ctx, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewValidationMongoRepository creates a new instance of ValidationMongoRepository. It also registers a cleanup function to assert the mocks expectations.
func NewValidationMongoRepository(t testing.TB) *ValidationMongoRepository {
	mock := &ValidationMongoRepository{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
