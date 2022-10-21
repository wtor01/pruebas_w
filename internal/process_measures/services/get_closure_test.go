package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Service_ProcessMeasures_GetClosure_Handle(t *testing.T) {

	type input struct {
		ctx context.Context
		dto ListClosureDto
	}

	type output struct {
		GetClosureMongo    process_measures.ProcessedMonthlyClosure
		GetClosureMongoErr error
	}

	type want struct {
		result process_measures.ProcessedMonthlyClosure
		err    error
	}

	tests := map[string]struct {
		input  input
		output output
		want   want
	}{
		"should be correct if id not empty": {
			input: input{
				ctx: context.Background(),
				dto: ListClosureDto{DistributorId: "", Id: "1234", Moment: "", Cups: ""},
			},
			output: output{
				GetClosureMongo:    process_measures.ProcessedMonthlyClosure{},
				GetClosureMongoErr: nil,
			},
			want: want{
				result: process_measures.ProcessedMonthlyClosure{},
				err:    nil,
			},
		},
		"should be correct if id empty": {
			input: input{
				ctx: context.Background(),
				dto: ListClosureDto{DistributorId: "", Id: "", Moment: "", Cups: ""},
			},
			output: output{
				GetClosureMongo:    process_measures.ProcessedMonthlyClosure{},
				GetClosureMongoErr: nil,
			},
			want: want{
				result: process_measures.ProcessedMonthlyClosure{},
				err:    nil,
			},
		},
		"should be error if getClosure error": {
			input: input{
				ctx: context.Background(),
				dto: ListClosureDto{DistributorId: "", Id: "", Moment: "", Cups: "", StartDate: time.Now(), EndDate: time.Now()},
			},
			output: output{
				GetClosureMongo:    process_measures.ProcessedMonthlyClosure{},
				GetClosureMongoErr: errors.New("error"),
			},
			want: want{
				result: process_measures.ProcessedMonthlyClosure{},
				err:    errors.New("error"),
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.ProcessMeasureClosureRepository)
			repo.Mock.On("GetClosure", testCase.input.ctx, process_measures.GetClosure{
				Id:            testCase.input.dto.Id,
				DistributorId: testCase.input.dto.DistributorId,
				CUPS:          testCase.input.dto.Cups,
				StartDate:     testCase.input.dto.StartDate,
				EndDate:       testCase.input.dto.EndDate,
				Moment:        process_measures.SelectMoment(testCase.input.dto.Moment),
			}).Return(testCase.output.GetClosureMongo, testCase.output.GetClosureMongoErr)

			srv := NewGetClosureService(repo)

			result, err := srv.Handler(testCase.input.ctx, testCase.input.dto)

			assert.Equal(t, testCase.want.result, result, testCase)
			assert.Equal(t, testCase.want.err, err, testCase)

		})
	}

}
