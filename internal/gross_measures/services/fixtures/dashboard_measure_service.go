package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

var MeasureShouldBeClientResult = map[string]int{
	clients.EquipmentTelemedidaType:  6,
	clients.EquipmentTelegestionType: 3,
	clients.EquipmentOtherType:       1,
}

var loc, _ = time.LoadLocation("Europe/Madrid")

var MeasureShouldBeRepoGetDashboardResult = []measures.DashboardMeasureI{
	measures.NewDashboardMeasure(
		5,
		measures.STM,
		measures.Curve,
		measures.Valid,
		time.Date(2022, 04, 02, 0, 0, 0, 0, loc),
	),
	measures.NewDashboardMeasure(
		3,
		measures.STG,
		measures.Curve,
		measures.Valid,
		time.Date(2022, 04, 02, 0, 0, 0, 0, loc),
	),
	measures.NewDashboardMeasure(
		1,
		measures.File,
		measures.BillingClosure,
		measures.Valid,
		time.Date(2022, 04, 02, 0, 0, 0, 0, loc),
	),
}

var MeasureShouldBeWant = measures.DashboardResult{
	Totals: measures.DashboardResultData{
		Telemedida: measures.DashboardResultTelemedidaData{
			Curva: measures.DashboardResultMeasureValue{
				Valid:            5,
				MeasuresShouldBe: 17280,
			},
			Closing: measures.DashboardResultMeasureValue{
				MeasuresShouldBe: 6,
			},
		},
		Telegestion: measures.DashboardResultTelegestionData{
			Curva: measures.DashboardResultMeasureValue{
				Valid:            3,
				MeasuresShouldBe: 8640,
			},
			Closing: measures.DashboardResultMeasureValue{
				MeasuresShouldBe: 3,
			},
			Resumen: measures.DashboardResultMeasureValue{
				MeasuresShouldBe: 90,
			},
		},
		Others: measures.DashboardResultOthersData{
			Closing: measures.DashboardResultMeasureValue{
				Valid:            1,
				MeasuresShouldBe: 1,
			},
		},
	},
	Daily: []measures.DashboardResultDailyData{
		{
			Date: time.Date(2022, 04, 01, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 02, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						Valid:            5,
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						Valid:            3,
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
				Others: measures.DashboardResultOthersData{
					Closing: measures.DashboardResultMeasureValue{
						Valid:            1,
						MeasuresShouldBe: 0,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 03, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 4, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 5, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 6, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 7, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 8, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 9, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 10, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 11, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 12, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 13, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 14, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 15, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 16, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 17, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 18, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 19, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 20, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 21, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 22, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 23, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 24, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 25, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 26, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 27, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 28, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 29, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
		{
			Date: time.Date(2022, 04, 30, 0, 0, 0, 0, loc),
			DashboardResultData: measures.DashboardResultData{
				Telemedida: measures.DashboardResultTelemedidaData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 576,
					},
				},
				Telegestion: measures.DashboardResultTelegestionData{
					Curva: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 288,
					},
					Resumen: measures.DashboardResultMeasureValue{
						MeasuresShouldBe: 3,
					},
				},
			},
		},
	},
}
