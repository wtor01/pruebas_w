package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_BillingMeasures_IsValidBillingBalancePeriod(t *testing.T) {

	type input struct {
		b      BillingMeasure
		period measures.PeriodKey
	}

	tests := map[string]struct {
		input input
		want  bool
	}{
		"should return false with invalid Period": {
			input: input{
				b:      BillingMeasure{},
				period: "fail",
			},
			want: false,
		},
		"should return false if Period is nil": {
			input: input{
				b:      BillingMeasure{},
				period: measures.P1,
			},
			want: false,
		},
		"should return false if Period BalanceValidation is invalid": {
			input: input{
				b: BillingMeasure{
					BillingBalance: BillingBalance{
						P1: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
					},
				},
				period: measures.P1,
			},
			want: false,
		},
		"should return true if P0 Period BalanceValidation is valid": {
			input: input{
				b: BillingMeasure{
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{
							BalanceValidationAI: measures.Valid,
						},
					},
				},
				period: measures.P0,
			},
			want: true,
		},
		"should return true if P1 Period BalanceValidation is valid": {
			input: input{
				b: BillingMeasure{
					BillingBalance: BillingBalance{
						P1: &BillingBalancePeriod{
							BalanceValidationAI: measures.Valid,
						},
					},
				},
				period: measures.P1,
			},
			want: true,
		},
		"should return true if P2 Period BalanceValidation is valid": {
			input: input{
				b: BillingMeasure{
					BillingBalance: BillingBalance{
						P2: &BillingBalancePeriod{
							BalanceValidationAI: measures.Valid,
						},
					},
				},
				period: measures.P2,
			},
			want: true,
		},
		"should return true if P3 Period BalanceValidation is valid": {
			input: input{
				b: BillingMeasure{
					BillingBalance: BillingBalance{
						P3: &BillingBalancePeriod{
							BalanceValidationAI: measures.Valid,
						},
					},
				},
				period: measures.P3,
			},
			want: true,
		},
		"should return true if P4 Period BalanceValidation is valid": {
			input: input{
				b: BillingMeasure{
					BillingBalance: BillingBalance{
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Valid,
						},
					},
				},
				period: measures.P4,
			},
			want: true,
		},
		"should return true if P5 Period BalanceValidation is valid": {
			input: input{
				b: BillingMeasure{
					BillingBalance: BillingBalance{
						P5: &BillingBalancePeriod{
							BalanceValidationAI: measures.Valid,
						},
					},
				},
				period: measures.P5,
			},
			want: true,
		},
		"should return true if P6 Period BalanceValidation is valid": {
			input: input{
				b: BillingMeasure{
					BillingBalance: BillingBalance{
						P6: &BillingBalancePeriod{
							BalanceValidationAI: measures.Valid,
						},
					},
				},
				period: measures.P6,
			},
			want: true,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			r := testCase.input.b.IsValidBillingBalancePeriod(testCase.input.period, measures.AI)
			assert.Equal(t, testCase.want, r, testName)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsAtrVsCurveValid(t *testing.T) {

	type input struct {
		b         BillingMeasure
		period    measures.PeriodKey
		magnitude measures.Magnitude
	}

	tests := map[string]struct {
		input input
		want  bool
	}{
		"should return false with invalid Period": {
			input: input{
				b:      BillingMeasure{},
				period: "fail",
			},
			want: false,
		},
		"should return false if Period is nil": {
			input: input{
				b:         BillingMeasure{},
				period:    measures.P1,
				magnitude: measures.AI,
			},
			want: false,
		},
		"should return false if Period AtrVsCurve is > 1KW": {
			input: input{
				b: BillingMeasure{
					AtrVsCurve: AtrVsCurve{
						P1: &measures.Values{
							AI: 2000,
						},
					},
				},
				period:    measures.P1,
				magnitude: measures.AI,
			},
			want: false,
		},
		"should return true if P0 Period AtrVsCurve is < 1KW": {
			input: input{
				b: BillingMeasure{
					AtrVsCurve: AtrVsCurve{
						P0: &measures.Values{
							AI: 1,
						},
					},
				},
				period:    measures.P0,
				magnitude: measures.AI,
			},
			want: true,
		},
		"should return true if P1 Period AtrVsCurve is < 1KW": {
			input: input{
				b: BillingMeasure{
					AtrVsCurve: AtrVsCurve{
						P1: &measures.Values{
							AI: 1,
						},
					},
				},
				period:    measures.P1,
				magnitude: measures.AI,
			},
			want: true,
		},
		"should return true if P2 Period AtrVsCurve is < 1KW": {
			input: input{
				b: BillingMeasure{
					AtrVsCurve: AtrVsCurve{
						P2: &measures.Values{
							AI: 1,
						},
					},
				},
				period:    measures.P2,
				magnitude: measures.AI,
			},
			want: true,
		},
		"should return true if P3 Period AtrVsCurve is < 1KW": {
			input: input{
				b: BillingMeasure{
					AtrVsCurve: AtrVsCurve{
						P3: &measures.Values{
							AI: 1,
						},
					},
				},
				period:    measures.P3,
				magnitude: measures.AI,
			},
			want: true,
		},
		"should return true if P4 Period AtrVsCurve is < 1KW": {
			input: input{
				b: BillingMeasure{
					AtrVsCurve: AtrVsCurve{
						P4: &measures.Values{
							AI: 1,
						},
					},
				},
				period:    measures.P4,
				magnitude: measures.AI,
			},
			want: true,
		},
		"should return true if P5 Period AtrVsCurve is < 1KW": {
			input: input{
				b: BillingMeasure{
					AtrVsCurve: AtrVsCurve{
						P5: &measures.Values{
							AI: 1,
						},
					},
				},
				period:    measures.P5,
				magnitude: measures.AI,
			},
			want: true,
		},
		"should return true if P6 Period AtrVsCurve is < 1KW": {
			input: input{
				b: BillingMeasure{
					AtrVsCurve: AtrVsCurve{
						P6: &measures.Values{
							AI: 1,
						},
					},
				},
				period:    measures.P6,
				magnitude: measures.AI,
			},
			want: true,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			r := testCase.input.b.IsAtrVsCurveValid(testCase.input.period, testCase.input.magnitude)
			assert.Equal(t, testCase.want, r, testName)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_BillingBalancePeriod_CalcAtrBalance(t *testing.T) {

	type input struct {
		b                BillingBalancePeriod
		actual, previous *measures.DailyReadingClosureCalendarPeriod
	}

	tests := map[string]struct {
		input input
		want  BillingBalancePeriod
	}{
		"should set measures.Invalid if actual and previous are nil": {
			input: input{
				b:        BillingBalancePeriod{},
				actual:   nil,
				previous: nil,
			},
			want: BillingBalancePeriod{
				BalanceValidationAI: measures.Invalid,
			},
		},
		"should set measures.Invalid if previous is nil": {
			input: input{
				b:        BillingBalancePeriod{},
				actual:   &measures.DailyReadingClosureCalendarPeriod{},
				previous: nil,
			},
			want: BillingBalancePeriod{
				BalanceValidationAI: measures.Invalid,
			},
		},
		"should set measures.Invalid if actual is nil": {
			input: input{
				b:        BillingBalancePeriod{},
				actual:   nil,
				previous: &measures.DailyReadingClosureCalendarPeriod{},
			},
			want: BillingBalancePeriod{
				BalanceValidationAI: measures.Invalid,
			},
		},
		"should set measures.Valid": {
			input: input{
				b: BillingBalancePeriod{},
				actual: &measures.DailyReadingClosureCalendarPeriod{
					ValidationStatus: measures.Valid,
				},
				previous: &measures.DailyReadingClosureCalendarPeriod{
					ValidationStatus: measures.Valid,
				},
			},
			want: BillingBalancePeriod{
				BalanceValidationAI: measures.Valid,
			},
		},
		"should set measures.Valid and subtraction actual - previous": {
			input: input{
				b: BillingBalancePeriod{},
				actual: &measures.DailyReadingClosureCalendarPeriod{
					Values: measures.Values{
						AI: 9,
						AE: 8,
						R1: 7,
						R2: 6,
						R3: 5,
						R4: 4,
					},
					ValidationStatus: measures.Valid,
				},
				previous: &measures.DailyReadingClosureCalendarPeriod{
					Values: measures.Values{
						AI: 8,
						AE: 5,
						R1: 5,
						R2: 0,
						R3: 5,
						R4: 1,
					},
					ValidationStatus: measures.Valid,
				},
			},
			want: BillingBalancePeriod{
				AI:                  1,
				BalanceValidationAI: measures.Valid,
			},
		},
		"should set measures.Invalid if actual is filled": {
			input: input{
				b: BillingBalancePeriod{},
				actual: &measures.DailyReadingClosureCalendarPeriod{
					Filled: true,
				},
				previous: &measures.DailyReadingClosureCalendarPeriod{},
			},
			want: BillingBalancePeriod{
				BalanceValidationAI: measures.Invalid,
			},
		},
		"should set measures.Invalid if previous is filled": {
			input: input{
				b:      BillingBalancePeriod{},
				actual: &measures.DailyReadingClosureCalendarPeriod{},
				previous: &measures.DailyReadingClosureCalendarPeriod{
					Filled: true,
				},
			},
			want: BillingBalancePeriod{
				BalanceValidationAI: measures.Invalid,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			testCase.input.b.CalcAtrBalance(testCase.input.actual, testCase.input.previous, measures.AI)
			assert.Equal(t, testCase.want, testCase.input.b, testName)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_BillingBalancePeriod_SetAtrBalance(t *testing.T) {
	type input struct {
		b                BillingBalancePeriod
		actual, previous *measures.DailyReadingClosureCalendarPeriod
	}

	tests := map[string]struct {
		input input
		want  BillingBalancePeriod
	}{
		"should set measures.Invalid if actual is filled": {
			input: input{
				b: BillingBalancePeriod{},
				actual: &measures.DailyReadingClosureCalendarPeriod{
					Filled: true,
				},
			},
			want: BillingBalancePeriod{
				BalanceValidationAI: measures.Invalid,
			},
		},
		"should set measures.Invalid if actual is nil": {
			input: input{
				b: BillingBalancePeriod{},
			},
			want: BillingBalancePeriod{
				BalanceValidationAI: measures.Invalid,
			},
		},
		"should set measures.Valid if actual is correct": {
			input: input{
				b: BillingBalancePeriod{},
				actual: &measures.DailyReadingClosureCalendarPeriod{
					Values: measures.Values{
						AI: 10000,
					},
					ValidationStatus: measures.Valid,
				},
			},
			want: BillingBalancePeriod{
				AI:                  10000,
				BalanceValidationAI: measures.Valid,
			},
		},
	}

	for name, _ := range tests {
		test := tests[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			test.input.b.SetAtrBalance(test.input.actual, measures.AI)
			assert.Equal(t, test.want, test.input.b)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsCurvePeriodCompletedByPeriod(t *testing.T) {

	type input struct {
		b      func() BillingMeasure
		period measures.PeriodKey
	}

	tests := map[string]struct {
		input input
		want  bool
	}{
		"should return false if have filled measure curve for this period": {
			input: input{
				b: func() BillingMeasure {
					b := NewBillingMeasure("",
						time.Date(2022, 07, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2022, 07, 2, 3, 0, 0, 0, time.UTC),
						"",
						"",
						[]measures.PeriodKey{measures.P1},
						[]measures.Magnitude{measures.AE},
						measures.TLG,
					)
					b.SetBillingLoadCurve([]BillingLoadCurve{
						{
							EndDate: time.Date(2022, 07, 1, 1, 0, 0, 0, time.UTC),
							Period:  measures.P1,
						},
						{
							EndDate: time.Date(2022, 07, 1, 1, 0, 0, 0, time.UTC),
							Period:  measures.P1,
							Origin:  measures.Filled,
						},
					})

					return b
				},
				period: measures.P1,
			},
			want: false,
		},
		"should return true if not have filled measure curve for this period": {
			input: input{
				b: func() BillingMeasure {
					b := NewBillingMeasure("",
						time.Date(2022, 07, 1, 24, 0, 0, 0, time.UTC),
						time.Date(2022, 07, 2, 3, 0, 0, 0, time.UTC),
						"",
						"",
						[]measures.PeriodKey{measures.P1},
						[]measures.Magnitude{measures.AE},
						measures.TLG,
					)
					b.SetBillingLoadCurve([]BillingLoadCurve{
						{
							EndDate: time.Date(2022, 07, 1, 1, 0, 0, 0, time.UTC),
							Period:  measures.P1,
						},
					})
					return b
				},
				period: measures.P1,
			},
			want: true,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			assert.Equal(t, testCase.want, testCase.input.b().IsCurvePeriodCompletedByPeriod(testCase.input.period), testName)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_NewBillingMeasure(t *testing.T) {
	type input struct {
		cups            string
		initDate        time.Time
		endDate         time.Time
		distributorCode string
		distributorId   string
		Periods         []measures.PeriodKey
		Magnitudes      []measures.Magnitude
		MeterType       measures.MeterType
	}
	type test struct {
		input input
		want  BillingMeasure
	}
	var loc, _ = time.LoadLocation("Europe/Madrid")

	tests := map[string]test{
		"should return ok for 1 period": {
			input: input{
				cups:            "CUPS",
				initDate:        time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				endDate:         time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				distributorCode: "0130",
				distributorId:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				Periods:         []measures.PeriodKey{measures.P1},
				Magnitudes:      []measures.Magnitude{measures.AE},
				MeterType:       measures.TLG,
			},
			want: BillingMeasure{
				DistributorCode:        "0130",
				DistributorID:          "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:                   "CUPS",
				PointType:              "",
				RegisterType:           "",
				GenerationDate:         time.Time{},
				EndDate:                time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				InitDate:               time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				Version:                "0",
				Origin:                 "",
				Status:                 Calculating,
				PreviousReadingClosure: measures.DailyReadingClosure{},
				ActualReadingClosure:   measures.DailyReadingClosure{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
				},
				BillingLoadCurve: nil,
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{},
					P1: &measures.Values{},
				},
				Periods:      []measures.PeriodKey{measures.P1},
				GraphHistory: map[string]*Graph{},
				Magnitudes:   []measures.Magnitude{measures.AE},
				MeterType:    measures.TLG,
			},
		},
		"should return ok for 2 period": {
			input: input{
				cups:            "CUPS",
				initDate:        time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				endDate:         time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				distributorCode: "0130",
				distributorId:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				Periods:         []measures.PeriodKey{measures.P1, measures.P2},
				Magnitudes:      []measures.Magnitude{measures.AE},
				MeterType:       measures.TLG,
			},
			want: BillingMeasure{
				DistributorCode:        "0130",
				DistributorID:          "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:                   "CUPS",
				PointType:              "",
				RegisterType:           "",
				GenerationDate:         time.Time{},
				EndDate:                time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				InitDate:               time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				Version:                "0",
				Origin:                 "",
				Status:                 Calculating,
				PreviousReadingClosure: measures.DailyReadingClosure{},
				ActualReadingClosure:   measures.DailyReadingClosure{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
				},
				BillingLoadCurve: nil,
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{},
					P1: &measures.Values{},
					P2: &measures.Values{},
				},
				Periods:      []measures.PeriodKey{measures.P1, measures.P2},
				GraphHistory: map[string]*Graph{},
				Magnitudes:   []measures.Magnitude{measures.AE},
				MeterType:    measures.TLG,
			},
		},
		"should return ok for 3 period": {
			input: input{
				cups:            "CUPS",
				initDate:        time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				endDate:         time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				distributorCode: "0130",
				distributorId:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				Periods:         []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
				Magnitudes:      []measures.Magnitude{measures.AE},
				MeterType:       measures.TLG,
			},
			want: BillingMeasure{
				DistributorCode:        "0130",
				DistributorID:          "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:                   "CUPS",
				PointType:              "",
				RegisterType:           "",
				GenerationDate:         time.Time{},
				EndDate:                time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				InitDate:               time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				Version:                "0",
				Origin:                 "",
				Status:                 Calculating,
				PreviousReadingClosure: measures.DailyReadingClosure{},
				ActualReadingClosure:   measures.DailyReadingClosure{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
				},
				BillingLoadCurve: nil,
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
				},
				Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
				GraphHistory: map[string]*Graph{},
				Magnitudes:   []measures.Magnitude{measures.AE},
				MeterType:    measures.TLG,
			},
		},
		"should return ok for 4 period": {
			input: input{
				cups:            "CUPS",
				initDate:        time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				endDate:         time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				distributorCode: "0130",
				distributorId:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				Periods:         []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
				Magnitudes:      []measures.Magnitude{measures.AE},
				MeterType:       measures.TLG,
			},
			want: BillingMeasure{
				DistributorCode:        "0130",
				DistributorID:          "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:                   "CUPS",
				PointType:              "",
				RegisterType:           "",
				GenerationDate:         time.Time{},
				EndDate:                time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				InitDate:               time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				Version:                "0",
				Origin:                 "",
				Status:                 Calculating,
				PreviousReadingClosure: measures.DailyReadingClosure{},
				ActualReadingClosure:   measures.DailyReadingClosure{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{},
				},
				BillingLoadCurve: nil,
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{},
				},
				Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
				GraphHistory: map[string]*Graph{},
				Magnitudes:   []measures.Magnitude{measures.AE},
				MeterType:    measures.TLG,
			},
		},
		"should return ok for 5 period": {
			input: input{
				cups:            "CUPS",
				initDate:        time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				endDate:         time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				distributorCode: "0130",
				distributorId:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				Periods:         []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5},
				Magnitudes:      []measures.Magnitude{measures.AE},
				MeterType:       measures.TLG,
			},
			want: BillingMeasure{
				DistributorCode:        "0130",
				DistributorID:          "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:                   "CUPS",
				PointType:              "",
				RegisterType:           "",
				GenerationDate:         time.Time{},
				EndDate:                time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				InitDate:               time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				Version:                "0",
				Origin:                 "",
				Status:                 Calculating,
				PreviousReadingClosure: measures.DailyReadingClosure{},
				ActualReadingClosure:   measures.DailyReadingClosure{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{},
					P5: &BillingBalancePeriod{},
				},
				BillingLoadCurve: nil,
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{},
					P5: &measures.Values{},
				},
				Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5},
				GraphHistory: map[string]*Graph{},
				Magnitudes:   []measures.Magnitude{measures.AE},
				MeterType:    measures.TLG,
			},
		},
		"should return ok for 6 period": {
			input: input{
				cups:            "CUPS",
				initDate:        time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				endDate:         time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				distributorCode: "0130",
				distributorId:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				Periods:         []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5, measures.P6},
				Magnitudes:      []measures.Magnitude{measures.AE},
				MeterType:       measures.TLG,
			},
			want: BillingMeasure{
				DistributorCode:        "0130",
				DistributorID:          "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:                   "CUPS",
				PointType:              "",
				RegisterType:           "",
				GenerationDate:         time.Time{},
				EndDate:                time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				InitDate:               time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				Version:                "0",
				Origin:                 "",
				Status:                 Calculating,
				PreviousReadingClosure: measures.DailyReadingClosure{},
				ActualReadingClosure:   measures.DailyReadingClosure{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{},
					P5: &BillingBalancePeriod{},
					P6: &BillingBalancePeriod{},
				},
				BillingLoadCurve: nil,
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{},
					P5: &measures.Values{},
					P6: &measures.Values{},
				},
				Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5, measures.P6},
				GraphHistory: map[string]*Graph{},
				Magnitudes:   []measures.Magnitude{measures.AE},
				MeterType:    measures.TLG,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			testCase.want.GenerateID()
			b := NewBillingMeasure(
				testCase.input.cups,
				testCase.input.initDate,
				testCase.input.endDate,
				testCase.input.distributorCode,
				testCase.input.distributorId,
				testCase.input.Periods,
				testCase.input.Magnitudes,
				testCase.input.MeterType,
			)
			testCase.want.GenerationDate = b.GenerationDate
			assert.Equal(
				t,
				testCase.want,
				b,
				testName)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_SetPreviousReadingClosure(t *testing.T) {

	want := BillingMeasureOk
	want.SetPreviousReadingClosure(PreviousReadingClosure)

	assert.Equal(t, PreviousReadingClosure, want.PreviousReadingClosure)
}

func Test_Unit_Domain_BillingMeasures_ActualReadingClosure(t *testing.T) {

	want := BillingMeasureOk
	want.SetActualReadingClosure(ActualReadingClosure)

	assert.Equal(t, ActualReadingClosure, want.ActualReadingClosure)
}

func Test_Unit_Domain_BillingMeasures_CalcAtrBalance(t *testing.T) {

	type input struct {
		periods     []measures.PeriodKey
		readingType measures.Type
	}
	tests := map[string]struct {
		input input
		want  func() BillingMeasure
	}{
		"should calc atr balance for period 1 well": {
			input: input{
				periods:     []measures.PeriodKey{measures.P1},
				readingType: measures.Absolute,
			},
			want: func() BillingMeasure {
				b := NewBillingMeasure(
					"ES0130000000357054DJ",
					time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
					time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
					"0130",
					"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
					[]measures.PeriodKey{measures.P1},
					[]measures.Magnitude{measures.AE},
					measures.TLG,
				)
				b.ReadingType = measures.Absolute
				b.SetPreviousReadingClosure(PreviousReadingClosure)
				b.SetActualReadingClosure(ActualReadingClosure)

				b.BillingBalance.P0 = BillingBalancePeriodP0
				b.BillingBalance.P1 = BillingBalancePeriodP1
				b.BillingBalance.Origin = ActualReadingClosure.Origin
				b.BillingBalance.EndDate = ActualReadingClosure.EndDate

				return b
			},
		},
		"should calc atr balance for period 2 well": {
			input: input{
				periods:     []measures.PeriodKey{measures.P1, measures.P2},
				readingType: measures.Absolute,
			},
			want: func() BillingMeasure {
				b := NewBillingMeasure(
					"ES0130000000357054DJ",
					time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
					time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
					"0130",
					"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
					[]measures.PeriodKey{measures.P1, measures.P2},
					[]measures.Magnitude{measures.AE},
					measures.TLG,
				)
				b.ReadingType = measures.Absolute
				b.SetPreviousReadingClosure(PreviousReadingClosure)
				b.SetActualReadingClosure(ActualReadingClosure)

				b.BillingBalance.P0 = BillingBalancePeriodP0
				b.BillingBalance.P1 = BillingBalancePeriodP1
				b.BillingBalance.P2 = BillingBalancePeriodP2
				b.BillingBalance.Origin = ActualReadingClosure.Origin
				b.BillingBalance.EndDate = ActualReadingClosure.EndDate

				return b
			},
		},
		"should calc atr balance for period 3 well": {
			input: input{
				periods:     []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
				readingType: measures.Absolute,
			},
			want: func() BillingMeasure {
				b := NewBillingMeasure(
					"ES0130000000357054DJ",
					time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
					time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
					"0130",
					"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
					[]measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					[]measures.Magnitude{measures.AE},
					measures.TLG,
				)
				b.ReadingType = measures.Absolute
				b.SetPreviousReadingClosure(PreviousReadingClosure)
				b.SetActualReadingClosure(ActualReadingClosure)

				b.BillingBalance.P0 = BillingBalancePeriodP0
				b.BillingBalance.P1 = BillingBalancePeriodP1
				b.BillingBalance.P2 = BillingBalancePeriodP2
				b.BillingBalance.P3 = BillingBalancePeriodP3
				b.BillingBalance.Origin = ActualReadingClosure.Origin
				b.BillingBalance.EndDate = ActualReadingClosure.EndDate

				return b
			},
		},
		"should calc atr balance for incremental reading type": {
			input: input{
				periods:     []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
				readingType: measures.Incremental,
			},
			want: func() BillingMeasure {
				b := NewBillingMeasure(
					"ES0130000000357054DJ",
					time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
					time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
					"0130",
					"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
					[]measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					[]measures.Magnitude{measures.AE},
					measures.TLG,
				)
				b.ReadingType = measures.Incremental

				b.SetPreviousReadingClosure(PreviousReadingClosure)
				b.SetActualReadingClosure(ActualReadingClosure)

				b.BillingBalance.P0 = &BillingBalancePeriod{
					AE:                  ActualReadingClosure.CalendarPeriods.P0.AE,
					BalanceValidationAE: measures.Valid,
				}
				b.BillingBalance.P1 = &BillingBalancePeriod{
					AE:                  ActualReadingClosure.CalendarPeriods.P1.AE,
					BalanceValidationAE: measures.Valid,
				}
				b.BillingBalance.P2 = &BillingBalancePeriod{
					AE:                  ActualReadingClosure.CalendarPeriods.P2.AE,
					BalanceValidationAE: measures.Valid,
				}
				b.BillingBalance.P3 = &BillingBalancePeriod{
					AE:                  ActualReadingClosure.CalendarPeriods.P3.AE,
					BalanceValidationAE: measures.Valid,
				}
				b.BillingBalance.Origin = ActualReadingClosure.Origin
				b.BillingBalance.EndDate = ActualReadingClosure.EndDate

				return b
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			b := NewBillingMeasure(
				"ES0130000000357054DJ",
				time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
				time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
				"0130",
				"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				testCase.input.periods,
				[]measures.Magnitude{measures.AE},
				measures.TLG,
			)
			b.ReadingType = testCase.input.readingType
			b.SetPreviousReadingClosure(PreviousReadingClosure)
			b.SetActualReadingClosure(ActualReadingClosure)
			b.CalcAtrBalance()
			want := testCase.want()
			want.GenerationDate = b.GenerationDate
			assert.Equal(t, want, b, testName)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_AreSomeCurveMeasureForPeriod(t *testing.T) {
	tests := map[string]struct {
		input             measures.PeriodKey
		getBillingMeasure func() BillingMeasure
		want              bool
	}{
		"should return false if not curve for period": {
			input: measures.P1,
			want:  false,
			getBillingMeasure: func() BillingMeasure {
				b := NewBillingMeasure(
					"ES0130000000357054DJ",
					time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
					time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
					"0130",
					"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
					[]measures.PeriodKey{measures.P1},
					[]measures.Magnitude{measures.AE},
					measures.TLG,
				)
				b.SetBillingLoadCurve([]BillingLoadCurve{
					{
						EndDate: time.Date(2022, 5, 31, 1, 0, 0, 0, loc).UTC(),
						Origin:  measures.STG,
						AI:      1458,
						AE:      9000,
						R1:      0,
						R2:      0,
						R3:      2416,
						R4:      1916,
						Period:  measures.P2,
					},
				})

				return b
			},
		},
		"should return false if only have curve with Filled origin for period": {
			input: measures.P2,
			want:  false,
			getBillingMeasure: func() BillingMeasure {
				b := NewBillingMeasure(
					"ES0130000000357054DJ",
					time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
					"0130",
					"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
					[]measures.PeriodKey{measures.P1, measures.P2},
					[]measures.Magnitude{measures.AE},
					measures.TLG,
				)
				b.SetBillingLoadCurve([]BillingLoadCurve{
					{
						EndDate: time.Date(2022, 5, 31, 1, 0, 0, 0, time.UTC),
						Origin:  measures.Filled,
						AI:      1458,
						AE:      9000,
						R1:      0,
						R2:      0,
						R3:      2416,
						R4:      1916,
						Period:  measures.P2,
					},
				})

				return b
			},
		},
		"should return true if have curve for period": {
			input: measures.P2,
			want:  true,
			getBillingMeasure: func() BillingMeasure {
				b := NewBillingMeasure(
					"ES0130000000357054DJ",
					time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
					"0130",
					"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
					[]measures.PeriodKey{measures.P1},
					[]measures.Magnitude{measures.AE},
					measures.TLG,
				)
				b.SetBillingLoadCurve([]BillingLoadCurve{
					{
						EndDate: time.Date(2022, 5, 31, 1, 0, 0, 0, time.UTC),
						Origin:  measures.Filled,
						AI:      1458,
						AE:      9000,
						R1:      0,
						R2:      0,
						R3:      2416,
						R4:      1916,
						Period:  measures.P2,
					},
					{
						EndDate: time.Date(2022, 5, 31, 2, 0, 0, 0, loc).UTC(),
						Origin:  measures.STG,
						AI:      1458,
						AE:      9000,
						R1:      0,
						R2:      0,
						R3:      2416,
						R4:      1916,
						Period:  measures.P2,
					},
				})

				return b
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, testCase.want, testCase.getBillingMeasure().AreSomeCurveMeasureForPeriod(testCase.input), testName)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_GetBalancePeriod(t *testing.T) {
	type OriginBillingMeasure struct {
		BillingMeasure *BillingMeasure
		PeriodKey      measures.PeriodKey
	}
	type Want struct {
		BillingBalancePeriod *BillingBalancePeriod
	}
	tests := map[string]struct {
		OriginBillingMeasure OriginBillingMeasure
		Want                 Want
	}{
		"Should Be Okay With Period0": {
			OriginBillingMeasure: OriginBillingMeasure{
				BillingMeasure: &BillingMeasure{
					GenerationDate:         time.Time{},
					EndDate:                time.Time{},
					InitDate:               time.Time{},
					PreviousReadingClosure: measures.DailyReadingClosure{},
					ActualReadingClosure:   measures.DailyReadingClosure{},
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{
							AI: 11.,
							AE: 22.,
							R1: 33.,
							R2: 44.,
							R3: 55.,
							R4: 66.,
						},
						P1: &BillingBalancePeriod{},
						P2: &BillingBalancePeriod{},
						P3: &BillingBalancePeriod{},
						P4: &BillingBalancePeriod{},
						P5: &BillingBalancePeriod{},
						P6: &BillingBalancePeriod{},
					},
					BillingLoadCurve: nil,
					AtrVsCurve:       AtrVsCurve{},
					Periods:          []measures.PeriodKey{},
					GraphHistory:     nil,
				},
				PeriodKey: measures.P0,
			},
			Want: Want{BillingBalancePeriod: &BillingBalancePeriod{
				AI: 11.,
				AE: 22.,
				R1: 33.,
				R2: 44.,
				R3: 55.,
				R4: 66.,
			},
			},
		},
		"Should Be Okay With Period 1": {
			OriginBillingMeasure: OriginBillingMeasure{
				BillingMeasure: &BillingMeasure{
					GenerationDate:         time.Time{},
					EndDate:                time.Time{},
					InitDate:               time.Time{},
					PreviousReadingClosure: measures.DailyReadingClosure{},
					ActualReadingClosure:   measures.DailyReadingClosure{},
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{},
						P1: &BillingBalancePeriod{
							AI: 11.,
							AE: 11.,
							R1: 11.,
							R2: 11.,
							R3: 11.,
							R4: 11.,
						},
						P2: &BillingBalancePeriod{},
						P3: &BillingBalancePeriod{},
						P4: &BillingBalancePeriod{},
						P5: &BillingBalancePeriod{},
						P6: &BillingBalancePeriod{},
					},
					BillingLoadCurve: nil,
					AtrVsCurve:       AtrVsCurve{},
					Periods:          []measures.PeriodKey{},
					GraphHistory:     nil,
				},
				PeriodKey: measures.P1,
			},
			Want: Want{BillingBalancePeriod: &BillingBalancePeriod{
				AI: 11.,
				AE: 11.,
				R1: 11.,
				R2: 11.,
				R3: 11.,
				R4: 11.,
			}},
		},
		"Should return nil": {
			OriginBillingMeasure: OriginBillingMeasure{
				BillingMeasure: &BillingMeasure{
					GenerationDate:         time.Time{},
					EndDate:                time.Time{},
					InitDate:               time.Time{},
					PreviousReadingClosure: measures.DailyReadingClosure{},
					ActualReadingClosure:   measures.DailyReadingClosure{},
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{},
						P1: &BillingBalancePeriod{
							AI: 11.,
							AE: 11.,
							R1: 11.,
							R2: 11.,
							R3: 11.,
							R4: 11.,
						},
						P2: &BillingBalancePeriod{},
						P3: &BillingBalancePeriod{},
						P4: &BillingBalancePeriod{},
						P5: &BillingBalancePeriod{},
						P6: &BillingBalancePeriod{},
					},
					BillingLoadCurve: nil,
					AtrVsCurve:       AtrVsCurve{},
					Periods:          []measures.PeriodKey{},
					GraphHistory:     nil,
				},
				PeriodKey: "Random Key",
			},
			Want: Want{
				BillingBalancePeriod: nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			billingBalanceP := testCase.OriginBillingMeasure.BillingMeasure.GetBalancePeriod(testCase.OriginBillingMeasure.PeriodKey)
			assert.Equal(t, testCase.Want.BillingBalancePeriod, billingBalanceP)
		})

	}
}

func Test_Unit_Domain_BillingMeasures_Calc_Atr_Vs_Curve(t *testing.T) {

	type input struct {
		billingMeasure BillingMeasure
	}

	type want struct {
		newBillingMeasure BillingMeasure
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"should change Balance and AtrVsCurve to P0 if period is 0": {
			input: input{billingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STG,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P0,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{},
			},
			},
			want: want{newBillingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				Periods:         []measures.PeriodKey{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STG,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P0,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
				},
			},
			},
		},
		"should doesn't change AtrVsCurve to P0 if Origin BillingLoadCurve is Filled": {
			input: input{billingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.Filled,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P0,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P0},
			},
			},
			want: want{newBillingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P0},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.Filled,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P0,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
			},
			},
		},
		"should change Balance and AtrVsCurve to P1 if period is 1": {
			input: input{billingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STM,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P1,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1},
			},
			},
			want: want{newBillingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STM,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P1,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
					P1: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1},
			},
			},
		},
		"should change Balance and AtrVsCurve to P2 if period is 2": {
			input: input{billingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STM,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P2,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: -12,
						AE: -8,
						R1: -2,
						R2: -1,
						R3: -2,
						R4: 0,
					},
					P1: &measures.Values{},
					P2: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2},
			},
			},
			want: want{newBillingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STM,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P2,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: -12,
						AE: -8,
						R1: -2,
						R2: -1,
						R3: -2,
						R4: 0,
					},
					P1: &measures.Values{},
					P2: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2},
			},
			},
		},
		"should change Balance and AtrVsCurve to P3 if period is 3": {
			input: input{billingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 20,
						AE: 15,
						R1: 7,
						R2: 19,
						R3: 5,
						R4: 5,
					},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STG,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P3,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 8,
						AE: 7,
						R1: 5,
						R2: 18,
						R3: 3,
						R4: 5,
					},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
			},
			},
			want: want{newBillingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 20,
						AE: 15,
						R1: 7,
						R2: 19,
						R3: 5,
						R4: 5,
					},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STG,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P3,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 8,
						AE: 7,
						R1: 5,
						R2: 18,
						R3: 3,
						R4: 5,
					},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
			},
			},
		},
		"should change Balance and AtrVsCurve to P4 if period is 4": {
			input: input{billingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.TPL,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P4,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
			},
			},
			want: want{newBillingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.TPL,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P4,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
			},
			},
		},
		"should change Balance and AtrVsCurve to P5 if period is 5": {
			input: input{billingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{},
					P5: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STM,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P5,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{},
					P5: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5},
			},
			},
			want: want{newBillingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{},
					P5: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STM,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P5,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{},
					P5: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5},
			},
			},
		},
		"should change Balance and AtrVsCurve to P6 if period is 6": {
			input: input{billingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{},
					P5: &BillingBalancePeriod{},
					P6: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STG,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P6,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: -12,
						AE: -8,
						R1: -2,
						R2: -1,
						R3: -2,
						R4: 0,
					},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{},
					P5: &measures.Values{},
					P6: &measures.Values{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5, measures.P6},
			},
			},
			want: want{newBillingMeasure: BillingMeasure{
				DistributorCode: "0130",
				DistributorID:   "5140cda8-1daa-4f52-b85d-bf29d45ca62a",
				CUPS:            "ES0130000000357054DJ",
				EndDate:         time.Time{},
				InitDate:        time.Time{},
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{},
					P1: &BillingBalancePeriod{},
					P2: &BillingBalancePeriod{},
					P3: &BillingBalancePeriod{},
					P4: &BillingBalancePeriod{},
					P5: &BillingBalancePeriod{},
					P6: &BillingBalancePeriod{
						AI: 14,
						AE: 9,
						R1: 5,
						R2: 16,
						R3: 3,
						R4: 1,
					},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.STG,
						AI:     12,
						AE:     8,
						R1:     2,
						R2:     1,
						R3:     2,
						R4:     0,
						Period: measures.P6,
					},
				},
				AtrVsCurve: AtrVsCurve{
					P0: &measures.Values{
						AI: -12,
						AE: -8,
						R1: -2,
						R2: -1,
						R3: -2,
						R4: 0,
					},
					P1: &measures.Values{},
					P2: &measures.Values{},
					P3: &measures.Values{},
					P4: &measures.Values{},
					P5: &measures.Values{},
					P6: &measures.Values{
						AI: 2,
						AE: 1,
						R1: 3,
						R2: 15,
						R3: 1,
						R4: 1,
					},
				},
				Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5, measures.P6},
			},
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			testCase.input.billingMeasure.CalcAtrVsCurve()
			assert.Equal(t, testCase.want.newBillingMeasure, testCase.input.billingMeasure, testName)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_Is_Empty_Central_Hours_Cch(t *testing.T) {

	tests := map[string]struct {
		b                *BillingMeasure
		expectedResponse bool
	}{
		"Should return false If first or last curve is Origin Filled": {
			b: &BillingMeasure{BillingLoadCurve: []BillingLoadCurve{
				{
					EndDate: time.Date(2022, 06, 5, 01, 00, 00, 0000, loc),
					Origin:  measures.Filled,
				},
				{
					EndDate: time.Date(2022, 06, 5, 02, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 03, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 04, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 05, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 06, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 07, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 8, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 9, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 10, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 11, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 12, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 13, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 14, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 15, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 16, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 17, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 18, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 19, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 20, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 21, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 22, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 23, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 24, 00, 00, 0000, loc),
					Origin:  "",
				},
			}},
			expectedResponse: false,
		},
		"Should return true If first or last curve isn't Origin Filled": {
			b: &BillingMeasure{BillingLoadCurve: []BillingLoadCurve{
				{
					EndDate: time.Date(2022, 06, 5, 01, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 02, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 03, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 04, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 05, 00, 00, 0000, loc),
					Origin:  measures.Filled,
				},
				{
					EndDate: time.Date(2022, 06, 5, 06, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 07, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 8, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 9, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 10, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 11, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 12, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 13, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 14, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 15, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 16, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 17, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 18, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 19, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 20, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 21, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 22, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 23, 00, 00, 0000, loc),
					Origin:  "",
				},
				{
					EndDate: time.Date(2022, 06, 5, 24, 00, 00, 0000, loc),
					Origin:  "",
				},
			}},
			expectedResponse: true,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			actualResponse := testCase.b.IsEmptyCentralHoursCch()
			assert.Equal(t, testCase.expectedResponse, actualResponse)
		})
	}

}
