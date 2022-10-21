package measures

import (
	"math"
	"sort"
	"time"
)

type DashboardMeasure struct {
	Count       int
	Origin      OriginType
	TypeMeasure ReadingType
	Status      Status
	Date        time.Time
}

func (d DashboardMeasure) SetResult(punteroReuslt *DashboardResultMeasureValue, punteroMap *DashboardResultMeasureValue) {

	if punteroReuslt == nil || punteroMap == nil {
		return
	}

	if d.Status == Valid {
		punteroReuslt.Valid += d.Count
		punteroMap.Valid += d.Count
		return
	}
	if d.Status == Invalid {
		punteroReuslt.Invalid += d.Count
		punteroMap.Invalid += d.Count
		return
	}
	if d.Status == Supervise {
		punteroReuslt.Supervise += d.Count
		punteroMap.Supervise += d.Count
		return
	}
}

type DashboardMeasureI interface {
	SetResult(result *DashboardResult, mapDashboardResultDailyData *map[string]DashboardResultDailyData)
}

func NewDashboardMeasure(
	count int,
	origin OriginType,
	typeMeasure ReadingType,
	status Status,
	date time.Time,
) DashboardMeasureI {
	switch origin {
	case STG:
		return DashboardMeasureTelegestion{
			DashboardMeasure{
				Count:       count,
				Origin:      origin,
				TypeMeasure: typeMeasure,
				Status:      status,
				Date:        date,
			},
		}
	case STM:
		return DashboardMeasureTelemedida{
			DashboardMeasure{
				Count:       count,
				Origin:      origin,
				TypeMeasure: typeMeasure,
				Status:      status,
				Date:        date,
			},
		}
	default:
		return DashboardMeasureOther{
			DashboardMeasure{
				Count:       count,
				Origin:      origin,
				TypeMeasure: typeMeasure,
				Status:      status,
				Date:        date,
			},
		}
	}
}

type DashboardMeasureTelegestion struct {
	DashboardMeasure
}

func (d DashboardMeasureTelegestion) SetResult(result *DashboardResult, datesMeasures *map[string]DashboardResultDailyData) {
	currentDateKey := d.Date.Format("2006-01-02")
	dateMeasure := (*datesMeasures)[currentDateKey]
	var punteroReuslt *DashboardResultMeasureValue
	var punteroMap *DashboardResultMeasureValue

	if d.TypeMeasure == Curve {
		punteroReuslt = &result.Totals.Telegestion.Curva
		punteroMap = &dateMeasure.Telegestion.Curva
	}
	if d.TypeMeasure == DailyClosure {
		punteroReuslt = &result.Totals.Telegestion.Closing
		punteroMap = &dateMeasure.Telegestion.Closing
	}
	if d.TypeMeasure == BillingClosure {
		punteroReuslt = &result.Totals.Telegestion.Resumen
		punteroMap = &dateMeasure.Telegestion.Resumen
	}
	d.DashboardMeasure.SetResult(punteroReuslt, punteroMap)
	(*datesMeasures)[currentDateKey] = dateMeasure
}

type DashboardMeasureTelemedida struct {
	DashboardMeasure
}

func (d DashboardMeasureTelemedida) SetResult(result *DashboardResult, datesMeasures *map[string]DashboardResultDailyData) {
	currentDateKey := d.Date.Format("2006-01-02")
	dateMeasure := (*datesMeasures)[currentDateKey]
	var punteroReuslt *DashboardResultMeasureValue
	var punteroMap *DashboardResultMeasureValue

	if d.TypeMeasure == Curve {
		punteroReuslt = &result.Totals.Telemedida.Curva
		punteroMap = &dateMeasure.Telemedida.Curva
	}
	if d.TypeMeasure == BillingClosure {
		punteroReuslt = &result.Totals.Telemedida.Closing
		punteroMap = &dateMeasure.Telemedida.Closing
	}

	d.DashboardMeasure.SetResult(punteroReuslt, punteroMap)
	(*datesMeasures)[currentDateKey] = dateMeasure
}

type DashboardMeasureOther struct {
	DashboardMeasure
}

