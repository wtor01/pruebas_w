package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Services_Closure_Put(t *testing.T) {

	type input struct {
		ctx     context.Context
		monthly UpdateClosureDto
	}

	type want struct {
		err error
	}

	type results struct {
		resultGetClosureOne    process_measures.ProcessedMonthlyClosure
		resultGetClosureOneErr error
		resultUpdateClosureErr error
	}

	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should be error in getClosureOne": {
			input: input{
				ctx:     context.Background(),
				monthly: UpdateClosureDto{},
			},
			results: results{
				resultGetClosureOne:    process_measures.ProcessedMonthlyClosure{},
				resultGetClosureOneErr: errors.New("error"),
				resultUpdateClosureErr: nil,
			},
			want: want{
				err: errors.New("error"),
			},
		},
		"should be error in updateClosure": {
			input: input{
				ctx:     context.Background(),
				monthly: UpdateClosureDto{},
			},
			results: results{
				resultGetClosureOne:    process_measures.ProcessedMonthlyClosure{},
				resultGetClosureOneErr: nil,
				resultUpdateClosureErr: errors.New("error"),
			},
			want: want{
				err: errors.New("error"),
			},
		},
		"should correct": {
			input: input{
				ctx:     context.Background(),
				monthly: UpdateClosureDto{},
			},
			results: results{
				resultGetClosureOne:    process_measures.ProcessedMonthlyClosure{},
				resultGetClosureOneErr: nil,
				resultUpdateClosureErr: nil,
			},
			want: want{
				err: nil,
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repoClosureOne := new(mocks.ProcessMeasureClosureRepository)

			repoClosureOne.Mock.On("GetClosureOne", testCase.input.ctx,
				testCase.input.monthly.Monthly.Id).Return(testCase.results.resultGetClosureOne, testCase.results.resultGetClosureOneErr)

			repoClosureOne.Mock.On("UpdateClosure", testCase.input.ctx,
				testCase.input.monthly.Monthly.Id, testCase.results.resultGetClosureOne).Return(testCase.results.resultUpdateClosureErr)

			srv := NewUpdateClosureService(repoClosureOne)

			err := srv.Handler(testCase.input.ctx,
				testCase.input.monthly)

			assert.Equal(t, testCase.want.err, err, testCase)

		})
	}
}
