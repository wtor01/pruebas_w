package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	event_mocks "bitbucket.org/sercide/data-ingestion/pkg/event/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Publish_Reprocessing_Meter_Service(t *testing.T) {
	type input struct {
		ctx context.Context
		dto process_measures.ReSchedulerMeterPayload
	}
	type want struct {
		err error
	}
	type results struct {
		MeterConfig    measures.MeterConfig
		errMeterConfig error
		errPublisher   error
		errPublish     error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should be correct": {
			input: input{
				ctx: context.Background(),
				dto: process_measures.ReSchedulerMeterPayload{ReadingType: measures.BillingClosure, Date: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC)},
			},
			want: want{
				err: nil,
			},
			results: results{
				MeterConfig:    measures.MeterConfig{StartDate: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC)},
				errMeterConfig: nil,
				errPublisher:   nil,
				errPublish:     nil,
			},
		},
		"should be error Publisher": {
			input: input{
				ctx: context.Background(),
				dto: process_measures.ReSchedulerMeterPayload{ReadingType: measures.BillingClosure, Date: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC)},
			},
			want: want{
				err: errors.New("error"),
			},
			results: results{
				MeterConfig:    measures.MeterConfig{StartDate: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC)},
				errMeterConfig: nil,
				errPublisher:   errors.New("error"),
				errPublish:     nil,
			},
		},
		"should be error Publish": {
			input: input{
				ctx: context.Background(),
				dto: process_measures.ReSchedulerMeterPayload{ReadingType: measures.BillingClosure, Date: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC)},
			},
			want: want{
				err: errors.New("number errors 1"),
			},
			results: results{
				MeterConfig:    measures.MeterConfig{StartDate: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2022, time.Month(9), 1, 0, 0, 0, 0, time.UTC)},
				errMeterConfig: nil,
				errPublisher:   nil,
				errPublish:     errors.New("error"),
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			mockPublisher := new(event_mocks.Publisher)
			mockPublisher.On("Close").Return(nil)
			mockPublisher.On("Publish", mock.Anything, "", mock.Anything, mock.Anything).Return(testCase.results.errPublish)
			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return mockPublisher, testCase.results.errPublisher
			}
			inventory := new(mocks.InventoryRepository)
			inventory.Mock.On("GetMeterConfigByMeter", testCase.input.ctx, measures.GetMeterConfigByMeterQuery{
				MeterSerialNumber: testCase.input.dto.MeterSerialNumber,
				Date:              testCase.input.dto.Date,
			}).Return(testCase.results.MeterConfig, testCase.results.errMeterConfig)
			srv := NewReprocessingMeterService(publisherCreator, inventory, "")
			err := srv.Handle(testCase.input.ctx, testCase.input.dto)

			assert.Equal(t, testCase.want.err, err, testCase)
		})
	}
}
