package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Service_Scheduler_List(t *testing.T) {

	type input struct {
		ctx context.Context
		dto ListSchedulerDto
	}

	type want struct {
		err    error
		count  int
		result []process_measures.Scheduler
	}

	type results struct {
		listSchedulerResponse      []process_measures.Scheduler
		listSchedulerResponseCount int
		listSchedulerResponseErr   error
	}

	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should return err if fail ListScheduler": {
			input: input{
				ctx: context.Background(),
				dto: ListSchedulerDto{},
			},
			want: want{
				err:    errors.New("err"),
				count:  0,
				result: []process_measures.Scheduler{},
			},
			results: results{
				listSchedulerResponse:      []process_measures.Scheduler{},
				listSchedulerResponseCount: 2,
				listSchedulerResponseErr:   errors.New("err"),
			},
		},
		"should return ok": {
			input: input{
				ctx: context.Background(),
				dto: ListSchedulerDto{},
			},
			want: want{
				err:   nil,
				count: 1,
				result: []process_measures.Scheduler{
					{
						ID:            "cbfc1a75-542b-4667-882b-f19863163866",
						DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
						Name:          "name",
						Description:   "description",
						SchedulerId:   "job-id",
						ServiceType:   "G-D",
						PointType:     "1",
						MeterType:     []string{"TLG"},
						ReadingType:   "curve",
						Format:        "* * 1 * *",
					},
				},
			},
			results: results{
				listSchedulerResponse: []process_measures.Scheduler{
					{
						ID:            "cbfc1a75-542b-4667-882b-f19863163866",
						DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
						Name:          "name",
						Description:   "description",
						SchedulerId:   "job-id",
						ServiceType:   "G-D",
						PointType:     "1",
						MeterType:     []string{"TLG"},
						ReadingType:   "curve",
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
			repo := new(mocks.SchedulerRepository)

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
