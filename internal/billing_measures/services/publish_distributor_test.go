package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clients_mocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	event_mocks "bitbucket.org/sercide/data-ingestion/pkg/event/mocks"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Unit_ProccessMeasure_Services_PublishDistributorService_Handle(t *testing.T) {

	tests := map[string]struct {
		want            error
		input           measures.SchedulerEventPayload
		publisherFail   error
		getDistributors []clients.Distributor
		schedulers      []billing_measures.Scheduler
		getSchedulerErr error
		publishErr      error
	}{
		"should return error if publisher fail": {
			want:  errors.New(""),
			input: measures.SchedulerEventPayload{}, publisherFail: errors.New("")},

		"should not publish messages if GetAllSchedulerConfigs return empty slice": {
			want:            nil,
			input:           measures.SchedulerEventPayload{},
			publisherFail:   nil,
			getDistributors: []clients.Distributor{},
			schedulers:      []billing_measures.Scheduler{},
			getSchedulerErr: nil,
		},
		"should call publish messages well": {
			want:            nil,
			input:           measures.SchedulerEventPayload{},
			publisherFail:   nil,
			getDistributors: []clients.Distributor{},
			schedulers:      []billing_measures.Scheduler{},
			getSchedulerErr: nil,
		},
	}

	topicName := "topic"
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			p := new(event_mocks.Publisher)
			p.On("Close").Return(nil)

			for _, d := range testCase.schedulers {
				m := billing_measures.NewProcessByDistributorEvent(measures.SchedulerEventPayload{
					DistributorId: d.DistributorId,
				})
				msg, _ := json.Marshal(m)
				attributes := make(map[string]string)
				attributes[event.EventTypeKey] = billing_measures.TypeDistributorProcess
				p.On("Publish", mock.Anything, topicName, msg, attributes).Return(testCase.publishErr)
			}

			g := new(mocks.BillingSchedulerRepository)
			g.On("SearchScheduler", mock.Anything).Return(testCase.schedulers, testCase.getSchedulerErr)

			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return p, testCase.publisherFail
			}

			invent := new(clients_mocks.Inventory)
			invent.On("GetAllDistributors", mock.Anything).Return(testCase.getDistributors, nil)

			schedulerService := NewSearchSchedulerService(g)
			h := NewPublishDistributorService(publisherCreator, schedulerService, invent, topicName)
			result := h.Handle(context.Background(), testCase.input)

			p.AssertNumberOfCalls(t, "Publish", len(testCase.schedulers))
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}
