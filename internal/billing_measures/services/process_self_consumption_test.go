package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services/fixtures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/apperrors"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Services_Billing_Measures_ProcessSelfConsumption(t *testing.T) {

	type mocksResult struct {
		billingMeasureRepositoryFindErr error
		billingMeasureRepositoryFind    billing_measures.BillingMeasure

		selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr error
		selfConsumptionRepositoryGetActiveSelfConsumptionByCUP    billing_measures.SelfConsumption

		billingMeasureRepositorySaveErr   error
		billingMeasureRepositorySaveInput billing_measures.BillingMeasure

		GetBySelfConsumptionBetweenDates    billing_measures.BillingSelfConsumption
		GetBySelfConsumptionBetweenDatesErr error

		billingSelfConsumptionRepositorySaveErr   error
		billingSelfConsumptionRepositorySaveInput billing_measures.BillingSelfConsumption
	}

	tests := map[string]struct {
		input       billing_measures.OnSaveBillingMeasurePayload
		want        error
		mocksResult mocksResult
	}{
		"Should return operational error if no billing measure": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: apperrors.NewAppError("err", false),
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: errors.New("err"),
			},
		},
		"Should return operational error if GetActiveSelfConsumptionByCUP return error": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: apperrors.NewAppError("err", false),
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					CUPS: "CUP",
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: errors.New("err"),
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP:    billing_measures.SelfConsumption{},
			},
		},
		"Should not change billing measure status if GetActiveSelfConsumptionByCUP return empty selfConsumption": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: nil,
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					CUPS:   "CUP",
					Status: billing_measures.Calculated,
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: nil,
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP:    billing_measures.SelfConsumption{},
			},
		},
		"Should save billingMeasure with status PendingCau if selfConsumption exist": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: apperrors.NewAppError("err", false),
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					CUPS: "CUP",
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: errors.New("err"),
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP: billing_measures.SelfConsumption{
					ID: "ID",
				},
				billingMeasureRepositorySaveInput: billing_measures.BillingMeasure{
					CUPS:   "CUP",
					Status: billing_measures.PendingCau,
				},
			},
		},
		"Should return error if save billing measure fail": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: apperrors.NewAppError("err", true),
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					CUPS: "CUP",
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: nil,
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP: billing_measures.SelfConsumption{
					ID: "ID",
				},
				billingMeasureRepositorySaveInput: billing_measures.BillingMeasure{
					CUPS:   "CUP",
					Status: billing_measures.PendingCau,
				},
				billingMeasureRepositorySaveErr: errors.New("err"),
			},
		},
		"Should create with status pending if is not Redy To Process": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: nil,
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: nil,
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP:    fixtures.SelfConsumptionRepositoryGetActiveSelfConsumptionByCUP,
				billingMeasureRepositorySaveErr:                           nil,
				billingMeasureRepositorySaveInput: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Status:        billing_measures.PendingCau,
				},
				GetBySelfConsumptionBetweenDates:    billing_measures.BillingSelfConsumption{},
				GetBySelfConsumptionBetweenDatesErr: nil,

				billingSelfConsumptionRepositorySaveErr:   nil,
				billingSelfConsumptionRepositorySaveInput: fixtures.BillingSelfConsumptionRepositorySaveInput,
			},
		},
		"Should update with status pending if is not Redy To Process": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: nil,
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: nil,
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP:    fixtures.SelfConsumptionRepositoryGetActiveSelfConsumptionByCUP,

				billingMeasureRepositorySaveErr: nil,
				billingMeasureRepositorySaveInput: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Status:        billing_measures.PendingCau,
				},
				GetBySelfConsumptionBetweenDates:    fixtures.BillingSelfConsumptionRepositorySaveInput,
				GetBySelfConsumptionBetweenDatesErr: nil,

				billingSelfConsumptionRepositorySaveErr:   nil,
				billingSelfConsumptionRepositorySaveInput: fixtures.BillingSelfConsumptionRepositorySaveInput,
			},
		},
		"Should add billinMeasure to billingSelfConsumption and update with status pending if is not Redy To Process": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: nil,
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: nil,
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP:    fixtures.SelfConsumptionRepositoryGetActiveSelfConsumptionByCUP,

				billingMeasureRepositorySaveErr: nil,
				billingMeasureRepositorySaveInput: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Status:        billing_measures.PendingCau,
				},
				GetBySelfConsumptionBetweenDates:    fixtures.GetBySelfConsumptionBetweenDates,
				GetBySelfConsumptionBetweenDatesErr: nil,

				billingSelfConsumptionRepositorySaveErr:   nil,
				billingSelfConsumptionRepositorySaveInput: fixtures.BillingSelfConsumptionRepositorySaveInput,
			},
		},
		"Should return error if billingSelfConsumptionRepository.Save fail": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: apperrors.NewAppError("err", true),
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: nil,
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP:    fixtures.SelfConsumptionRepositoryGetActiveSelfConsumptionByCUP,

				billingMeasureRepositorySaveErr: nil,
				billingMeasureRepositorySaveInput: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Status:        billing_measures.PendingCau,
				},
				GetBySelfConsumptionBetweenDates:    billing_measures.BillingSelfConsumption{},
				GetBySelfConsumptionBetweenDatesErr: nil,

				billingSelfConsumptionRepositorySaveErr:   errors.New("err"),
				billingSelfConsumptionRepositorySaveInput: fixtures.BillingSelfConsumptionRepositorySaveInput,
			},
		},
		"Should return error if billingSelfConsumptionRepository.Save when start process": {
			input: billing_measures.OnSaveBillingMeasurePayload{
				InitDate:         time.Time{},
				EndDate:          time.Time{},
				CUPS:             "CUP",
				BillingMeasureId: "ID",
			},
			want: apperrors.NewAppError("err", true),
			mocksResult: mocksResult{
				billingMeasureRepositoryFindErr: nil,
				billingMeasureRepositoryFind: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				},
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr: nil,
				selfConsumptionRepositoryGetActiveSelfConsumptionByCUP:    fixtures.SelfConsumptionRepositoryGetActiveSelfConsumptionByCUPonePoint,

				billingMeasureRepositorySaveErr: nil,
				billingMeasureRepositorySaveInput: billing_measures.BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Status:        billing_measures.PendingCau,
				},
				GetBySelfConsumptionBetweenDates:    billing_measures.BillingSelfConsumption{},
				GetBySelfConsumptionBetweenDatesErr: nil,

				billingSelfConsumptionRepositorySaveErr:   errors.New("err"),
				billingSelfConsumptionRepositorySaveInput: fixtures.BillingSelfConsumptionRepositorySaveInputStartedonPoint,
			},
		},
	}

	for testsName, _ := range tests {
		testCase := tests[testsName]
		t.Run(testsName, func(t *testing.T) {
			t.Parallel()
			billingMeasureRepository := new(mocks.BillingMeasureRepository)
			selfConsumptionRepository := new(mocks.SelfConsumptionRepository)
			billingSelfConsumptionRepository := new(mocks.BillingSelfConsumptionRepository)

			billingMeasureRepository.On(
				"Find", mock.Anything,
				testCase.input.BillingMeasureId,
			).Return(
				testCase.mocksResult.billingMeasureRepositoryFind,
				testCase.mocksResult.billingMeasureRepositoryFindErr,
			)

			billingMeasureRepository.On(
				"Save", mock.Anything,
				testCase.mocksResult.billingMeasureRepositorySaveInput,
			).Return(
				testCase.mocksResult.billingMeasureRepositorySaveErr,
			)

			selfConsumptionRepository.On(
				"GetActiveSelfConsumptionByCUP",
				mock.Anything,
				billing_measures.QueryGetActiveSelfConsumptionByCUP{
					CUP:       testCase.mocksResult.billingMeasureRepositoryFind.CUPS,
					StartDate: testCase.mocksResult.billingMeasureRepositoryFind.InitDate,
					EndDate:   testCase.mocksResult.billingMeasureRepositoryFind.EndDate,
				}).Return(
				testCase.mocksResult.selfConsumptionRepositoryGetActiveSelfConsumptionByCUP,
				testCase.mocksResult.selfConsumptionRepositoryGetActiveSelfConsumptionByCUPErr,
			)

			billingSelfConsumptionRepository.On(
				"GetBySelfConsumptionBetweenDates",
				mock.Anything,
				billing_measures.QueryGetBillingSelfConsumption{
					SelfconsumptionId: testCase.mocksResult.selfConsumptionRepositoryGetActiveSelfConsumptionByCUP.ID,
					StartDate:         testCase.mocksResult.billingMeasureRepositoryFind.InitDate,
					EndDate:           testCase.mocksResult.billingMeasureRepositoryFind.EndDate,
				}).Return(
				testCase.mocksResult.GetBySelfConsumptionBetweenDates,
				testCase.mocksResult.GetBySelfConsumptionBetweenDatesErr,
			)

			billingSelfConsumptionRepository.On(
				"Save",
				mock.Anything,
				mock.MatchedBy(func(compare billing_measures.BillingSelfConsumption) bool {
					testCase.mocksResult.billingSelfConsumptionRepositorySaveInput.Id = compare.Id
					testCase.mocksResult.billingSelfConsumptionRepositorySaveInput.GenerationDate = compare.GenerationDate
					return assert.Equal(t, compare, testCase.mocksResult.billingSelfConsumptionRepositorySaveInput)
				}),
			).Return(
				testCase.mocksResult.billingSelfConsumptionRepositorySaveErr,
			)

			location, _ := time.LoadLocation("Europe/Madrid")
			srv := NewProcessSelfConsumption(billingMeasureRepository, selfConsumptionRepository, billingSelfConsumptionRepository, location, nil)
			err := srv.Handler(context.Background(), testCase.input)

			assert.Equal(t, testCase.want, err)
		})

	}
}
