package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"time"
)

func createFloat(x float64) *float64 {
	return &x
}

var BillingSelfConsumption_fixtures_C_B = billing_measures.BillingSelfConsumption{Points: []billing_measures.BillingSelfConsumptionPoint{
	billing_measures.BillingSelfConsumptionPoint{
		ServicePointType: "Consumo", CUPS: "1234", InitDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, billing_measures.BillingSelfConsumptionPoint{
		ServicePointType: "G-D", CUPS: "4321", MvhReceived: []billing_measures.BillingSelfConsumptionPointMvhReceived(nil),
		InitDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}},
	Curve: []billing_measures.BillingSelfConsumptionCurve{billing_measures.BillingSelfConsumptionCurve{
		EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{
			{
				CUPS: "1234",
				AI:   10,
				AE:   20,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHGN: createFloat(6),
					EHCR: createFloat(4),
					EHEX: createFloat(0),
					EHAU: createFloat(6),
					EHSA: nil,
					EHDC: createFloat(10),
				},
			},
			{
				CUPS: "4321",
				AI:   30,
				AE:   40,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHGN: createFloat(10),
					EHCR: nil,
					EHEX: createFloat(0),
					EHAU: nil,
					EHSA: createFloat(0),
					EHDC: nil,
				},
			},
		},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHGN: createFloat(10),
			EHCR: nil,
			EHEX: createFloat(4),
			EHAU: nil,
			EHSA: createFloat(0),
			EHDC: nil,
		},
	}}}

var BillingSelfConsumption_fixtures_D_E1 = billing_measures.BillingSelfConsumption{
	Points: []billing_measures.BillingSelfConsumptionPoint{
		billing_measures.BillingSelfConsumptionPoint{
			ServicePointType: "Consumo", CUPS: "1234", InitDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, billing_measures.BillingSelfConsumptionPoint{
			ServicePointType: "G-D", CUPS: "4321", MvhReceived: []billing_measures.BillingSelfConsumptionPointMvhReceived(nil),
			InitDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}},
	Curve: []billing_measures.BillingSelfConsumptionCurve{billing_measures.BillingSelfConsumptionCurve{
		EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{
			{
				CUPS: "1234",
				AI:   10,
				AE:   20,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHGN: createFloat(6),
					EHCR: nil,
					EHEX: nil,
					EHAU: nil,
					EHSA: nil,
					EHDC: nil,
				},
			},
			{
				CUPS: "4321",
				AI:   30,
				AE:   40,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHGN: createFloat(40),
					EHCR: nil,
					EHEX: createFloat(0),
					EHAU: nil,
					EHSA: createFloat(30),
					EHDC: nil,
				},
			},
		},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHGN: createFloat(40),
			EHCR: nil,
			EHEX: nil,
			EHAU: nil,
			EHSA: createFloat(30),
			EHDC: nil,
		},
	}},
}

var BillingSelfConsumption_fixtures_E1 = billing_measures.BillingSelfConsumption{Points: []billing_measures.BillingSelfConsumptionPoint{
	billing_measures.BillingSelfConsumptionPoint{
		ServicePointType: "Consumo", CUPS: "1234", InitDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, billing_measures.BillingSelfConsumptionPoint{
		ServicePointType: "G-D", CUPS: "4321", MvhReceived: []billing_measures.BillingSelfConsumptionPointMvhReceived(nil),
		InitDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}},
	Curve: []billing_measures.BillingSelfConsumptionCurve{billing_measures.BillingSelfConsumptionCurve{
		EndDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{
			{
				CUPS: "1234",
				AI:   10,
				AE:   20,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHGN: nil,
					EHCR: createFloat(-10),
					EHEX: nil,
					EHAU: nil,
					EHSA: nil,
					EHDC: nil,
				},
			},
			{
				CUPS: "4321",
				AI:   30,
				AE:   40,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHGN: createFloat(40),
					EHCR: nil,
					EHEX: createFloat(0),
					EHAU: createFloat(0),
					EHSA: createFloat(0),
					EHDC: nil,
				},
			},
		},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHGN: createFloat(40),
			EHCR: nil,
			EHEX: nil,
			EHAU: nil,
			EHSA: createFloat(0),
			EHDC: nil,
		},
	}}}

