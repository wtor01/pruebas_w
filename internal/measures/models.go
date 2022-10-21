package measures

import "time"

type MeterType string

const (
	TLG   MeterType = "TLG"
	TLM   MeterType = "TLM"
	OTHER MeterType = "OTHER"
)

type PriorityContract string

const (
	PriorityContract1 PriorityContract = "1"
	PriorityContract2 PriorityContract = "2"
	PriorityContract3 PriorityContract = "3"
)

type MeterConfigTlgCode string

const (
	TLG_OP_CURVE   MeterConfigTlgCode = "TLG_OP_CURVE"
	NO_TLG         MeterConfigTlgCode = "NO_TLG"
	TLG_OP_NOCURVE MeterConfigTlgCode = "TLG_OP_NOCURVE"
	TLG_NOOP       MeterConfigTlgCode = "TLG_NOOP"
)

type PeriodKey string

const (
	P0 PeriodKey = "P0"
	P1 PeriodKey = "P1"
	P2 PeriodKey = "P2"
	P3 PeriodKey = "P3"
	P4 PeriodKey = "P4"
	P5 PeriodKey = "P5"
	P6 PeriodKey = "P6"
)

type Magnitude string

const (
	AI Magnitude = "AI"
	AE Magnitude = "AE"
	R1 Magnitude = "R1"
	R2 Magnitude = "R2"
	R3 Magnitude = "R3"
	R4 Magnitude = "R4"
)

var ValidMagnitudes = []Magnitude{AI, AE, R1, R2, R3, R4}

var ValidPeriodsCurve = []PeriodKey{P1, P2, P3, P4, P5, P6}

func (p PeriodKey) ToNumber() (int, bool) {
	switch p {
	case P0:
		return 0, true
	case P1:
		return 1, true
	case P2:
		return 2, true
	case P3:
		return 3, true
	case P4:
		return 4, true
	case P5:
		return 5, true
	case P6:
		return 6, true
	default:
		return 0, true
	}
}

type Type string

const (
	Incremental Type = "INC"
	Absolute    Type = "ABS"
)

type ReadingType string

var ReadingTypes = []ReadingType{Curve, BillingClosure, DailyClosure}

const (
	Curve          ReadingType = "curve"
	BillingClosure ReadingType = "billing_closure"
	DailyClosure   ReadingType = "daily_closure"
)

type RegisterType string

const (
	Hourly      RegisterType = "HOURLY"
	QuarterHour RegisterType = "QUARTER"
	NoneType    RegisterType = "NONE"
	Both        RegisterType = "BOTH"
)

type MeasureCurveReadingType string

const (
	HourlyMeasureCurveReadingType  MeasureCurveReadingType = "HOURLY"
	QuarterMeasureCurveReadingType MeasureCurveReadingType = "QUARTER"
)

type PointType string

const (
	PointType1 PointType = "1"
	PointType2 PointType = "2"
	PointType3 PointType = "3"
	PointType4 PointType = "4"
	PointType5 PointType = "5"
)

type MeasurePointType string

const (
	MeasurePointTypeP MeasurePointType = "P"
	MeasurePointTypeR MeasurePointType = "R"
	MeasurePointTypeC MeasurePointType = "C"
	MeasurePointTypeT MeasurePointType = "T"
)

type EquipmentType string

const (
	Main      EquipmentType = "MAIN"
	Redundant EquipmentType = "REDUNDANT"
	Receipt   EquipmentType = "RECEIPT"
)

type OriginType string

const (
	STG                   OriginType = "STG"
	STM                   OriginType = "STM"
	TPL                   OriginType = "TPL"
	File                  OriginType = "File"
	Manual                OriginType = "MANUAL"
	CSV                   OriginType = "CSV"
	Auto                  OriginType = "AUTO"
	Visual                OriginType = "VISUAL"
	Filled                OriginType = "FILLED"
	CalculatedWithQuarter OriginType = "CALCULATED_QUARTER"
)

type Status string

const (
	Valid     Status = "VAL"
	Invalid   Status = "INV"
	Supervise Status = "SUPERV"
	Alert     Status = "ALERT"
	None      Status = "NONE"
)

var StatusValue = map[Status]int{
	Invalid:   1,
	Supervise: 2,
	Alert:     3,
	Valid:     4,
	None:      5,
}

type StatusProcessCurve string

// StatusProcessCurveValue used to have the order of states for process curve
var StatusProcessCurveValue = map[Status]int{
	None:      1,
	Invalid:   2,
	Supervise: 3,
	Alert:     4,
	Valid:     5,
}

func (s Status) Value(statusValue map[Status]int) int {
	v, _ := statusValue[s]

	return v
}

func (s Status) Compare(b map[Status]int, a Status) bool {
	return s.Value(b) < a.Value(b)
}

type QualityCode string

