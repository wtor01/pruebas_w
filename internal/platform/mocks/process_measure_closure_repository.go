// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	process_measures "bitbucket.org/sercide/data-ingestion/internal/process_measures"
	mock "github.com/stretchr/testify/mock"
)

// ProcessMeasureClosureRepository is an autogenerated mock type for the ProcessMeasureClosureRepository type
type ProcessMeasureClosureRepository struct {
	mock.Mock
}

// CreateClosure provides a mock function with given fields: ctx, monthly
func (_m *ProcessMeasureClosureRepository) CreateClosure(ctx context.Context, monthly process_measures.ProcessedMonthlyClosure) error {
	ret := _m.Called(ctx, monthly)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, process_measures.ProcessedMonthlyClosure) error); ok {
		r0 = rf(ctx, monthly)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetClosure provides a mock function with given fields: ctx, query
func (_m *ProcessMeasureClosureRepository) GetClosure(ctx context.Context, query process_measures.GetClosure) (process_measures.ProcessedMonthlyClosure, error) {
	ret := _m.Called(ctx, query)

	var r0 process_measures.ProcessedMonthlyClosure
	if rf, ok := ret.Get(0).(func(context.Context, process_measures.GetClosure) process_measures.ProcessedMonthlyClosure); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(process_measures.ProcessedMonthlyClosure)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, process_measures.GetClosure) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetClosureOne provides a mock function with given fields: ctx, id
func (_m *ProcessMeasureClosureRepository) GetClosureOne(ctx context.Context, id string) (process_measures.ProcessedMonthlyClosure, error) {
	ret := _m.Called(ctx, id)

	var r0 process_measures.ProcessedMonthlyClosure
	if rf, ok := ret.Get(0).(func(context.Context, string) process_measures.ProcessedMonthlyClosure); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(process_measures.ProcessedMonthlyClosure)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResume provides a mock function with given fields: ctx, query
func (_m *ProcessMeasureClosureRepository) GetResume(ctx context.Context, query process_measures.GetResume) (process_measures.ResumesProcessMonthlyClosure, error) {
	ret := _m.Called(ctx, query)

	var r0 process_measures.ResumesProcessMonthlyClosure
	if rf, ok := ret.Get(0).(func(context.Context, process_measures.GetResume) process_measures.ResumesProcessMonthlyClosure); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(process_measures.ResumesProcessMonthlyClosure)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, process_measures.GetResume) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateClosure provides a mock function with given fields: ctx, id, monthly
func (_m *ProcessMeasureClosureRepository) UpdateClosure(ctx context.Context, id string, monthly process_measures.ProcessedMonthlyClosure) error {
	ret := _m.Called(ctx, id, monthly)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, process_measures.ProcessedMonthlyClosure) error); ok {
		r0 = rf(ctx, id, monthly)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewProcessMeasureClosureRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewProcessMeasureClosureRepository creates a new instance of ProcessMeasureClosureRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProcessMeasureClosureRepository(t mockConstructorTestingTNewProcessMeasureClosureRepository) *ProcessMeasureClosureRepository {
	mock := &ProcessMeasureClosureRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}