func (d DashboardMeasureOther) SetResult(result *DashboardResult, datesMeasures *map[string]DashboardResultDailyData) {
	currentDateKey := d.Date.Format("2006-01-02")
	dateMeasure := (*datesMeasures)[currentDateKey]
	var punteroReuslt *DashboardResultMeasureValue
	var punteroMap *DashboardResultMeasureValue

	if d.TypeMeasure == BillingClosure {
		punteroReuslt = &result.Totals.Others.Closing
		punteroMap = &dateMeasure.Others.Closing
	}

	d.DashboardMeasure.SetResult(punteroReuslt, punteroMap)
	(*datesMeasures)[currentDateKey] = dateMeasure
}

type DashboardResultMeasureValue struct {
	Valid            int `json:"valid"`
	Invalid          int `json:"invalid"`
	Supervise        int `json:"supervise"`
	MeasuresShouldBe int `json:"should_be"`
}

type DashboardResultTelegestionData struct {
	Curva   DashboardResultMeasureValue `json:"curva"`
	Closing DashboardResultMeasureValue `json:"closing"`
	Resumen DashboardResultMeasureValue `json:"resumen"`
}

type DashboardResultTelemedidaData struct {
	Curva   DashboardResultMeasureValue `json:"curva"`
	Closing DashboardResultMeasureValue `json:"closing"`
}

type DashboardResultOthersData struct {
	Closing DashboardResultMeasureValue `json:"closing"`
}

type DashboardResultData struct {
	Telegestion DashboardResultTelegestionData `json:"telegestion"`
	Telemedida  DashboardResultTelemedidaData  `json:"telemedida"`
	Others      DashboardResultOthersData      `json:"others"`
}

type DashboardResultDailyData struct {
	Date time.Time `json:"Date"`
	DashboardResultData
}

type DashboardResult struct {
	Totals DashboardResultData        `json:"totals"`
	Daily  []DashboardResultDailyData `json:"daily"`
}

type MeasureShouldBeValue struct {
	Daily int
	Total int
}

type MeasureShouldBe struct {
	Telegestion struct {
		Curva   MeasureShouldBeValue
		Closing MeasureShouldBeValue
		Resumen MeasureShouldBeValue
	}
	Telemedida struct {
		Curva   MeasureShouldBeValue
		Closing MeasureShouldBeValue
	}
	Others struct {
		Closing MeasureShouldBeValue
	}
}

type MeasureCount struct {
	Count int
	Total int
}

func GetDashboardResultDailyData(startDate time.Time, endDate time.Time, measureShouldBeInfo MeasureShouldBe) map[string]DashboardResultDailyData {
	datesMeasures := make(map[string]DashboardResultDailyData)
	diff := endDate.Sub(startDate)

	days := math.Ceil(diff.Hours() / 24)

	for i := 0; i < int(days); i++ {
		date := startDate.Add(time.Hour * 24 * time.Duration(i))
		datesMeasures[date.Format("2006-01-02")] = DashboardResultDailyData{
			Date: date,
			DashboardResultData: DashboardResultData{
				Telegestion: DashboardResultTelegestionData{
					Curva: DashboardResultMeasureValue{
						MeasuresShouldBe: measureShouldBeInfo.Telegestion.Curva.Daily,
					},
					Closing: DashboardResultMeasureValue{
						MeasuresShouldBe: measureShouldBeInfo.Telegestion.Closing.Daily,
					},
					Resumen: DashboardResultMeasureValue{
						MeasuresShouldBe: measureShouldBeInfo.Telegestion.Resumen.Daily,
					},
				},
				Telemedida: DashboardResultTelemedidaData{
					Curva: DashboardResultMeasureValue{
						MeasuresShouldBe: measureShouldBeInfo.Telemedida.Curva.Daily,
					},
					Closing: DashboardResultMeasureValue{
						MeasuresShouldBe: measureShouldBeInfo.Telemedida.Closing.Daily,
					},
				},
				Others: DashboardResultOthersData{
					Closing: DashboardResultMeasureValue{
						MeasuresShouldBe: measureShouldBeInfo.Others.Closing.Daily,
					},
				},
			},
		}
	}
	return datesMeasures
}

