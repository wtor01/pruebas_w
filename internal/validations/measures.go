package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

type Periods struct {
	measures.Values
}
type QueryClosedCupsMeasureOnDate struct {
	CUPS string
	Date time.Time
}
type QueryCurveCupsMeasureOnDate struct {
	CUPS      string
	StartDate time.Time
	EndDate   time.Time
}
type ProcessedMonthlyClosure struct {
	DistributorCode   string                          `json:"distributor_code" bson:"distributor_code"`
	DistributorID     string                          `json:"distributor_id" bson:"distributor_id"`
	CUPS              string                          `json:"cups" bson:"cups"`
	EndDate           time.Time                       `json:"end_date" bson:"end_date"`
	StartDate         time.Time                       `json:"start_date" bson:"start_date"`
	MeterSerialNumber string                          `json:"meter_serial_number" bson:"meter_serial_number"`
	GenerationDate    time.Time                       `json:"generation_date" bson:"generation_date"`
	ReadingDate       time.Time                       `json:"reading_date" bson:"reading_date"`
	ServiceType       string                          `json:"service_type" bson:"service_type"`
	PointType         string                          `json:"point_type" bson:"point_type"`
	Origin            measures.OriginType             `json:"origin" bson:"origin"`
	MeasurePointType  measures.MeasurePointType       `json:"measure_point_type" bson:"measure_point_type"`
	ContractNumber    string                          `json:"contract_number" bson:"contract_number"`
	CalendarPeriods   ProcessedMonthlyClosureCalendar `json:"calendar_periods" bson:"calendar_periods"`
	QualityCode       measures.QualityCode            `json:"quality_code" bson:"quality_code"`
	ValidationStatus  measures.Status                 `json:"validation_status" bson:"validation_status"`
	TariffId          string                          `json:"tariff_id" bson:"tariff_id"`
	CalendarCode      string                          `json:"calendar_code" bson:"calendar_code"`
	MeterType         string                          `json:"meter_type"  bson:"meter_type"`
	Coefficient       string                          `json:"coefficient" bson:"coefficient"`
	P1Demand          float64                         `json:"p1_demand" bson:"p1_demand"`
	P2Demand          float64                         `json:"p2_demand" bson:"p2_demand"`
	P3Demand          float64                         `json:"p3_demand" bson:"p3_demand"`
	P4Demand          float64                         `json:"p4_demand" bson:"p4_demand"`
	P5Demand          float64                         `json:"p5_demand" bson:"p5_demand"`
	P6Demand          float64                         `json:"p6_demand" bson:"p6_demand"`
	Magnitudes        []measures.Magnitude            `json:"magnitudes" bson:"magnitudes"`
	InvalidationCodes []string                        `json:"invalidation_codes" bson:"invalidation_codes"`
}
type ProcessedMonthlyClosureCalendar struct {
	P0 *ProcessedMonthlyClosurePeriod `json:"P0" bson:"P0"`
	P1 *ProcessedMonthlyClosurePeriod `json:"P1" bson:"P1"`
	P2 *ProcessedMonthlyClosurePeriod `json:"P2" bson:"P2"`
	P3 *ProcessedMonthlyClosurePeriod `json:"P3" bson:"P3"`
	P4 *ProcessedMonthlyClosurePeriod `json:"P4" bson:"P4"`
	P5 *ProcessedMonthlyClosurePeriod `json:"P5" bson:"P5"`
	P6 *ProcessedMonthlyClosurePeriod `json:"P6" bson:"P6"`
}
type ProcessedMonthlyClosurePeriod struct {
	ProcessedDailyClosurePeriod `bson:",inline"`
	AIi                         float64 `json:"AIi" bson:"AIi"`
	AEi                         float64 `json:"AEi" bson:"AEi"`
	R1i                         float64 `json:"R1i" bson:"R1i"`
	R2i                         float64 `json:"R2i" bson:"R2i"`
	R3i                         float64 `json:"R3i" bson:"R3i"`
	R4i                         float64 `json:"R4i" bson:"R4i"`
}
type ProcessedDailyClosurePeriod struct {
	measures.Values          `bson:",inline"`
	InvalidationCodes        []string        `json:"invalidation_codes" bson:"invalidation_codes"`
	ValidationStatus         measures.Status `json:"validation_status" bson:"validation_status"`
	ValidationStatusIsManual bool            `json:"validation_status_is_manual" bson:"validation_status_is_manual"`
	Filled                   bool            `json:"filled" bson:"filled"`
}
type ProcessedLoadCurve struct {
	Id                       string                    `bson:"_id"`
	DistributorCode          string                    `json:"distributor_code" bson:"distributor_code"`
	DistributorID            string                    `json:"distributor_id" bson:"distributor_id"`
	CUPS                     string                    `json:"cups" bson:"cups"`
	Period                   measures.PeriodKey        `json:"period" bson:"period"`
	EndDate                  time.Time                 `json:"end_date" bson:"end_date"`
	MeterSerialNumber        string                    `json:"meter_serial_number" bson:"meter_serial_number"`
	GenerationDate           time.Time                 `json:"generation_date" bson:"generation_date"`
	ReadingDate              time.Time                 `json:"reading_date" bson:"reading_date"`
	Type                     measures.RegisterType     `json:"type" bson:"type"`
	ServiceType              measures.ServiceType      `json:"service_type" bson:"service_type"`
	PointType                measures.PointType        `json:"point_type" bson:"point_type"`
	Origin                   measures.OriginType       `json:"origin" bson:"origin"`
	MeasurePointType         measures.MeasurePointType `json:"measure_point_type" bson:"measure_point_type"`
	PriorityContract         measures.PriorityContract `json:"priority_contract" bson:"priority_contract"`
	MeterConfigId            string                    `json:"meter_config_id" bson:"meter_config_id"`
	AI                       float64                   `json:"AI" bson:"AI"`
	AE                       float64                   `json:"AE" bson:"AE"`
	R1                       float64                   `json:"R1" bson:"R1"`
	R2                       float64                   `json:"R2" bson:"R2"`
	R3                       float64                   `json:"R3" bson:"R3"`
	R4                       float64                   `json:"R4" bson:"R4"`
	InvalidationCodes        []string                  `json:"invalidation_codes" bson:"invalidation_codes"`
	ValidationStatus         measures.Status           `json:"validation_status" bson:"validation_status"`
	ValidationStatusIsManual bool                      `json:"validation_status_is_manual" bson:"validation_status_is_manual"`
	SeasonId                 string                    `json:"season_id" bson:"season_id"`
	DayTypeId                string                    `json:"day_type_id" bson:"day_type_id"`
	P1Demand                 int                       `json:"p1_demand" bson:"p1_demand"`
	P2Demand                 int                       `json:"p2_demand" bson:"p2_demand"`
	P3Demand                 int                       `json:"p3_demand" bson:"p3_demand"`
	P4Demand                 int                       `json:"p4_demand" bson:"p4_demand"`
	P5Demand                 int                       `json:"p5_demand" bson:"p5_demand"`
	P6Demand                 int                       `json:"p6_demand" bson:"p6_demand"`
	Magnitudes               []measures.Magnitude      `json:"magnitudes" bson:"magnitudes"`
}
