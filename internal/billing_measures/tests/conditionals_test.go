package tests

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHaveBillingHistory(t *testing.T) {

	type input struct {
		b          *billing_measures.BillingMeasure
		ContextTlg *billing_measures.GraphContext
	}

	type output struct {
		billingHistory billing_measures.BillingMeasure
		lastHistoryErr error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   bool
	}{
		"should be true, have history": {
			input: input{
				b:          &billing_measures.BillingMeasure{},
				ContextTlg: &billing_measures.GraphContext{},
			},
			output: output{
				billingHistory: billing_measures.BillingMeasure{
					Id: "1",
				},
				lastHistoryErr: nil,
			},
			want: true,
		},
		"should be false, no history, error": {
			input: input{
				b:          &billing_measures.BillingMeasure{},
				ContextTlg: &billing_measures.GraphContext{},
			},
			output: output{
				billingHistory: billing_measures.BillingMeasure{},
				lastHistoryErr: errors.New("no: results"),
			},
			want: false,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			billingMeasureRepository := mocks.NewBillingMeasureRepository(t)
			billingMeasureRepository.On("LastHistory", ctx, billing_measures.QueryLastHistory{
				CUPS:       test.input.b.CUPS,
				InitDate:   test.input.b.InitDate,
				EndDate:    test.input.b.EndDate,
				Periods:    []measures.PeriodKey{measures.P1},
				Magnitudes: []measures.Magnitude{measures.AI},
			}).Return(test.output.billingHistory, test.output.lastHistoryErr)

			haveBillingHistory := billing_measures.NewHasBillingHistory(
				test.input.b,
				billingMeasureRepository,
				measures.P1,
				measures.AI,
				test.input.ContextTlg,
			)

			result := haveBillingHistory.Eval(ctx)

			if result {
				assert.Equal(t, test.input.ContextTlg.LastHistory, test.output.billingHistory)
			} else {
				assert.Empty(t, test.input.ContextTlg.LastHistory)
			}
			assert.Equal(t, result, test.want)
		})
	}
}
