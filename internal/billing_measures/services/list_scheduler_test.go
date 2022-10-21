package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Service_Billing_Scheduler_List(t *testing.T) {
	type input struct {
		ctx context.Context
		dto ListSchedulerDto
	}
	type want struct {
		err    error
		count  int
		result []billing_measures.Scheduler
	}
	type results struct {
		listSchedulerResponse      []billing_measures.Scheduler
		listSchedulerResponseCount int
		listSchedulerResponseErr   error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"Should return err if fail ListScheduler": {
			input: input{
				ctx: context.Background(),
				dto: ListSchedulerDto{},
			},
			want: want{
				err:    errors.New("ListScheduler Have Failed"),
				count:  0,
				result: []billing_measures.Scheduler{},
			},
			results: results{
				listSchedulerResponse:      []billing_measures.Scheduler{},
				listSchedulerResponseCount: 2,
				listSchedulerResponseErr:   errors.New("ListScheduler Have Failed"),
			},
		},
		"Should return ok": {
			input: input{
				ctx: context.Background(),
				dto: ListSchedulerDto{},
			},
			want: want{
				err:   nil,
				count: 1,
				result: []billing_measures.Scheduler{
					{
						ID:            "cbfc1a75-542b-4667-882b-f19863163866",
						DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
						Name:          "name",
						SchedulerId:   "job-id",
						ServiceType:   "G-D",
						PointType:     "1",
						MeterType:     []string{"TLG"},
						ProcessType:   "D-C TLG",
						Format:        "* * 1 * *",
					},
				},
			},
			results: results{
				listSchedulerResponse: []billing_measures.Scheduler{
					{
						ID:            "cbfc1a75-542b-4667-882b-f19863163866",
						DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
						Name:          "name",
						SchedulerId:   "job-id",
						ServiceType:   "G-D",
						PointType:     "1",
						MeterType:     []string{"TLG"},
						ProcessType:   "D-C TLG",
						Format:        "* * 1 * *",
					},
				},
				listSchedulerResponseCount: 1,
				listSchedulerResponseErr:   nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.BillingSchedulerRepository)

			repo.Mock.On("ListScheduler", testCase.input.ctx, db.Pagination{
				Limit:  testCase.input.dto.Limit,
				Offset: testCase.input.dto.Offset,
			}).Return(testCase.results.listSchedulerResponse, testCase.results.listSchedulerResponseCount, testCase.results.listSchedulerResponseErr)
			srv := NewListSchedulerService(repo)
			result, count, err := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.want.result, result, testCase)
			assert.Equal(t, testCase.want.count, count, testCase)
			assert.Equal(t, testCase.want.err, err, testCase)
		})
	}
}
