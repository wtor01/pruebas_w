package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services/fixtures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Service_ProcessMeasures_ListDashboardCups_Handle(t *testing.T) {

	type input struct {
		Dto ListDashboardCupsDTO
	}

	type output struct {
		GetMetersConfig      measures.GetMetersAndCountResult
		GetMetersConfigErr   error
		ListDashboardCups    map[string]*measures.DashboardCupsReading
		ListDashboardCupsErr error
	}

	type want struct {
		result measures.DashboardListCups
		err    error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be err, GetMetersErr": {
			input: input{
				Dto: NewListDashboardCupsDTO("DistributorX", 10, 0, "TLG", time.Date(2022, 4, 30, 22, 0, 0, 0, time.UTC), time.Date(2022, 5, 30, 22, 0, 0, 0, time.UTC)),
			},
			output: output{
				GetMetersConfig:    measures.GetMetersAndCountResult{},
				GetMetersConfigErr: errors.New("err GetMeters"),
			},
			want: want{
				err: errors.New("err GetMeters"),
			},
		},
		"Should be err, GetCupsMeasures err": {
			input: input{
				Dto: NewListDashboardCupsDTO("DistributorX", 10, 0, "TLG", time.Date(2022, 4, 30, 22, 0, 0, 0, time.UTC), time.Date(2022, 5, 30, 22, 0, 0, 0, time.UTC)),
			},
			output: output{
				GetMetersConfig:      fixtures.RESULT_GET_METERS_AND_COUNT,
				ListDashboardCupsErr: errors.New("err GetCupsMeasures"),
			},
			want: want{
				err: errors.New("err GetCupsMeasures"),
			},
		},
		"Should be return correct values": {
			input: input{
				Dto: NewListDashboardCupsDTO("DistributorX", 10, 0, "TLG", time.Date(2022, 4, 30, 22, 0, 0, 0, time.UTC), time.Date(2022, 5, 30, 22, 0, 0, 0, time.UTC)),
			},
			output: output{
				GetMetersConfig:   fixtures.RESULT_GET_METERS_AND_COUNT,
				ListDashboardCups: fixtures.RESULT_LIST_CUPS,
			},
			want: want{
				result: fixtures.RESULT_DASHBOARD_CUPS_LIST,
				err:    nil,
			},
		},
	}

	loc, _ := time.LoadLocation("Europe/Madrid")
	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			inventoryRepository := new(mocks.InventoryRepository)
			inventoryRepository.On("GetMetersAndCountByDistributorId", mock.Anything, measures.GetMetersAndCountByDistributorIdQuery{
				DistributorId: test.input.Dto.DistributorId,
				Type:          measures.MeterType(test.input.Dto.Type),
				StartDate:     test.input.Dto.StartDate,
				EndDate:       test.input.Dto.EndDate,
				Limit:         test.input.Dto.Limit,
				Offset:        test.input.Dto.Offset,
			}).Return(test.output.GetMetersConfig, test.output.GetMetersConfigErr)

			dashboardInventory := new(mocks.ProcessMeasureDashboardRepository)
			dashboardInventory.On("GetCupsMeasures", mock.Anything, process_measures.ListCupsQuery{
				Cups:      []string{"ESXXXX01", "ESXXXX02", "ESXXXX03", "ESXXXX04", "ESXXXX05"},
				StartDate: test.input.Dto.StartDate,
				EndDate:   test.input.Dto.EndDate,
			}).Return(test.output.ListDashboardCups, test.output.ListDashboardCupsErr)

			srv := NewListDashboardCupsService(inventoryRepository, dashboardInventory, loc)
			result, err := srv.Handle(context.Background(), test.input.Dto)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.result, result)
		})
	}
}
