package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	client_mocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	event_mocks "bitbucket.org/sercide/data-ingestion/pkg/event/mocks"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Reprocessing_Publish_Distributor(t *testing.T) {
	type input struct {
		ctx context.Context
		dto measures.SchedulerEventPayload
	}
	type want struct {
		err error
	}
	type results struct {
		startDate          time.Time
		errGetDate         error
		errSetDate         error
		getDistributors    []clients.Distributor
		errGetDistributors error
		errPublish         error
		publisherFail      error
		searchScheduler    []process_measures.Scheduler
		errSearchScheduler error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should be correct": {
			input: input{
				ctx: context.Background(),
				dto: measures.SchedulerEventPayload{},
			},
			results: results{
				startDate:          time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC),
				errGetDate:         nil,
				errSetDate:         nil,
				errPublish:         nil,
				errGetDistributors: nil,
				publisherFail:      nil,
				getDistributors:    []clients.Distributor{clients.Distributor{}},
			},
			want: want{
				err: nil,
			},
		},
		"should be correct 1 distributor": {
			input: input{
				ctx: context.Background(),
				dto: measures.SchedulerEventPayload{DistributorId: "TEST"},
			},
			results: results{
				startDate:          time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC),
				errGetDate:         nil,
				errSetDate:         nil,
				errPublish:         nil,
				errGetDistributors: nil,
				publisherFail:      nil,
				getDistributors:    []clients.Distributor{clients.Distributor{}},
			},
			want: want{
				err: nil,
			},
		},
		"should be error GetDate": {
			input: input{
				ctx: context.Background(),
				dto: measures.SchedulerEventPayload{},
			},
			results: results{
				startDate:          time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC),
				errGetDate:         errors.New("error"),
				errSetDate:         nil,
				errPublish:         nil,
				errGetDistributors: nil,
				publisherFail:      nil,
				getDistributors:    []clients.Distributor{clients.Distributor{}},
			},
			want: want{
				err: errors.New("no redis cache start date"),
			},
		},
		"should be error 1st SetDate": {
			input: input{
				ctx: context.Background(),
				dto: measures.SchedulerEventPayload{},
			},
			results: results{
				startDate:          time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC),
				errGetDate:         errors.New("error"),
				errSetDate:         errors.New("error"),
				errPublish:         nil,
				errGetDistributors: nil,
				publisherFail:      nil,
				getDistributors:    []clients.Distributor{clients.Distributor{}},
			},
			want: want{
				err: errors.New("error"),
			},
		},
		"should be error 2nd SetDate": {
			input: input{
				ctx: context.Background(),
				dto: measures.SchedulerEventPayload{},
			},
			results: results{
				startDate:          time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC),
				errGetDate:         nil,
				errSetDate:         errors.New("error"),
				errPublish:         nil,
				errGetDistributors: nil,
				publisherFail:      nil,
				getDistributors:    []clients.Distributor{clients.Distributor{}},
			},
			want: want{
				err: errors.New("error"),
			},
		},
		"should be error get Distributor": {
			input: input{
				ctx: context.Background(),
				dto: measures.SchedulerEventPayload{},
			},
			results: results{
				startDate:          time.Date(2022, time.Month(9), 10, 10, 0, 0, 0, time.UTC),
				errGetDate:         nil,
				errSetDate:         nil,
				errPublish:         nil,
				errGetDistributors: errors.New("error"),
				publisherFail:      nil,
				getDistributors:    []clients.Distributor{clients.Distributor{}},
			},
			want: want{
				err: errors.New("error"),
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			generatorDate := func() time.Time {
				return time.Date(2022, time.Month(9), 10, 0, 0, 0, 0, time.UTC)
			}
			mockPublisher := new(event_mocks.Publisher)
			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return mockPublisher, testCase.results.errPublish
			}
			schedulerClientMock := new(mocks.SchedulerRepository)
			inventoryRepository := new(client_mocks.Inventory)
			redis := new(mocks.ReprocessingDateRepository)
			redis.Mock.On("GetDate", testCase.input.ctx, getReprocessingDistributorRedis).Return(testCase.results.errGetDate, testCase.results.startDate)
			redis.Mock.On("SetDate", testCase.input.ctx, getReprocessingDistributorRedis, generatorDate()).Return(testCase.results.errSetDate)
			inventoryRepository.Mock.On("GetAllDistributors", testCase.input.ctx).Return(testCase.results.getDistributors, testCase.results.errGetDistributors)
			schedulerClientMock.Mock.On("SearchScheduler", testCase.input.ctx, process_measures.SearchScheduler{}).Return(testCase.results.searchScheduler, testCase.results.errSearchScheduler)
			mockPublisher.Mock.On("Close").Return(nil)
			for range testCase.results.getDistributors {
				m := process_measures.NewReprocessingProcessByDistributorEvent(process_measures.ReSchedulerEventPayload{StartDate: testCase.results.startDate, EndDate: generatorDate(), DistributorId: testCase.input.dto.DistributorId})
				msg, _ := json.Marshal(m)
				attributes := make(map[string]string)
				attributes[event.EventTypeKey] = m.Type
				mockPublisher.Mock.On("Publish", mock.Anything, "", msg, attributes).Return(testCase.results.errPublish)

			}
			srv := NewReprocessingPublishDistributorService(publisherCreator, NewSearchSchedulerService(schedulerClientMock), inventoryRepository, "", redis, generatorDate)
			err := srv.Handle(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.want.err, err, testCase)

		})

	}

}
