package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

var Actual_Closure_Complete = measures.DailyReadingClosure{
	DistributorCode:  "0130",
	DistributorID:    "DistributorX",
	CUPS:             "ESXXXXXXXXXX",
	ServiceType:      "D-C",
	PointType:        "4",
	Origin:           measures.STM,
	MeasurePointType: measures.MeasurePointTypeR,
	ContractNumber:   "1",
	CalendarPeriods: measures.DailyReadingClosureCalendar{
		P0: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 28000,
				AE: 28000,
				R1: 28000,
			},
		},
		P1: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 10000,
				AE: 10000,
				R1: 10000,
			},
		},
		P2: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 12000,
				AE: 12000,
				R1: 12000,
			},
		},
		P3: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 6000,
				AE: 6000,
				R1: 6000,
			},
		},
	},
}

var Actual_Closure_Invalid_Balance = measures.DailyReadingClosure{
	DistributorCode:  "0130",
	DistributorID:    "DistributorX",
	CUPS:             "ESXXXXXXXXXX",
	ServiceType:      "D-C",
	PointType:        "4",
	Origin:           measures.STM,
	MeasurePointType: measures.MeasurePointTypeR,
	ContractNumber:   "1",
	CalendarPeriods: measures.DailyReadingClosureCalendar{
		P0: &measures.DailyReadingClosureCalendarPeriod{
			Filled: true,
		},
		P1: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 10000,
				AE: 10000,
				R1: 10000,
			},
		},
		P2: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 12000,
				AE: 12000,
				R1: 12000,
			},
		},
		P3: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 6000,
				AE: 6000,
				R1: 6000,
			},
		},
	},
}

var Actual_Closure_Close_ATR_Invalid = measures.DailyReadingClosure{
	DistributorCode:  "0130",
	DistributorID:    "DistributorX",
	CUPS:             "ESXXXXXXXXXX",
	ServiceType:      "D-C",
	PointType:        "4",
	Origin:           measures.STM,
	MeasurePointType: measures.MeasurePointTypeR,
	ContractNumber:   "1",
	CalendarPeriods: measures.DailyReadingClosureCalendar{
		P0: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 28000,
				AE: 28000,
				R1: 28000,
			},
		},
		P1: &measures.DailyReadingClosureCalendarPeriod{
			Filled: true,
		},
		P2: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 12000,
				AE: 12000,
				R1: 12000,
			},
		},
		P3: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 6000,
				AE: 6000,
				R1: 6000,
			},
		},
	},
}

var Actual_Closure_Close_ATR_Invalid_2 = measures.DailyReadingClosure{
	DistributorCode:  "0130",
	DistributorID:    "DistributorX",
	CUPS:             "ESXXXXXXXXXX",
	ServiceType:      "D-C",
	PointType:        "4",
	Origin:           measures.STM,
	MeasurePointType: measures.MeasurePointTypeR,
	ContractNumber:   "1",
	CalendarPeriods: measures.DailyReadingClosureCalendar{
		P0: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 28000,
				AE: 28000,
				R1: 28000,
			},
		},
		P1: &measures.DailyReadingClosureCalendarPeriod{
			Filled: true,
		},
		P2: &measures.DailyReadingClosureCalendarPeriod{
			Filled: true,
		},
		P3: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 6000,
				AE: 6000,
				R1: 6000,
			},
		},
	},
}

var Actual_Closure_Close_ATR_Balance_Empty = measures.DailyReadingClosure{
	DistributorCode:  "0130",
	DistributorID:    "DistributorX",
	CUPS:             "ESXXXXXXXXXX",
	ServiceType:      "D-C",
	PointType:        "4",
	Origin:           measures.STM,
	MeasurePointType: measures.MeasurePointTypeR,
	ContractNumber:   "1",
	CalendarPeriods: measures.DailyReadingClosureCalendar{
		P0: &measures.DailyReadingClosureCalendarPeriod{
			Filled: true,
		},
		P1: &measures.DailyReadingClosureCalendarPeriod{
			Filled: true,
		},
		P2: &measures.DailyReadingClosureCalendarPeriod{
			Filled: true,
		},
		P3: &measures.DailyReadingClosureCalendarPeriod{
			Filled: true,
		},
	},
}

var Previous_Closure_Complete = measures.DailyReadingClosure{
	DistributorCode:  "0130",
	DistributorID:    "DistributorX",
	CUPS:             "ESXXXXXXXXXX",
	ServiceType:      "D-C",
	PointType:        "4",
	Origin:           measures.STM,
	MeasurePointType: measures.MeasurePointTypeR,
	ContractNumber:   "1",
	CalendarPeriods: measures.DailyReadingClosureCalendar{
		P0: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 19000,
				AE: 19000,
				R1: 19000,
			},
		},
		P1: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 7000,
				AE: 7000,
				R1: 7000,
			},
		},
		P2: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 8000,
				AE: 8000,
				R1: 8000,
			},
		},
		P3: &measures.DailyReadingClosureCalendarPeriod{
			Values: measures.Values{
				AI: 4000,
				AE: 4000,
				R1: 4000,
			},
		},
	},
}
