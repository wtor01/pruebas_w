package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
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

func Test_Publish_Reprocessing_Service_Point_Service(t *testing.T) {
	type input struct {
		ctx context.Context
		dto process_measures.ReSchedulerEventPayload
	}
	type want struct {
		err error
	}
	type results struct {
		ListGrossMeasuresFromGenerationDate    gross_measures.MeasureCurveMeterSerialNumber
		ErrListGrossMeasuresFromGenerationDate error
		errPublish                             error
		publisherFail                          error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should be correct": {
			input: input{
				ctx: context.Background(),
				dto: process_measures.ReSchedulerEventPayload{ReadingType: measures.DailyClosure, StartDate: time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC), EndDate: time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC), Limit: 1},
			},
			results: results{
				ListGrossMeasuresFromGenerationDate:    gross_measures.MeasureCurveMeterSerialNumber{Year: 2022, Month: 9, Day: 10},
				ErrListGrossMeasuresFromGenerationDate: nil,
				errPublish:                             nil,
			},
			want: want{
				err: nil,
			},
		},
		"should be error": {
			input: input{
				ctx: context.Background(),
				dto: process_measures.ReSchedulerEventPayload{ReadingType: measures.DailyClosure, StartDate: time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC), EndDate: time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC), Limit: 1},
			},
			results: results{
				ListGrossMeasuresFromGenerationDate:    gross_measures.MeasureCurveMeterSerialNumber{},
				ErrListGrossMeasuresFromGenerationDate: errors.New("error"),
				errPublish:                             nil,
			},
			want: want{
				err: errors.New("error"),
			},
		},

		"should be correct but error Publish": {
			input: input{
				ctx: context.Background(),
				dto: process_measures.ReSchedulerEventPayload{ReadingType: measures.DailyClosure, StartDate: time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC), EndDate: time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC), Limit: 1},
			},
			results: results{
				ListGrossMeasuresFromGenerationDate:    gross_measures.MeasureCurveMeterSerialNumber{},
				ErrListGrossMeasuresFromGenerationDate: nil,
				errPublish:                             errors.New("error"),
			},
			want: want{
				err: errors.New("number errors 1"),
			},
		},
		"should be error Publish": {
			input: input{
				ctx: context.Background(),
				dto: process_measures.ReSchedulerEventPayload{ReadingType: measures.BillingClosure, StartDate: time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC), EndDate: time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC), Limit: 1},
			},
			results: results{
				ListGrossMeasuresFromGenerationDate:    gross_measures.MeasureCurveMeterSerialNumber{},
				ErrListGrossMeasuresFromGenerationDate: nil,
				errPublish:                             nil,
				publisherFail:                          errors.New("error"),
			},
			want: want{
				err: errors.New("error"),
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {
			mockPublisher := new(event_mocks.Publisher)
			mockPublisher.On("Close").Return(nil)

			mockGrossRepository := new(mocks.GrossMeasureRepository)
			res := make([]gross_measures.MeasureCurveMeterSerialNumber, 0)
			res = append(res, testCase.results.ListGrossMeasuresFromGenerationDate)
			mockGrossRepository.Mock.On("ListGrossMeasuresFromGenerationDate", mock.Anything,
				gross_measures.QueryListForProcessCurveGenerationDate{
					ReadingType:   testCase.input.dto.ReadingType,
					DistributorId: testCase.input.dto.DistributorId,
					StartDate:     testCase.input.dto.StartDate,
					EndDate:       testCase.input.dto.EndDate,
					Limit:         testCase.input.dto.Limit,
					Offset:        testCase.input.dto.Offset,
				}).Return(res, testCase.results.ErrListGrossMeasuresFromGenerationDate)

			mockPublisher.On("Publish", mock.Anything, "", mock.Anything, mock.Anything).Return(testCase.results.errPublish)

			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return mockPublisher, testCase.results.publisherFail
			}
			srv := NewReprocessingServicePointService(publisherCreator, "", mockGrossRepository, 100)
			err := srv.Handle(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.want.err, err, testCase)

		})

	}
}
