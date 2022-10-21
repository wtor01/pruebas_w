package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Service_BillingMeasures_GetBillingSelfConsumptionByCau_Handler(t *testing.T) {

	type input struct {
		dto GetBillingSelfConsumptionDto
	}

	type output struct {
		GetBillingSelfConsumption    []billing_measures.BillingSelfConsumption
		GetBillingSelfConsumptionErr error
	}

	type want struct {
		err    error
		result []billing_measures.BillingSelfConsumptionUnit
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be bbdd error": {
			input: input{
				dto: NewGetBillingSelfConsumptionDto(
					"CAU_ID",
					"DistributorX",
					time.Date(2022, 6, 30, 22, 0, 0, 0, time.UTC),
					time.Date(2022, 7, 30, 22, 0, 0, 0, time.UTC),
				),
			},
			output: output{
				GetBillingSelfConsumption:    []billing_measures.BillingSelfConsumption{},
				GetBillingSelfConsumptionErr: errors.New("err bbdd"),
			},
			want: want{
				err:    errors.New("err bbdd"),
				result: []billing_measures.BillingSelfConsumptionUnit{},
			},
		},
		"Should be ok": {
			input: input{
				dto: NewGetBillingSelfConsumptionDto(
					"CAU_ID",
					"DistributorX",
					time.Date(2022, 6, 30, 22, 0, 0, 0, time.UTC),
					time.Date(2022, 7, 30, 22, 0, 0, 0, time.UTC),
				),
			},
			output: output{
				GetBillingSelfConsumption: []billing_measures.BillingSelfConsumption{
					{
						SurplusType:   billing_measures.DcSurplusType,
						CAU:           "CAU_ID",
						Name:          "CAU_NAME",
						DistributorId: "DistributorX",
						InitDate:      time.Date(2022, 6, 30, 22, 0, 0, 0, time.UTC),
						EndDate:       time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
						Status:        billing_measures.CalculatedSelfConsumptionStatus,
						Config: billing_measures.BillingSelfConsumptionConfig{
							CauID:        "CAU_ID",
							CnmcTypeDesc: "TypeDesc",
							ConfType:     "A",
						},
						Points: []billing_measures.BillingSelfConsumptionPoint{
							{
								ServicePointType: billing_measures.FronteraServicePointType,
								CUPS:             "CUPS",
								InitDate:         time.Date(2020, 1, 1, 23, 0, 0, 0, time.UTC),
								EndDate:          time.Date(2050, 1, 1, 23, 0, 0, 0, time.UTC),
							},
						},
						Curve: []billing_measures.BillingSelfConsumptionCurve{
							{
								EndDate: time.Date(2022, 6, 30, 23, 0, 0, 0, time.UTC),
								Points: []billing_measures.BillingSelfConsumptionCurvePoint{
									{
										CUPS: "CUPS",
										AI:   20,
										AE:   5,
									},
								},
								BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
									EHGN: utils.CreateFloat64(10),
									EHCR: utils.CreateFloat64(5),
									EHEX: utils.CreateFloat64(0),
									EHAU: utils.CreateFloat64(10),
									EHSA: utils.CreateFloat64(5),
									EHDC: utils.CreateFloat64(0),
								},
							},
							{
								EndDate: time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC),
								Points: []billing_measures.BillingSelfConsumptionCurvePoint{
									{
										CUPS: "CUPS",
										AI:   20,
										AE:   5,
									},
								},
								BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
									EHGN: utils.CreateFloat64(10),
									EHCR: utils.CreateFloat64(5),
									EHEX: utils.CreateFloat64(-10),
									EHAU: utils.CreateFloat64(10),
									EHSA: utils.CreateFloat64(5),
									EHDC: utils.CreateFloat64(0),
								},
							},
							{
								EndDate: time.Date(2022, 7, 1, 1, 0, 0, 0, time.UTC),
								Points: []billing_measures.BillingSelfConsumptionCurvePoint{
									{
										CUPS: "CUPS",
										AI:   20,
										AE:   5,
									},
								},
								BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
									EHGN: utils.CreateFloat64(10),
									EHCR: utils.CreateFloat64(5),
									EHEX: utils.CreateFloat64(2),
									EHAU: utils.CreateFloat64(10),
									EHSA: utils.CreateFloat64(5),
									EHDC: utils.CreateFloat64(0),
								},
							},
							{
								EndDate: time.Date(2022, 7, 1, 2, 0, 0, 0, time.UTC),
								Points: []billing_measures.BillingSelfConsumptionCurvePoint{
									{
										CUPS: "CUPS",
										AI:   20,
										AE:   5,
									},
								},
								BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
									EHGN: utils.CreateFloat64(10),
									EHCR: utils.CreateFloat64(5),
									EHEX: utils.CreateFloat64(5),
									EHAU: utils.CreateFloat64(10),
									EHSA: utils.CreateFloat64(5),
									EHDC: utils.CreateFloat64(0),
								},
							},
							{
								EndDate: time.Date(2022, 7, 1, 3, 0, 0, 0, time.UTC),
								Points: []billing_measures.BillingSelfConsumptionCurvePoint{
									{
										CUPS: "CUPS",
										AI:   20,
										AE:   5,
									},
								},
								BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
									EHGN: utils.CreateFloat64(10),
									EHCR: utils.CreateFloat64(5),
									EHEX: utils.CreateFloat64(20),
									EHAU: utils.CreateFloat64(10),
									EHSA: utils.CreateFloat64(5),
									EHDC: utils.CreateFloat64(0),
								},
							},
						},
						GraphHistory: &billing_measures.Graph{
							Dict:       nil,
							From:       nil,
							To:         nil,
							Algorithms: nil,
							StartedAt:  time.Time{},
							FinishedAt: time.Time{},
						},
					},
				},
			},
			want: want{
				result: []billing_measures.BillingSelfConsumptionUnit{
					{
						InitDate: time.Date(2022, 6, 30, 22, 0, 0, 0, time.UTC),
						EndDate:  time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
						Status:   billing_measures.CalculatedSelfConsumptionStatus,
						CauInfo: billing_measures.BillingSelfConsumptionCau{
							Id:         "CAU_ID",
							Name:       "CAU_NAME",
							Points:     1,
							UnitType:   "TypeDesc",
							ConfigType: "A",
						},
						Totals: billing_measures.BillingSelfConsumptionTotals{
							GrossGeneration:    75,
							NetGeneration:      50,
							SelfConsumption:    50,
							NetworkConsumption: 25,
							AuxConsumption:     25,
						},
						NetGeneration: []billing_measures.BillingSelfConsumptionNetGeneration{
							{
								Date:      "2022-07-01",
								Net:       50,
								Excedente: 17,
							},
						},
						UnitConsumption: []billing_measures.BillingSelfConsumptionUnitConsumption{
							{
								Date:    "2022-07-01",
								Network: 25,
								Self:    50,
								Aux:     25,
							},
						},
						CalendarConsumption: []billing_measures.BillingSelfConsumptionCalendarConsumption{
							{
								Date:   "2022-07-01",
								Energy: 17,
								Values: []billing_measures.CalendarConsumptionValues{
									{
										Hour:   "01:00",
										Energy: 0,
									},
									{
										Hour:   "02:00",
										Energy: -10,
									},
									{
										Hour:   "03:00",
										Energy: 2,
									},
									{
										Hour:   "04:00",
										Energy: 5,
									},
									{
										Hour:   "05:00",
										Energy: 20,
									},
								},
							},
						},
						Cups: []billing_measures.BillingSelfConsumptionCups{
							{
								Cups:        "CUPS",
								Type:        billing_measures.FronteraServicePointType,
								Consumption: 100,
								Generation:  25,
								StartDate:   "02-01-2020",
								EndDate:     "02-01-2050",
							},
						},
					},
				},
			},
		},
	}

	loc, _ := time.LoadLocation("Europe/Madrid")
	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			billingRepository := new(mocks.BillingSelfConsumptionRepository)
			billingRepository.On("GetSelfConsumptionByCau", mock.Anything, billing_measures.QueryGetBillingSelfConsumptionByCau{
				CauId:         test.input.dto.CauId,
				DistributorId: test.input.dto.DistributorId,
				StartDate:     test.input.dto.StartDate,
				EndDate:       test.input.dto.EndDate,
			}).Return(test.output.GetBillingSelfConsumption, test.output.GetBillingSelfConsumptionErr)
			srv := NewGetBillingSelfConsumptionByCauService(billingRepository, loc)

			result, err := srv.Handler(context.Background(), test.input.dto)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.result, result)
		})
	}
}
