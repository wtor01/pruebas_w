package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clients_mocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_GrossMeasure_Services_InsertMeasureBase_SetInsertMeasureMetadata(t *testing.T) {

	type input struct {
		measuresBase []gross_measures.GrossMeasureBase
	}

	type want struct {
		measuresBase []gross_measures.GrossMeasureBase
		result       error
	}

	type output struct {
		getDistributorByCdosResponse clients.Distributor
		getDistributorByCdosErr      error
	}
	endDate := time.Date(2021, 7, 14, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input  input
		want   want
		output output
	}{
		"Should set distributor data and generate id": {
			input: input{
				measuresBase: []gross_measures.GrossMeasureBase{
					&gross_measures.MeasureCloseWrite{
						EndDate:           endDate,
						Type:              measures.Incremental,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: "MeterSerialNumber",
						File:              "File",
						DistributorCDOS:   "CDOS",
					},
				},
			},
			want: want{
				measuresBase: []gross_measures.GrossMeasureBase{
					&gross_measures.MeasureCloseWrite{
						EndDate:           endDate,
						Type:              measures.Incremental,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: "MeterSerialNumber",
						File:              "File",
						DistributorCDOS:   "CDOS",
						DistributorID:     "DistributorID",
						Id:                fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s%s%s%s%s%s", "DistributorID", endDate.Format("2006-01-02_15:04:05"), "MeterSerialNumber", "File", measures.Incremental, measures.BillingClosure)))),
					},
				},
			},
			output: output{
				getDistributorByCdosResponse: clients.Distributor{
					ID: "DistributorID",
				},
			},
		},
		"Should return err if measures have diff DistributorCDOS": {
			input: input{
				measuresBase: []gross_measures.GrossMeasureBase{
					&gross_measures.MeasureCloseWrite{
						EndDate:           endDate,
						Type:              measures.Incremental,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: "MeterSerialNumber",
						File:              "File",
						DistributorCDOS:   "CDOS",
					},
					&gross_measures.MeasureCloseWrite{
						EndDate:           endDate,
						Type:              measures.Incremental,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: "MeterSerialNumber",
						File:              "File",
						DistributorCDOS:   "fail",
					},
				},
			},
			want: want{
				measuresBase: []gross_measures.GrossMeasureBase{
					&gross_measures.MeasureCloseWrite{
						EndDate:           endDate,
						Type:              measures.Incremental,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: "MeterSerialNumber",
						File:              "File",
						DistributorCDOS:   "CDOS",
					},
					&gross_measures.MeasureCloseWrite{
						EndDate:           endDate,
						Type:              measures.Incremental,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: "MeterSerialNumber",
						File:              "File",
						DistributorCDOS:   "fail",
					},
				},
				result: ErrDiffDistributors,
			},
			output: output{},
		},
		"Should return err if Distributor not found": {
			input: input{
				measuresBase: []gross_measures.GrossMeasureBase{
					&gross_measures.MeasureCloseWrite{
						EndDate:           endDate,
						Type:              measures.Incremental,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: "MeterSerialNumber",
						File:              "File",
						DistributorCDOS:   "CDOS",
					},
				},
			},
			want: want{
				measuresBase: []gross_measures.GrossMeasureBase{
					&gross_measures.MeasureCloseWrite{
						EndDate:           endDate,
						Type:              measures.Incremental,
						ReadingType:       measures.BillingClosure,
						MeterSerialNumber: "MeterSerialNumber",
						File:              "File",
						DistributorCDOS:   "CDOS",
					},
				},
				result: ErrDistributor,
			},
			output: output{
				getDistributorByCdosErr: ErrDistributor,
			},
		},
		"Should return err if empty measures": {
			input: input{
				measuresBase: []gross_measures.GrossMeasureBase{},
			},
			want: want{
				measuresBase: []gross_measures.GrossMeasureBase{},
				result:       ErrDistributor,
			},
			output: output{
				getDistributorByCdosErr: ErrDistributor,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			validationsClient := new(clients_mocks.Validation)
			inventoryClient := new(clients_mocks.Inventory)

			if len(testCase.input.measuresBase) != 0 {
				inventoryClient.On("GetDistributorByCdos", mock.Anything, testCase.input.measuresBase[0].GetDistributorCDOS()).Return(testCase.output.getDistributorByCdosResponse, testCase.output.getDistributorByCdosErr)
			}

			svc := InsertMeasureBase{
				validationClient: validationsClient,
				inventoryClient:  inventoryClient,
				generatorDate:    time.Now,
			}

			result := svc.SetInsertMeasureMetadata(ctx, testCase.input.measuresBase)

			assert.ElementsMatch(t, testCase.input.measuresBase, testCase.want.measuresBase)
			assert.Equal(t, testCase.want.result, result, "result")
		})
	}
}