var QualityCodeValue = map[QualityCode]int{
	Absent:     1,
	Incomplete: 2,
	Complete:   3,
}

func (s QualityCode) Value(statusValue map[QualityCode]int) int {
	v, _ := statusValue[s]

	return v
}

func (s QualityCode) Compare(b map[QualityCode]int, a QualityCode) bool {
	return s.Value(b) < a.Value(b)
}

const (
	Complete   QualityCode = "COMPLETE"
	Incomplete QualityCode = "INCOMPLETE"
	Absent     QualityCode = "ABSENT"
)

type ServiceType string

const (
	DcServiceType ServiceType = "D-C"
	GdServiceType ServiceType = "G-D"
	DdServiceType ServiceType = "D-D"
)

type ValidationType string

const (
	Gross   ValidationType = "INM"
	Process ValidationType = "PROC"
	Billing ValidationType = "COHE"
)

type ClosureType string

const (
	Monthly   ClosureType = "MONTHLY"
	Daily     ClosureType = "DAILY"
	Other     ClosureType = "OTHER"
	NoClosure ClosureType = "NO_CLOSURE"
)

type Values struct {
	AI float64 `json:"AI" bson:"AI"`
	AE float64 `json:"AE" bson:"AE"`
	R1 float64 `json:"R1" bson:"R1"`
	R2 float64 `json:"R2" bson:"R2"`
	R3 float64 `json:"R3" bson:"R3"`
	R4 float64 `json:"R4" bson:"R4"`
}

type ValuesMonthly struct {
	Values
	AIi float64 `json:"AIi" bson:"AIi"`
	AEi float64 `json:"AEi" bson:"AEi"`
	R1i float64 `json:"R1i" bson:"R1i"`
	R2i float64 `json:"R2i" bson:"R2i"`
	R3i float64 `json:"R3i" bson:"R3i"`
	R4i float64 `json:"R4i" bson:"R4i"`
}

func (v Values) GetMagnitude(m Magnitude) float64 {
	switch m {
	case AI:
		return v.AI
	case AE:
		return v.AE
	case R1:
		return v.R1
	case R2:
		return v.R2
	case R3:
		return v.R3
	case R4:
		return v.R4
	default:
		return .0
	}
}

func (v *Values) SumMagnitude(m Magnitude, value float64) {
	switch m {
	case AI:
		v.AI += value
	case AE:
		v.AE += value
	case R1:
		v.R1 += value
	case R2:
		v.R2 += value
	case R3:
		v.R3 += value
	case R4:
		v.R4 += value
	}
}

type DailyReadingClosureCalendarPeriod struct {
	Values                   `bson:",inline"`
	Filled                   bool     `json:"filled" bson:"filled"`
	InvalidationCodes        []string `json:"invalidation_codes" bson:"invalidation_codes"`
	ValidationStatus         Status   `json:"validation_status" bson:"validation_status"`
	ValidationStatusIsManual bool     `json:"validation_status_is_manual" bson:"validation_status_is_manual"`
}

type DailyReadingClosureCalendar struct {
	P0 *DailyReadingClosureCalendarPeriod `json:"P0" bson:"P0"`
	P1 *DailyReadingClosureCalendarPeriod `json:"P1" bson:"P1"`
	P2 *DailyReadingClosureCalendarPeriod `json:"P2" bson:"P2"`
	P3 *DailyReadingClosureCalendarPeriod `json:"P3" bson:"P3"`
	P4 *DailyReadingClosureCalendarPeriod `json:"P4" bson:"P4"`
	P5 *DailyReadingClosureCalendarPeriod `json:"P5" bson:"P5"`
	P6 *DailyReadingClosureCalendarPeriod `json:"P6" bson:"P6"`
}

func NewDailyReadingClosureCalendar(periods []PeriodKey) *DailyReadingClosureCalendar {
	calendar := DailyReadingClosureCalendar{}
	for _, period := range periods {
		calendar.FillClosurePeriod(period)
	}

	return &calendar
}

func (closurePeriods DailyReadingClosureCalendar) GetClosurePeriod(p PeriodKey) *DailyReadingClosureCalendarPeriod {

	switch p {
	case P0:
		return closurePeriods.P0
	case P1:
		return closurePeriods.P1
	case P2:
		return closurePeriods.P2
	case P3:
		return closurePeriods.P3
	case P4:
		return closurePeriods.P4
	case P5:
		return closurePeriods.P5
	case P6:
		return closurePeriods.P6
	}

	return nil
}

