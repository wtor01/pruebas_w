package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clients_mocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
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

func Test_Unit_Services_Aggregations_ProcessAggregationInitService(t *testing.T) {
	type expected struct {
		getAllDistributorsResult []clients.Distributor
		getAllDistributorsErr    error
		publishAllEventsErr      error
		shouldReturnEarly        bool
	}
	today := time.Now()

	tests := map[string]struct {
		input    aggregations.ConfigScheduler
		expected expected
		want     error
	}{
		"should not process if date not in range": {
			input: aggregations.ConfigScheduler{
				Date: today,
				Config: aggregations.Config{
					EndDate:   today.AddDate(0, -1, 0),
					StartDate: today.AddDate(0, -2, 0),
				},
			},
			expected: expected{
				getAllDistributorsResult: []clients.Distributor{},
				getAllDistributorsErr:    nil,
				shouldReturnEarly:        true,
			},
			want: nil,
		},
		"should not process if end date is empty and not in range": {
			input: aggregations.ConfigScheduler{
				Date: today,
				Config: aggregations.Config{
					StartDate: today.AddDate(0, 1, 0),
				},
			},
			expected: expected{
				getAllDistributorsResult: []clients.Distributor{},
				getAllDistributorsErr:    nil,
				shouldReturnEarly:        true,
			},
			want: nil,
		},
		"should return err if GetAllDistributors fail": {
			input: aggregations.ConfigScheduler{
				Date: today,
				Config: aggregations.Config{
					StartDate: today.AddDate(0, -1, 0),
				},
			},
			expected: expected{
				getAllDistributorsResult: []clients.Distributor{},
				getAllDistributorsErr:    errors.New(""),
			},
			want: errors.New(""),
		},
		"should publish messages": {
			input: aggregations.ConfigScheduler{
				Date: today,
				Config: aggregations.Config{
					EndDate:   today.AddDate(0, 1, 0),
					StartDate: today.AddDate(0, -1, 0),
				},
			},
			expected: expected{
				getAllDistributorsResult: []clients.Distributor{
					{
						ID:   "604930e9-84a7-4c7e-891e-436477f22e64",
						CDOS: "0132",
					},
					{
						ID:   "b84c1fbe-5832-4796-a0df-dbfc16e59850",
						CDOS: "0165",
					},
				},
				getAllDistributorsErr: nil,
			},
			want: nil,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			topicName := "topic"
			p := new(event_mocks.Publisher)
			p.On("Close").Return(nil)
			for _, d := range testCase.expected.getAllDistributorsResult {
				dto := testCase.input
				dto.DistributorId = d.ID
				dto.DistributorCDOS = d.CDOS

				m := aggregations.NewSchedulerDistributorEvent(dto)
				msg, _ := json.Marshal(m)
				attributes := make(map[string]string)
				attributes[event.EventTypeKey] = m.Type
				p.On("Publish", mock.Anything, topicName, msg, attributes).Return(testCase.expected.publishAllEventsErr)
			}

			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return p, nil
			}

			inventoryClient := new(clients_mocks.Inventory)
			inventoryClient.On("GetAllDistributors", mock.Anything).Return(testCase.expected.getAllDistributorsResult, testCase.expected.getAllDistributorsErr)

			svc := NewProcessAggregationInitService(publisherCreator, inventoryClient, topicName)
			result := svc.Handler(context.Background(), testCase.input)

			if testCase.expected.shouldReturnEarly {
				p.AssertNumberOfCalls(t, "Publish", 0)
				inventoryClient.AssertNumberOfCalls(t, "GetAllDistributors", 0)
			} else {
				inventoryClient.AssertNumberOfCalls(t, "GetAllDistributors", 1)
				p.AssertNumberOfCalls(t, "Publish", len(testCase.expected.getAllDistributorsResult))
			}

			assert.Equal(t, testCase.want, result)
		})
	}
}
