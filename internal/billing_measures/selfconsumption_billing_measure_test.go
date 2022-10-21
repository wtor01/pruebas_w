package billing_measures

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_BillingMeasures_BillingSelfConsumption_NewBillingSelfConsumption(t *testing.T) {
	type input struct {
		selfConsumption SelfConsumption
		billingMeasure  BillingMeasure
	}
	type want struct {
		result BillingSelfConsumption
		err    error
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"Should return ok": {
			input: input{
				selfConsumption: SelfConsumption{
					ID:            "ID",
					CAU:           "CAU",
					Name:          "Name",
					StatusID:      0,
					StatusName:    "",
					CcaaId:        0,
					Ccaa:          "",
					InitDate:      time.Time{},
					EndDate:       time.Time{},
					DistributorId: "",
					Configs: []ConfigSelfConsumption{
						{
							ID:                    "Id",
							CauID:                 "CauID",
							StatusID:              1,
							StatusName:            "StatusName",
							InitDate:              time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
							EndDate:               time.Date(2022, 02, 01, 00, 00, 0, 0, time.UTC),
							CnmcTypeId:            2,
							CnmcTypeName:          "CnmcTypeName",
							CnmcTypeDesc:          "CnmcTypeDesc",
							ConfType:              "ConfType",
							ConfTypeDescription:   "ConfTypeDescription",
							ConsumerType:          "ConsumerType",
							ParticipantNumber:     3,
							ConnType:              "ConnType",
							Excedents:             false,
							Compensation:          true,
							GenerationPot:         4,
							GroupSubgroup:         5,
							AntivertType:          "AntivertType",
							SolarZoneId:           6,
							SolarZoneNum:          7,
							SolarZoneName:         "SolarZoneName",
							TechnologyId:          "TechnologyId",
							TechnologyDescription: "TechnologyDescription",
						},
						{
							ID:       "Id 2",
							InitDate: time.Date(2022, 02, 01, 00, 00, 0, 0, time.UTC),
							EndDate:  time.Date(2022, 03, 01, 00, 00, 0, 0, time.UTC),
						},
					},
					Points: []PointSelfConsumption{
						{
							ID:               "Id 1",
							ServicePointType: "ServicePointType",
							CUPS:             "CUPS_1",
							InstalationFlag:  1,
							WithoutmeterFlag: 2,
							Exent1Flag:       3,
							Exent2Flag:       4,
							PartitionCoeff:   5,
							InitDate:         time.Date(2022, 01, 01, 01, 00, 0, 0, time.UTC),
							EndDate:          time.Date(2022, 01, 01, 02, 00, 0, 0, time.UTC),
						},
						{
							ID:               "Id 2",
							ServicePointType: "ServicePointType",
							CUPS:             "CUPS_2",
							InstalationFlag:  1,
							WithoutmeterFlag: 2,
							Exent1Flag:       3,
							Exent2Flag:       4,
							PartitionCoeff:   5,
							InitDate:         time.Date(2022, 01, 01, 01, 00, 0, 0, time.UTC),
							EndDate:          time.Date(2022, 01, 01, 02, 00, 0, 0, time.UTC),
						},
					},
				},
				billingMeasure: BillingMeasure{
					Id:            "BillingMeasureId",
					DistributorID: "DistributorID",
					CUPS:          "CUPS_1",
					InitDate:      time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:       time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				},
			},
			want: want{
				result: BillingSelfConsumption{
					SelfconsumptionId: "ID",
					CAU:               "CAU",
					Name:              "Name",
					DistributorId:     "",
					InitDate:          time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:           time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Status:            PendingPoints,
					Config: BillingSelfConsumptionConfig{
						Id:                    "Id",
						CauID:                 "CauID",
						StatusID:              1,
						StatusName:            "StatusName",
						InitDate:              time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
						EndDate:               time.Date(2022, 02, 01, 00, 00, 0, 0, time.UTC),
						CnmcTypeId:            2,
						CnmcTypeName:          "CnmcTypeName",
						CnmcTypeDesc:          "CnmcTypeDesc",
						ConfType:              "ConfType",
						ConfTypeDescription:   "ConfTypeDescription",
						ConsumerType:          "ConsumerType",
						ParticipantNumber:     3,
						ConnType:              "ConnType",
						Excedents:             false,
						Compensation:          true,
						GenerationPot:         4,
						GroupSubgroup:         5,
						AntivertType:          "AntivertType",
						SolarZoneId:           6,
						SolarZoneNum:          7,
						SolarZoneName:         "SolarZoneName",
						TechnologyId:          "TechnologyId",
						TechnologyDescription: "TechnologyDescription",
					},
					Points: []BillingSelfConsumptionPoint{
						{
							ID:               "Id 1",
							ServicePointType: "ServicePointType",
							CUPS:             "CUPS_1",
							InstalationFlag:  1,
							WithoutmeterFlag: 2,
							Exent1Flag:       3,
							Exent2Flag:       4,
							PartitionCoeff:   5,
							InitDate:         time.Date(2022, 01, 01, 01, 00, 0, 0, time.UTC),
							EndDate:          time.Date(2022, 01, 01, 02, 00, 0, 0, time.UTC),
							MvhReceived: []BillingSelfConsumptionPointMvhReceived{
								{
									BillingMeasureId: "BillingMeasureId",
									InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
									EndDate:          time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
								},
							},
						},
						{
							ID:               "Id 2",
							ServicePointType: "ServicePointType",
							CUPS:             "CUPS_2",
							InstalationFlag:  1,
							WithoutmeterFlag: 2,
							Exent1Flag:       3,
							Exent2Flag:       4,
							PartitionCoeff:   5,
							InitDate:         time.Date(2022, 01, 01, 01, 00, 0, 0, time.UTC),
							EndDate:          time.Date(2022, 01, 01, 02, 00, 0, 0, time.UTC),
							MvhReceived:      make([]BillingSelfConsumptionPointMvhReceived, 0),
						},
					},
					Curve: make([]BillingSelfConsumptionCurve, 0),
				},
				err: nil,
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			result, err := NewBillingSelfConsumption(testCase.input.selfConsumption, testCase.input.billingMeasure)

			assert.NotEmpty(t, result.Id)
			assert.NotEmpty(t, result.GenerationDate)

			testCase.want.result.Id = result.Id
			testCase.want.result.GenerationDate = result.GenerationDate

			assert.Equal(t, testCase.want.result, result)
			assert.Equal(t, testCase.want.err, err)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_BillingSelfConsumption_SetDates(t *testing.T) {
	type input struct {
		billingSelfConsumption BillingSelfConsumption
		billingMeasure         BillingMeasure
	}

	tests := map[string]struct {
		input input
		want  BillingSelfConsumption
	}{
		"Should not update dates": {
			input: input{
				billingSelfConsumption: BillingSelfConsumption{
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				},
				billingMeasure: BillingMeasure{
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
				},
			},
			want: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
			},
		},
		"Should update dates": {
			input: input{
				billingSelfConsumption: BillingSelfConsumption{
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
				},
				billingMeasure: BillingMeasure{
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				},
			},
			want: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			testCase.input.billingSelfConsumption.SetDates(testCase.input.billingMeasure)

			assert.Equal(t, testCase.want, testCase.input.billingSelfConsumption)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_BillingSelfConsumption_AddBillingMeasure(t *testing.T) {
	type input struct {
		billingSelfConsumption BillingSelfConsumption
		billingMeasure         BillingMeasure
	}

	tests := map[string]struct {
		input input
		want  BillingSelfConsumption
	}{
		"Should add billing measure": {
			input: input{
				billingSelfConsumption: BillingSelfConsumption{
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Points: []BillingSelfConsumptionPoint{
						{

							CUPS:        "CUP_1",
							MvhReceived: []BillingSelfConsumptionPointMvhReceived{},
						},
					},
				},
				billingMeasure: BillingMeasure{
					Id:       "billingMeasure id",
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
					CUPS:     "CUP_1",
					BillingLoadCurve: []BillingLoadCurve{
						{
							EndDate: time.Date(2022, 01, 01, 00, 00, 00, 00, time.UTC),
							AI:      1,
							AE:      2,
						},
					},
				},
			},
			want: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				Points: []BillingSelfConsumptionPoint{
					{

						CUPS: "CUP_1",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{
							{
								BillingMeasureId: "billingMeasure id",
								InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
							},
						},
					},
				},
				Curve: []BillingSelfConsumptionCurve{
					{
						EndDate: time.Date(2022, 01, 01, 00, 00, 00, 00, time.UTC),
						Points: []BillingSelfConsumptionCurvePoint{
							{
								CUPS: "CUP_1",
								AI:   1,
								AE:   2,
							},
						},
					},
				},
			},
		},
		"Should add billing measure and update curves": {
			input: input{
				billingSelfConsumption: BillingSelfConsumption{
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Points: []BillingSelfConsumptionPoint{
						{

							CUPS:        "CUP_1",
							MvhReceived: []BillingSelfConsumptionPointMvhReceived{},
						},
					},
					Curve: []BillingSelfConsumptionCurve{
						{
							EndDate: time.Date(2022, 01, 01, 00, 00, 00, 00, time.UTC),
							Points: []BillingSelfConsumptionCurvePoint{
								{
									CUPS: "CUP_1",
									AI:   1,
									AE:   2,
								},
							},
						},
					},
				},
				billingMeasure: BillingMeasure{
					Id:       "billingMeasure id",
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
					CUPS:     "CUP_1",
					BillingLoadCurve: []BillingLoadCurve{
						{
							EndDate: time.Date(2022, 01, 01, 00, 00, 00, 00, time.UTC),
							AI:      10,
							AE:      20,
						},
						{
							EndDate: time.Date(2022, 01, 02, 00, 00, 00, 00, time.UTC),
							AI:      1,
							AE:      2,
						},
					},
				},
			},
			want: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				Points: []BillingSelfConsumptionPoint{
					{

						CUPS: "CUP_1",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{
							{
								BillingMeasureId: "billingMeasure id",
								InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
							},
						},
					},
				},
				Curve: []BillingSelfConsumptionCurve{
					{
						EndDate: time.Date(2022, 01, 01, 00, 00, 00, 00, time.UTC),
						Points: []BillingSelfConsumptionCurvePoint{
							{
								CUPS: "CUP_1",
								AI:   10,
								AE:   20,
							},
						},
					},
					{
						EndDate: time.Date(2022, 01, 02, 00, 00, 00, 00, time.UTC),
						Points: []BillingSelfConsumptionCurvePoint{
							{
								CUPS: "CUP_1",
								AI:   1,
								AE:   2,
							},
						},
					},
				},
			},
		},
		"Should add billing measure in existing MvhReceived and update dates": {
			input: input{
				billingSelfConsumption: BillingSelfConsumption{
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
					Points: []BillingSelfConsumptionPoint{
						{

							CUPS: "CUP_1",
							MvhReceived: []BillingSelfConsumptionPointMvhReceived{
								{
									BillingMeasureId: "billingMeasure id",
									InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
									EndDate:          time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				billingMeasure: BillingMeasure{
					Id:       "billingMeasure id 2",
					InitDate: time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					CUPS:     "CUP_1",
				},
			},
			want: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				Points: []BillingSelfConsumptionPoint{
					{

						CUPS: "CUP_1",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{
							{
								BillingMeasureId: "billingMeasure id",
								InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
							},
							{
								BillingMeasureId: "billingMeasure id 2",
								InitDate:         time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
							},
						},
					},
				},
			},
		},
		"Should add billing measure with diff cups": {
			input: input{
				billingSelfConsumption: BillingSelfConsumption{
					InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					Points: []BillingSelfConsumptionPoint{
						{

							CUPS: "CUP_1",
							MvhReceived: []BillingSelfConsumptionPointMvhReceived{
								{
									BillingMeasureId: "billingMeasure id",
									InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
									EndDate:          time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
								},
							},
						},
						{

							CUPS:        "CUP_2",
							MvhReceived: []BillingSelfConsumptionPointMvhReceived{},
						},
					},
				},
				billingMeasure: BillingMeasure{
					Id:       "billingMeasure id 2",
					InitDate: time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
					EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
					CUPS:     "CUP_2",
				},
			},
			want: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				Points: []BillingSelfConsumptionPoint{
					{

						CUPS: "CUP_1",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{
							{
								BillingMeasureId: "billingMeasure id",
								InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
							},
						},
					},
					{

						CUPS: "CUP_2",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{
							{
								BillingMeasureId: "billingMeasure id 2",
								InitDate:         time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
							},
						},
					},
				},
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			testCase.input.billingSelfConsumption.AddBillingMeasure(testCase.input.billingMeasure)

			assert.Equal(t, testCase.want, testCase.input.billingSelfConsumption)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_BillingSelfConsumption_IsRedyToProcess(t *testing.T) {
	tests := map[string]struct {
		input BillingSelfConsumption
		want  bool
	}{
		"Should return false 1 cup incomplete mvh received": {
			input: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				Points: []BillingSelfConsumptionPoint{
					{

						CUPS: "CUP_1",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{
							{
								BillingMeasureId: "billingMeasure id",
								InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
							},
						},
					},
				},
			},
			want: false,
		},
		"Should return false 0 cup": {
			input: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				Points: []BillingSelfConsumptionPoint{
					{

						CUPS:        "CUP_1",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{},
					},
				},
			},
			want: false,
		},
		"Should return true": {
			input: BillingSelfConsumption{
				InitDate: time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
				EndDate:  time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
				Points: []BillingSelfConsumptionPoint{
					{

						CUPS: "CUP_1",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{
							{
								BillingMeasureId: "billingMeasure id",
								InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
							},
							{
								BillingMeasureId: "billingMeasure id",
								InitDate:         time.Date(2022, 01, 15, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
							},
						},
					},
					{

						CUPS: "CUP_2",
						MvhReceived: []BillingSelfConsumptionPointMvhReceived{
							{
								BillingMeasureId: "billingMeasure id 2",
								InitDate:         time.Date(2022, 01, 01, 00, 00, 0, 0, time.UTC),
								EndDate:          time.Date(2022, 01, 31, 00, 00, 0, 0, time.UTC),
							},
						},
					},
				},
			},
			want: true,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			assert.Equal(t, testCase.want, testCase.input.IsRedyToProcess())
		})
	}
}
