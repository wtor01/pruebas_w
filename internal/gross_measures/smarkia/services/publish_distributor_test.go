package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	client_mocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	event_mocks "bitbucket.org/sercide/data-ingestion/pkg/event/mocks"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Unit_Smarkia_Services_PublishDistributorService_Handle(t *testing.T) {

	tests := map[string]struct {
		want                        error
		input                       PublishDistributorDTO
		publisherFail               error
		getDistributorResponse      []clients.Distributor
		getDistributorResponseCount int
		getDistributorResponseErr   error
		publishErr                  error
	}{
		"should return error if publisher fail": {
			want: errors.New(""),
			input: PublishDistributorDTO{
				ProcessName: "test",
			},
			publisherFail: errors.New(""),
			getDistributorResponse: []clients.Distributor{
				{SmarkiaID: "1", CDOS: "0130"},
			},
		},
		"should return err if distributorService.GetDistributors fail": {
			want: errors.New(""),
			input: PublishDistributorDTO{
				ProcessName: "test",
			},
			getDistributorResponse: []clients.Distributor{
				{SmarkiaID: "1", CDOS: "0130"},
			},
			getDistributorResponseErr: errors.New(""),
		},
		"should not publish messages if distributorService.GetDistributors return empty slice": {
			want: nil,
			input: PublishDistributorDTO{
				ProcessName: "test",
			},
			publisherFail:             nil,
			getDistributorResponse:    []clients.Distributor{},
			getDistributorResponseErr: nil,
		},
		"should call publish messages well": {
			want: nil,
			input: PublishDistributorDTO{
				ProcessName: "test",
			},
			publisherFail: nil,
			getDistributorResponse: []clients.Distributor{
				{SmarkiaID: "1", CDOS: "0130"},
			},
			getDistributorResponseCount: 1,
			getDistributorResponseErr:   nil,
		},
		"should return err if some message publish fail": {
			want: errors.New("number errors 1"),
			input: PublishDistributorDTO{
				ProcessName: "test",
			},
			publisherFail: nil,
			getDistributorResponse: []clients.Distributor{
				{SmarkiaID: "1", CDOS: "0130"},
			},
			getDistributorResponseCount: 1,
			getDistributorResponseErr:   nil,
			publishErr:                  errors.New(""),
		},
	}

	topicName := "topic"
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			p := new(event_mocks.Publisher)
			p.On("Close").Return(nil)
			for _, d := range testCase.getDistributorResponse {
				m := smarkia.MessageDistributorProcess{
					Type: smarkia.TypeDistributorProcess,
					Payload: smarkia.MessageDistributorProcessPayload{
						DistributorId:   d.ID,
						ProcessName:     testCase.input.ProcessName,
						DistributorCDOS: d.CDOS,
						SmarkiaId:       d.SmarkiaID,
					},
				}
				msg, _ := json.Marshal(m)
				attributes := make(map[string]string)
				attributes[event.EventTypeKey] = m.Type
				p.On("Publish", mock.Anything, topicName, msg, attributes).Return(testCase.publishErr)
			}
			inventoryClient := new(client_mocks.Inventory)
			inventoryClient.On("ListDistributors", ctx, clients.ListDistributorsDto{
				Limit: 1000,
				Values: map[string]string{
					"is_smarkia_active": "1",
				},
			}).Return(testCase.getDistributorResponse, testCase.getDistributorResponseCount, testCase.getDistributorResponseErr)
			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return p, testCase.publisherFail
			}
			h := NewPublishDistributorService(publisherCreator, inventoryClient, topicName)
			result := h.Handle(ctx, testCase.input)

			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}