func (closurePeriods *DailyReadingClosureCalendar) FillClosurePeriod(period PeriodKey) {

	switch period {
	case P0:
		closurePeriods.P0 = &DailyReadingClosureCalendarPeriod{}
	case P1:
		closurePeriods.P1 = &DailyReadingClosureCalendarPeriod{}
	case P2:
		closurePeriods.P2 = &DailyReadingClosureCalendarPeriod{}
	case P3:
		closurePeriods.P3 = &DailyReadingClosureCalendarPeriod{}
	case P4:
		closurePeriods.P4 = &DailyReadingClosureCalendarPeriod{}
	case P5:
		closurePeriods.P5 = &DailyReadingClosureCalendarPeriod{}
	case P6:
		closurePeriods.P6 = &DailyReadingClosureCalendarPeriod{}
	}
}

func (closurePeriods *DailyReadingClosureCalendar) SetMagnitudeToPeriod(period PeriodKey, magnitude Magnitude, value float64) {
	closurePeriod := closurePeriods.GetClosurePeriod(period)
	if closurePeriod == nil {
		return
	}
	switch magnitude {
	case AI:
		closurePeriod.AI = value
	case AE:
		closurePeriod.AE = value
	case R1:
		closurePeriod.R1 = value
	case R2:
		closurePeriod.R2 = value
	case R3:
		closurePeriod.R3 = value
	case R4:
		closurePeriod.R4 = value
	}
}

func (closurePeriods *DailyReadingClosureCalendar) SetFilledToPeriod(period PeriodKey) {
	closurePeriod := closurePeriods.GetClosurePeriod(period)
	if closurePeriod == nil {
		return
	}
	closurePeriod.Filled = true
}

func (closurePeriods *DailyReadingClosureCalendar) GetMagnitudeToPeriod(period PeriodKey, magnitude Magnitude) float64 {
	closurePeriod := closurePeriods.GetClosurePeriod(period)
	if closurePeriod == nil {
		return .0
	}
	return closurePeriod.Values.GetMagnitude(magnitude)
}

type DailyReadingClosure struct {
	ClosureType       ClosureType                 `json:"closure_type" bson:"closure_type"`
	DistributorCode   string                      `json:"distributor_code" bson:"distributor_code"`
	DistributorID     string                      `json:"distributor_id" bson:"distributor_id"`
	CUPS              string                      `json:"cups" bson:"cups"`
	InitDate          time.Time                   `json:"init_date,omitempty" bson:"init_date,omitempty"`
	EndDate           time.Time                   `json:"end_date" bson:"end_date"`
	MeterSerialNumber string                      `json:"meter_serial_number" bson:"meter_serial_number"`
	GenerationDate    time.Time                   `json:"generation_date" bson:"generation_date"`
	ReadingDate       time.Time                   `json:"reading_date" bson:"reading_date"`
	ServiceType       string                      `json:"service_type" bson:"service_type"`
	PointType         string                      `json:"point_type" bson:"point_type"`
	Origin            OriginType                  `json:"origin" bson:"origin"`
	MeasurePointType  MeasurePointType            `json:"measure_point_type" bson:"measure_point_type"`
	ContractNumber    string                      `json:"contract_number" bson:"contract_number"`
	CalendarPeriods   DailyReadingClosureCalendar `json:"calendar_periods" bson:"calendar_periods"`
	QualityCode       QualityCode                 `json:"quality_code" bson:"quality_code"`
	ValidationStatus  Status                      `json:"validation_status" bson:"validation_status"`
	TariffId          string                      `json:"tariff_id" bson:"tariff_id"`
	CalendarCode      string                      `json:"calendar_code" bson:"calendar_code"`
	Coefficient       string                      `json:"coefficient" bson:"coefficient"`
	P1Demand          float64                     `json:"p1_demand" bson:"p1_demand"`
	P2Demand          float64                     `json:"p2_demand" bson:"p2_demand"`
	P3Demand          float64                     `json:"p3_demand" bson:"p3_demand"`
	P4Demand          float64                     `json:"p4_demand" bson:"p4_demand"`
	P5Demand          float64                     `json:"p5_demand" bson:"p5_demand"`
	P6Demand          float64                     `json:"p6_demand" bson:"p6_demand"`
	Magnitudes        []Magnitude                 `json:"magnitudes" bson:"magnitudes"`
}

func (d DailyReadingClosure) GetCalendarPeriod(period PeriodKey) *DailyReadingClosureCalendarPeriod {
	switch period {
	case P0:
		return d.CalendarPeriods.P0
	case P1:
		return d.CalendarPeriods.P1
	case P2:
		return d.CalendarPeriods.P2
	case P3:
		return d.CalendarPeriods.P3
	case P4:
		return d.CalendarPeriods.P4
	case P5:
		return d.CalendarPeriods.P5
	case P6:
		return d.CalendarPeriods.P6
	default:
		return nil
	}
}

func (d DailyReadingClosure) GetCalendarPeriodMagnitude(period PeriodKey, magnitude Magnitude) float64 {
	c := d.GetCalendarPeriod(period)

	if c == nil {
		return .0
	}

	return c.Values.GetMagnitude(magnitude)
}
