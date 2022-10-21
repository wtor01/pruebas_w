package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Get_Resume(t *testing.T) {
	type input struct {
		ctx context.Context
		dto ListResumeDto
	}
	type want struct {
		err    error
		result process_measures.ResumesProcessMonthlyClosure
	}
	type results struct {
		getResume    process_measures.ResumesProcessMonthlyClosure
		getResumeErr error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should return ok": {
			input: input{
				ctx: context.Background(),
				dto: ListResumeDto{},
			},
			want: want{
				err:    nil,
				result: process_measures.ResumesProcessMonthlyClosure{},
			},
			results: results{
				getResume:    process_measures.ResumesProcessMonthlyClosure{},
				getResumeErr: nil,
			},
		},
		"should return err": {
			input: input{
				ctx: context.Background(),
				dto: ListResumeDto{},
			},
			want: want{
				err:    errors.New("error"),
				result: process_measures.ResumesProcessMonthlyClosure{},
			},
			results: results{
				getResume:    process_measures.ResumesProcessMonthlyClosure{},
				getResumeErr: errors.New("error"),
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.ProcessMeasureClosureRepository)
			repo.Mock.On("GetResume", mock.Anything, process_measures.GetResume{
				Cups:          testCase.input.dto.Cups,
				DistributorId: testCase.input.dto.DistributorId,
				StartDate:     testCase.input.dto.StartDate,
				EndDate:       testCase.input.dto.EndDate.AddDate(0, 0, 1),
			}).Return(testCase.results.getResume, testCase.results.getResumeErr)
			srv := NewGetResumeService(repo, time.UTC)
			result, err := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.want.result, result, testCase)
			assert.Equal(t, testCase.want.err, err, testCase)

		})
	}
}
