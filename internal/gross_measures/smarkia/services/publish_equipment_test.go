package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
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

func Test_Unit_Smarkia_Services_PublishEquipmentService_Handle(t *testing.T) {

	tests := map[string]struct {
		want                  error
		input                 PublishEquipmentDto
		publisherFail         error
		publishErr            error
		getEquipmentResult    []smarkia.Equipment
		getEquipmentResultErr error
	}{
		"should return error if publisher fail": {
			want:          errors.New(""),
			input:         PublishEquipmentDto{},
			publisherFail: errors.New(""),
		},
		"should return error if api.GetEquipments fail": {
			want: errors.New(""),
			input: PublishEquipmentDto{
				DistributorId: "distributorX",
				SmarkiaId:     "1",
				ProcessName:   "curva",
				CtId:          "ctid",
			},
			getEquipmentResultErr: errors.New(""),
			publisherFail:         errors.New(""),
		},
		"should not publish messages if api.GetEquipments return empty slice": {
			want: nil,
			input: PublishEquipmentDto{
				DistributorId: "distributorX",
				ProcessName:   "curva",
				SmarkiaId:     "1",
				CtId:          "ctid",
			},
			getEquipmentResult:    []smarkia.Equipment{},
			getEquipmentResultErr: nil,
			publisherFail:         nil,
		},
		"should call publish messages well": {
			want: nil,
			input: PublishEquipmentDto{
				DistributorId: "distributorX",
				ProcessName:   "curva",
				SmarkiaId:     "1",
				CtId:          "ctid",
			},
			publisherFail: nil,
			getEquipmentResult: []smarkia.Equipment{
				{ID: "1"},
			},
			getEquipmentResultErr: nil,
		},
		"should return err if some message publish fail": {
			want: errors.New("number errors 1"),
			input: PublishEquipmentDto{
				DistributorId: "distributorX",
				ProcessName:   "curva",
				SmarkiaId:     "1",
				CtId:          "ctid",
			},
			publishErr: errors.New(""),
			getEquipmentResult: []smarkia.Equipment{
				{ID: "1"},
			},
			getEquipmentResultErr: nil,
		},
	}

	topicName := "topic"
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			p := new(event_mocks.Publisher)
			p.On("Close").Return(nil)

			for _, equipment := range testCase.getEquipmentResult {
				m := smarkia.NewEquipmentProcessEvent(smarkia.MessageEquipmentDto{
					ProcessName:     testCase.input.ProcessName,
					DistributorId:   testCase.input.DistributorId,
					SmarkiaId:       testCase.input.SmarkiaId,
					CtId:            testCase.input.CtId,
					DistributorCDOS: testCase.input.DistributorCDOS,
					Date:            testCase.input.Date,
				}, equipment.ID, "")
				msg, _ := json.Marshal(m)
				attributes := make(map[string]string)
				attributes[event.EventTypeKey] = m.Type
				p.On("Publish", mock.Anything, topicName, msg, attributes).Return(testCase.publishErr)
			}

			g := new(mocks.GetEquipmentser)

			g.On("GetEquipments", ctx, smarkia.GetEquipmentsQuery{
				Id: testCase.input.CtId,
			}).Return(testCase.getEquipmentResult, testCase.getEquipmentResultErr)

			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return p, testCase.publisherFail
			}

			h := NewPublishEquipmentService(publisherCreator, g, topicName)
			result := h.Handle(ctx, testCase.input)

			p.AssertNumberOfCalls(t, "Publish", len(testCase.getEquipmentResult))
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}
