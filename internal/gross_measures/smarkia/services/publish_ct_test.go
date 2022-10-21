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

func Test_Unit_Smarkia_Services_PublishCtService_Handle(t *testing.T) {

	tests := map[string]struct {
		want            error
		input           PublishCtDTO
		publisherFail   error
		publishErr      error
		getCTsResult    []smarkia.Ct
		getCTsResultErr error
	}{
		"should return error if publisher fail": {
			want:          errors.New(""),
			input:         PublishCtDTO{},
			publisherFail: errors.New(""),
			getCTsResult: []smarkia.Ct{
				{ID: "1"},
			},
		},
		"should return error if api.GetCTs fail": {
			want: errors.New(""),
			input: PublishCtDTO{
				DistributorId: "distributorX",
				ProcessName:   "curva",
			},
			getCTsResultErr: errors.New(""),
			publisherFail:   errors.New(""),
		},
		"should not publish messages if api.GetCTs return empty slice": {
			want: nil,
			input: PublishCtDTO{
				DistributorId: "distributorX",
				SmarkiaId:     "1",
				ProcessName:   "curva",
			},
			getCTsResult:    []smarkia.Ct{},
			getCTsResultErr: nil,
			publisherFail:   nil,
		},
		"should call publish messages well": {
			want: nil,
			input: PublishCtDTO{
				DistributorId: "distributorX",
				SmarkiaId:     "1",
				ProcessName:   "curva",
			},
			publisherFail: nil,
			getCTsResult: []smarkia.Ct{
				{ID: "1"},
			},
			getCTsResultErr: nil,
		},
		"should return err if some message publish fail": {
			want: errors.New("number errors 1"),
			input: PublishCtDTO{
				DistributorId: "distributorX",
				SmarkiaId:     "1",
				ProcessName:   "curva",
			},
			publishErr: errors.New(""),
			getCTsResult: []smarkia.Ct{
				{ID: "1"},
			},
			getCTsResultErr: nil,
		},
	}

	topicName := "topic"
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			p := new(event_mocks.Publisher)
			p.On("Close").Return(nil)

			for _, ct := range testCase.getCTsResult {
				m := smarkia.NewCtProcessEvent(
					testCase.input.DistributorId,
					testCase.input.SmarkiaId,
					testCase.input.ProcessName,
					testCase.input.DistributorCDOS,
					ct.ID, testCase.input.Date,
				)
				msg, _ := json.Marshal(m)
				attributes := make(map[string]string)
				attributes[event.EventTypeKey] = m.Type
				p.On("Publish", mock.Anything, topicName, msg, attributes).Return(testCase.publishErr)
			}

			g := new(mocks.GetCTer)

			g.On("GetCTs", ctx, testCase.input.SmarkiaId).Return(testCase.getCTsResult, testCase.getCTsResultErr)

			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return p, testCase.publisherFail
			}

			h := NewPublishCtService(publisherCreator, g, topicName)
			result := h.Handle(ctx, testCase.input)

			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}