var BillingSelfConsumption_fixtures_A_real_example_1_input = billing_measures.BillingSelfConsumption{
	Points: []billing_measures.BillingSelfConsumptionPoint{{CUPS: "1234", ServicePointType: billing_measures.ConsumoServicePointType}, {CUPS: "4321", ServicePointType: billing_measures.GdServicePointType}},
	Curve: []billing_measures.BillingSelfConsumptionCurve{{
		EndDate: time.Date(2022, time.July, 12, 1, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   4000,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 2, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   2000,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 3, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   6000,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 4, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   7000,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 5, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   9000,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 6, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   9000,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 7, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   10000,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 8, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   26000,
			AE:   1000,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 9, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   31000,
			AE:   6000,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 10, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   15000,
			AE:   10000,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 11, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   22000,
			AE:   29000,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 12, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   13000,
			AE:   34000,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 13, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   11000,
			AE:   31000,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 14, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   9000,
			AE:   22000,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 15, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   15000,
			AE:   14000,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 16, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   17000,
			AE:   10000,
		}},
	},
	},
}

var BillingSelfConsumption_fixtures_A_real_example_1_output = billing_measures.BillingSelfConsumption{
	Points: []billing_measures.BillingSelfConsumptionPoint{{CUPS: "1234", ServicePointType: billing_measures.ConsumoServicePointType}, {CUPS: "4321", ServicePointType: billing_measures.GdServicePointType}},
	Curve: []billing_measures.BillingSelfConsumptionCurve{{
		EndDate: time.Date(2022, time.July, 12, 1, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   4000,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(4000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(4000),
			EHEX: createFloat(0)},
	}, {
		EndDate: time.Date(2022, time.July, 12, 2, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   2000,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(2000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(2000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 3, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   6000,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(6000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(6000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 4, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   7000,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(7000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(7000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 5, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   9000,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(9000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(9000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 6, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   9000,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(9000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(9000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 7, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   10000,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(10000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(10000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 8, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   26000,
			AE:   1000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(25000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(25000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 9, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   31000,
			AE:   6000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(25000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(25000),
			EHEX: createFloat(0),
		},
	}, {

		EndDate: time.Date(2022, time.July, 12, 10, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   15000,
			AE:   10000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(5000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(5000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 11, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   22000,
			AE:   29000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHEX: createFloat(7000),
				EHCR: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHEX: createFloat(7000),
			EHCR: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 12, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   13000,
			AE:   34000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHEX: createFloat(21000),
				EHCR: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHEX: createFloat(21000),
			EHCR: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 13, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   11000,
			AE:   31000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHEX: createFloat(20000),
				EHCR: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHEX: createFloat(20000),
			EHCR: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 14, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   9000,
			AE:   22000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHEX: createFloat(13000),
				EHCR: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHEX: createFloat(13000),
			EHCR: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 15, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   15000,
			AE:   14000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(1000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(1000),
			EHEX: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 16, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1234",
			AI:   17000,
			AE:   10000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHCR: createFloat(7000),
				EHEX: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHCR: createFloat(7000),
			EHEX: createFloat(0),
		},
	},
	},
}

var BillingSelfConsumption_fixtures_C_B_real_example_1_input = billing_measures.BillingSelfConsumption{
	Points: []billing_measures.BillingSelfConsumptionPoint{
		{CUPS: "1", ServicePointType: billing_measures.ConsumoServicePointType},
		{CUPS: "2", ServicePointType: billing_measures.ConsumoServicePointType},
		{CUPS: "3", ServicePointType: billing_measures.ConsumoServicePointType},
		{CUPS: "4", ServicePointType: billing_measures.GdServicePointType},
		{CUPS: "5", ServicePointType: billing_measures.GdServicePointType}},
	Curve: []billing_measures.BillingSelfConsumptionCurve{{
		EndDate: time.Date(2022, time.July, 12, 1, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1",
			AI:   4000,
		}, {
			CUPS: "2",
			AI:   1000,
		}, {
			CUPS: "3",
			AI:   6000,
		}, {
			CUPS: "4",
			AI:   500,
			AE:   0,
		}, {
			CUPS: "5",
			AI:   200,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 2, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1",
			AI:   2000,
		}, {
			CUPS: "2",
			AI:   1000,
		}, {
			CUPS: "3",
			AI:   3000,
		}, {
			CUPS: "4",
			AI:   500,
			AE:   0,
		}, {
			CUPS: "5",
			AI:   200,
			AE:   0,
		}},
	}, {
		EndDate: time.Date(2022, time.July, 12, 3, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1",
			AI:   6000,
		}, {
			CUPS: "2",
			AI:   1000,
		}, {
			CUPS: "3",
			AI:   4000,
		}, {
			CUPS: "4",
			AI:   500,
			AE:   0,
		}, {
			CUPS: "5",
			AI:   200,
			AE:   0,
		}},
	},
		{
			EndDate: time.Date(2022, time.July, 12, 3, 0, 0, 0, time.UTC),
			Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
				CUPS: "1",
				AI:   7000,
			}, {
				CUPS: "2",
				AI:   1000,
			}, {
				CUPS: "3",
				AI:   4000,
			}, {
				CUPS: "4",
				AI:   500,
				AE:   0,
			}, {
				CUPS: "5",
				AI:   200,
				AE:   0,
			}},
		},
	},
}

var BillingSelfConsumption_fixtures_C_B_real_example_1_output = billing_measures.BillingSelfConsumption{
	Points: []billing_measures.BillingSelfConsumptionPoint{
		{CUPS: "1", ServicePointType: billing_measures.ConsumoServicePointType},
		{CUPS: "2", ServicePointType: billing_measures.ConsumoServicePointType},
		{CUPS: "3", ServicePointType: billing_measures.ConsumoServicePointType},
		{CUPS: "4", ServicePointType: billing_measures.GdServicePointType},
		{CUPS: "5", ServicePointType: billing_measures.GdServicePointType}},
	Curve: []billing_measures.BillingSelfConsumptionCurve{{
		EndDate: time.Date(2022, time.July, 12, 1, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1",
			AI:   4000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHDC: createFloat(4000),
				EHCR: createFloat(4000),

				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		},
			{
				CUPS: "2",
				AI:   1000,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHDC: createFloat(1000),
					EHCR: createFloat(1000),

					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}, {
				CUPS: "3",
				AI:   6000,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHDC: createFloat(6000),
					EHCR: createFloat(6000),

					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}, {
				CUPS: "4",
				AI:   500,
				AE:   0,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{EHSA: createFloat(500),

					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}, {
				CUPS: "5",
				AI:   200,
				AE:   0,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHSA: createFloat(200),
					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHDC: createFloat(11000),
			EHCR: createFloat(6000),
			EHSA: createFloat(700),
			EHAU: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 2, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1",
			AI:   2000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHDC: createFloat(2000),
				EHCR: createFloat(2000),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}, {
			CUPS: "2",
			AI:   1000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHDC: createFloat(1000),
				EHCR: createFloat(1000),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}, {
			CUPS: "3",
			AI:   3000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHDC: createFloat(3000),
				EHCR: createFloat(3000),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}, {
			CUPS: "4",
			AI:   500,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHSA: createFloat(500),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}, {
			CUPS: "5",
			AI:   200,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHSA: createFloat(200),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHSA: createFloat(700),
			EHDC: createFloat(6000),
			EHAU: createFloat(0),
		},
	}, {
		EndDate: time.Date(2022, time.July, 12, 3, 0, 0, 0, time.UTC),
		Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
			CUPS: "1",
			AI:   6000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHDC: createFloat(6000),
				EHCR: createFloat(6000),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}, {
			CUPS: "2",
			AI:   1000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHDC: createFloat(1000),
				EHCR: createFloat(1000),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}, {
			CUPS: "3",
			AI:   4000,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHDC: createFloat(4000),
				EHCR: createFloat(4000),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}, {
			CUPS: "4",
			AI:   500,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHSA: createFloat(500),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}, {
			CUPS: "5",
			AI:   200,
			AE:   0,
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHSA: createFloat(200),
				EHGN: createFloat(0),
				EHEX: createFloat(0),
				EHAU: createFloat(0),
			},
		}},
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHSA: createFloat(700),
			EHDC: createFloat(11000),
			EHCR: createFloat(11000),
			EHAU: createFloat(0),
		},
	},
		{
			EndDate: time.Date(2022, time.July, 12, 3, 0, 0, 0, time.UTC),
			Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
				CUPS: "1",
				AI:   7000,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHDC: createFloat(9000),
					EHCR: createFloat(7000),
					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}, {
				CUPS: "2",
				AI:   1000,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHDC: createFloat(3000),
					EHCR: createFloat(1000),
					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}, {
				CUPS: "3",
				AI:   4000,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHDC: createFloat(7000),
					EHCR: createFloat(4000),
					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}, {
				CUPS: "4",
				AI:   500,
				AE:   0,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHSA: createFloat(500),
					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}, {
				CUPS: "5",
				AI:   200,
				AE:   0,
				BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
					EHSA: createFloat(200),
					EHGN: createFloat(0),
					EHEX: createFloat(0),
					EHAU: createFloat(0),
				},
			}},
			BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
				EHSA: createFloat(700),
				EHDC: createFloat(12000),
				EHCR: createFloat(12000),
				EHAU: createFloat(0),
			},
		},
	},
}
