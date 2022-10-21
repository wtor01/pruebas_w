package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Service_ProcessMeasures_GetDashboardMeasures_Handle(t *testing.T) {
	type input struct {
		Dto GetDashboardMeasuresDTO
	}

	type output struct {
		GetDashboard    []measures.DashboardMeasureI
		GetDashboardErr error
		GroupMeters     map[measures.MeterType]map[measures.RegisterType]measures.MeasureCount
		GroupMetersErr  error
	}

	type want struct {
		result measures.DashboardResult
		err    error
	}

	testCase := map[string]struct {
		input  input
		output output
		want   want
	}{
		"should be err, get dashboard error": {
			input: input{
				Dto: NewGetDashboardMeasuresDTO("DistributorX", time.Date(2022, 4, 30, 22, 0, 0, 0, time.UTC), time.Date(2022, 5, 30, 22, 0, 0, 0, time.UTC)),
			},
			output: output{
				GetDashboard:    []measures.DashboardMeasureI{},
				GetDashboardErr: errors.New("err get dashboard"),
			},
			want: want{
				result: measures.DashboardResult{},
				err:    errors.New("err get dashboard"),
			},
		},
		"should be err, groupMeters error": {
			input: input{
				Dto: NewGetDashboardMeasuresDTO("DistributorX", time.Date(2022, 4, 30, 22, 0, 0, 0, time.UTC), time.Date(2022, 5, 30, 22, 0, 0, 0, time.UTC)),
			},
			output: output{
				GetDashboard:   []measures.DashboardMeasureI{},
				GroupMeters:    map[measures.MeterType]map[measures.RegisterType]measures.MeasureCount{},
				GroupMetersErr: errors.New("err group meters"),
			},
			want: want{
				result: measures.DashboardResult{},
				err:    errors.New("err group meters"),
			},
		},

		//TODO: SHOULDBE CORRECT
	}

	loc, _ := time.LoadLocation("Europe/Madrid")
	for name := range testCase {
		test := testCase[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			dashboardRepository := new(mocks.ProcessMeasureDashboardRepository)
			dashboardRepository.On("GetDashboard", mock.Anything, process_measures.GetDashboardQuery{
				DistributorId: test.input.Dto.DistributorID,
				StartDate:     test.input.Dto.StartDate,
				EndDate:       test.input.Dto.EndDate,
			}).Return(test.output.GetDashboard, test.output.GetDashboardErr)

			inventoryRepository := new(mocks.InventoryRepository)
			inventoryRepository.On("GroupMetersByType", mock.Anything, measures.GroupMetersByTypeQuery{
				DistributorId: test.input.Dto.DistributorID,
				StartDate:     test.input.Dto.StartDate,
				EndDate:       test.input.Dto.EndDate,
			}).Return(test.output.GroupMeters, test.output.GroupMetersErr)

			srv := NewGetDashboardMeasures(dashboardRepository, inventoryRepository, loc)

			result, err := srv.Handle(ctx, test.input.Dto)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.result, result)
		})
	}
}
