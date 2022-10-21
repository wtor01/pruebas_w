package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	event_mocks "bitbucket.org/sercide/data-ingestion/pkg/event/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Execute_Services(t *testing.T) {
	type input struct {
		ctx context.Context
		dto DtoServiceExecute
	}

	type want struct {
		err error
	}
	type results struct {
		meterConfigResult []measures.MeterConfig
		errMeterConfig    error
		errPublisher      error
		errPublish        error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should be correct": {
			input: input{
				ctx: context.Background(),
				dto: DtoServiceExecute{ReadingType: measures.Curve},
			},
			want: want{err: nil},
			results: results{
				meterConfigResult: []measures.MeterConfig{measures.MeterConfig{}},
				errMeterConfig:    nil,
				errPublisher:      nil,
				errPublish:        nil,
			},
		},
		"should be error Meter": {
			input: input{
				ctx: context.Background(),
				dto: DtoServiceExecute{
					ReadingType: measures.Curve,
					StartDate:   time.Date(2022, time.Month(9), 22, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2022, time.Month(9), 22, 0, 0, 0, 0, time.UTC),
				},
			},
			want: want{err: errors.New("error")},
			results: results{
				meterConfigResult: []measures.MeterConfig{measures.MeterConfig{}},
				errMeterConfig:    errors.New("error"),
				errPublisher:      nil,
				errPublish:        nil,
			},
		},
		"should be error Publisher": {
			input: input{
				ctx: context.Background(),
				dto: DtoServiceExecute{
					ReadingType: measures.DailyClosure,
					StartDate:   time.Date(2022, time.Month(9), 22, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2022, time.Month(9), 22, 0, 0, 0, 0, time.UTC),
				},
			},
			want: want{err: errors.New("error")},
			results: results{
				meterConfigResult: []measures.MeterConfig{measures.MeterConfig{}},
				errMeterConfig:    nil,
				errPublisher:      errors.New("error"),
				errPublish:        nil,
			},
		},
		"should be error publish": {
			input: input{
				ctx: context.Background(),
				dto: DtoServiceExecute{
					ReadingType: measures.BillingClosure,
					StartDate:   time.Date(2022, time.Month(9), 22, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2022, time.Month(9), 22, 0, 0, 0, 0, time.UTC),
				},
			},
			want: want{err: errors.New("number errors 1")},
			results: results{
				meterConfigResult: []measures.MeterConfig{measures.MeterConfig{}},
				errMeterConfig:    nil,
				errPublisher:      nil,
				errPublish:        errors.New("error"),
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {

			p := new(event_mocks.Publisher)
			measureRepository := new(mocks.InventoryRepository)
			measureRepository.Mock.On("ListMeterConfigByCups", testCase.input.ctx, measures.ListMeterConfigByCups{
				CUPS:          testCase.input.dto.Cups,
				StartDate:     testCase.input.dto.StartDate,
				EndDate:       testCase.input.dto.EndDate,
				DistributorId: testCase.input.dto.DistributorId,
			}).Return(testCase.results.meterConfigResult, testCase.results.errMeterConfig)
			p.Mock.On("Publish", mock.Anything, "", mock.Anything, mock.Anything).Return(testCase.results.errPublish)
			p.Mock.On("Close").Return(nil)
			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return p, testCase.results.errPublisher
			}
			srv := NewExecuteServices(measureRepository, publisherCreator, "")
			result := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.want.err, result)
		})
	}
}