func GetDashboardResult(dashboardMeasures []DashboardMeasureI, measureShouldBeInfo MeasureShouldBe, datesMeasures map[string]DashboardResultDailyData) DashboardResult {
	result := DashboardResult{
		Totals: DashboardResultData{
			Telegestion: DashboardResultTelegestionData{
				Curva: DashboardResultMeasureValue{
					MeasuresShouldBe: measureShouldBeInfo.Telegestion.Curva.Total,
				},
				Closing: DashboardResultMeasureValue{
					MeasuresShouldBe: measureShouldBeInfo.Telegestion.Closing.Total,
				},
				Resumen: DashboardResultMeasureValue{
					MeasuresShouldBe: measureShouldBeInfo.Telegestion.Resumen.Total,
				},
			},
			Telemedida: DashboardResultTelemedidaData{
				Curva: DashboardResultMeasureValue{
					MeasuresShouldBe: measureShouldBeInfo.Telemedida.Curva.Total,
				},
				Closing: DashboardResultMeasureValue{
					MeasuresShouldBe: measureShouldBeInfo.Telemedida.Closing.Total,
				},
			},
			Others: DashboardResultOthersData{
				Closing: DashboardResultMeasureValue{
					MeasuresShouldBe: measureShouldBeInfo.Others.Closing.Total,
				},
			},
		},
		Daily: nil,
	}
	a := &result
	b := &datesMeasures

	for _, d := range dashboardMeasures {
		d.SetResult(a, b)
	}

	for _, data := range datesMeasures {
		result.Daily = append(result.Daily, data)
	}

	sort.Slice(result.Daily, func(i, j int) bool {
		return result.Daily[i].Date.Before(result.Daily[j].Date)
	})

	return result
}

func CalcByCurveType(curveType RegisterType, metersCount int) int {
	switch curveType {
	case Hourly:
		return 24 * metersCount
	case QuarterHour:
		return (24 * 4) * metersCount
	case Both:
		return 24 * (24 * 4) * metersCount
	default:
		return 0
	}
}

type DashboardCupsValues struct {
	Valid     int `json:"valid"`
	Invalid   int `json:"invalid"`
	Supervise int `json:"supervise"`
	None      int `json:"none"`
	Total     int `json:"total"`
	ShouldBe  int `json:"should_be"`
}

func (d *DashboardCupsValues) setStatus(status Status, value int) {
	switch status {
	case Valid:
		d.Valid = value
	case Invalid:
		d.Invalid = value
	case Supervise:
		d.Supervise = value
	default:
		d.None = value
	}
}
func (d *DashboardCupsValues) setTotal(value int) {
	d.Total = value
}
func (d *DashboardCupsValues) setShouldBe(value int) {
	d.ShouldBe = value
}

type DashboardCupsReading struct {
	Curve   DashboardCupsValues
	Daily   DashboardCupsValues
	Monthly DashboardCupsValues
}

func (d *DashboardCupsReading) SetReading(readingType ReadingType, status Status, value int) {
	switch readingType {
	case Curve:
		d.Curve.setStatus(status, value)
	case DailyClosure:
		d.Daily.setStatus(status, value)
	case BillingClosure:
		d.Monthly.setStatus(status, value)
	}
}

func (d *DashboardCupsReading) SetTotal(readingType ReadingType, value int) {
	switch readingType {
	case Curve:
		d.Curve.setTotal(value)
	case DailyClosure:
		d.Daily.setTotal(value)
	case BillingClosure:
		d.Monthly.setTotal(value)
	}
}

func (d *DashboardCupsReading) SetShouldBe(readingType ReadingType, value int) {
	switch readingType {
	case Curve:
		d.Curve.setShouldBe(value)
	case DailyClosure:
		d.Daily.setShouldBe(value)
	case BillingClosure:
		d.Monthly.setShouldBe(value)
	}
}

type DashboardCups struct {
	Cups   string
	Values DashboardCupsReading
}

type DashboardListCups struct {
	Cups  []DashboardCups
	Total int
}
