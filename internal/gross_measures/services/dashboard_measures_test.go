package services

import (
	mocks2 "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/services/fixtures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Measure_Services_DashboardService_Handle(t *testing.T) {

	loc, _ := time.LoadLocation("Europe/Madrid")

	tests := map[string]struct {
		want                      measures.DashboardResult
		wantErr                   error
		input                     DashboardServiceDTO
		repoGetDashboardResult    []measures.DashboardMeasureI
		repoGetDashboardResultErr error
		clientResult              map[string]int
		clientResultErr           error
	}{
		"should return err if GetDashboard fail": {
			want:                      measures.DashboardResult{},
			wantErr:                   errors.New(""),
			input:                     DashboardServiceDTO{},
			repoGetDashboardResult:    nil,
			repoGetDashboardResultErr: errors.New(""),
			clientResult:              nil,
			clientResultErr:           nil,
		},
		"should return err if ListMeasureEquipmentByDistributorId fail": {
			want:                      measures.DashboardResult{},
			wantErr:                   errors.New(""),
			input:                     DashboardServiceDTO{},
			repoGetDashboardResult:    make([]measures.DashboardMeasureI, 0),
			repoGetDashboardResultErr: nil,
			clientResult:              nil,
			clientResultErr:           errors.New(""),
		},
		"measureShouldBe": {
			want:                      fixtures.MeasureShouldBeWant,
			wantErr:                   nil,
			repoGetDashboardResult:    fixtures.MeasureShouldBeRepoGetDashboardResult,
			repoGetDashboardResultErr: nil,
			input:                     NewDashboardServiceDTO("test", time.Date(2022, 04, 01, 0, 0, 0, 0, loc), time.Date(2022, 04, 30, 0, 0, 0, 0, loc)),
			clientResult:              fixtures.MeasureShouldBeClientResult,
			clientResultErr:           nil,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			dashboardRepository := new(mocks.DashboardRepository)
			clientInventory := new(mocks2.Inventory)

			clientInventory.On("GroupMetersByType", mock.Anything, testCase.input.DistributorID).Return(testCase.clientResult, testCase.clientResultErr)

			dashboardRepository.On("GetDashboard", ctx, gross_measures.GetDashboardQuery{
				DistributorId: testCase.input.DistributorID,
				StartDate:     testCase.input.StartDate,
				EndDate:       testCase.input.EndDate,
			}).Return(testCase.repoGetDashboardResult, testCase.repoGetDashboardResultErr)

			h := NewDashboardService(dashboardRepository, clientInventory, loc)
			result, err := h.Handle(ctx, testCase.input)

			assert.Equal(t, testCase.want, result, testCase)
			assert.Equal(t, testCase.wantErr, err, testCase)
		})
	}
}
