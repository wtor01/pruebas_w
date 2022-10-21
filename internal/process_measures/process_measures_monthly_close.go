package process_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"crypto/sha256"
	"fmt"
	"time"
)

type ProcessedMonthlyClosurePeriod struct {
	ProcessedDailyClosurePeriod `bson:",inline"`
	AIi                         float64 `json:"AIi" bson:"AIi"`
	AEi                         float64 `json:"AEi" bson:"AEi"`
	R1i                         float64 `json:"R1i" bson:"R1i"`
	R2i                         float64 `json:"R2i" bson:"R2i"`
	R3i                         float64 `json:"R3i" bson:"R3i"`
	R4i                         float64 `json:"R4i" bson:"R4i"`
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

func (p ProcessedMonthlyClosurePeriod) GetMagnitude(magnitude measures.Magnitude) (float64, float64) {

	switch magnitude {
	case measures.AI:
		return p.AI, p.AIi
	case measures.AE:
		return p.AE, p.AEi
	case measures.R1:
		return p.R1, p.R1i
	case measures.R2:
		return p.R2, p.R2i
	case measures.R3:
		return p.R3, p.R3i
	case measures.R4:
		return p.R4, p.R4i
	}
	return 0, 0
}

func (p ProcessedMonthlyClosureCalendar) GetPeriodValues(periodValue measures.PeriodKey) *ProcessedMonthlyClosurePeriod {

	switch periodValue {
	case measures.P0:
		return p.P0
	case measures.P1:
		return p.P1
	case measures.P2:
		return p.P2
	case measures.P3:
		return p.P3
	case measures.P4:
		return p.P4
	case measures.P5:
		return p.P5
	case measures.P6:
		return p.P6
	}

	return nil
}

type ResumesProcessMonthlyClosure struct {
	Previous *ProcessedMonthlyClosure
	Next     *ProcessedMonthlyClosure
}
type ProcessedMonthlyClosure struct {
	Id                          string                          `bson:"_id"`
	DistributorCode             string                          `json:"distributor_code" bson:"distributor_code"`
	DistributorID               string                          `json:"distributor_id" bson:"distributor_id"`
	CUPS                        string                          `json:"cups" bson:"cups"`
	EndDate                     time.Time                       `json:"end_date" bson:"end_date"`
	StartDate                   time.Time                       `json:"start_date" bson:"start_date"`
	MeterSerialNumber           string                          `json:"meter_serial_number" bson:"meter_serial_number"`
	GenerationDate              time.Time                       `json:"generation_date" bson:"generation_date"`
	ReadingDate                 time.Time                       `json:"reading_date" bson:"reading_date"`
	ServiceType                 measures.ServiceType            `json:"service_type" bson:"service_type"`
	PointType                   measures.PointType              `json:"point_type" bson:"point_type"`
	Origin                      measures.OriginType             `json:"origin" bson:"origin"`
	MeasurePointType            measures.MeasurePointType       `json:"measure_point_type" bson:"measure_point_type"`
	ContractNumber              string                          `json:"contract_number" bson:"contract_number"`
	CalendarPeriods             ProcessedMonthlyClosureCalendar `json:"calendar_periods" bson:"calendar_periods"`
	QualityCode                 measures.QualityCode            `json:"quality_code" bson:"quality_code"`
	ValidationStatus            measures.Status                 `json:"validation_status" bson:"validation_status"`
	TariffId                    string                          `json:"tariff_id" bson:"tariff_id"`
	CalendarCode                string                          `json:"calendar_code" bson:"calendar_code"`
	MeterType                   measures.MeterType              `json:"meter_type"  bson:"meter_type"`
	Coefficient                 string                          `json:"coefficient" bson:"coefficient"`
	P1Demand                    float64                         `json:"p1_demand" bson:"p1_demand"`
	P2Demand                    float64                         `json:"p2_demand" bson:"p2_demand"`
	P3Demand                    float64                         `json:"p3_demand" bson:"p3_demand"`
	P4Demand                    float64                         `json:"p4_demand" bson:"p4_demand"`
	P5Demand                    float64                         `json:"p5_demand" bson:"p5_demand"`
	P6Demand                    float64                         `json:"p6_demand" bson:"p6_demand"`
	Magnitudes                  []measures.Magnitude            `json:"magnitudes" bson:"magnitudes"`
	InvalidationCodes           []string                        `json:"invalidation_codes" bson:"invalidation_codes"`
	Periods                     []measures.PeriodKey            `json:"-" bson:"-"`
	LastProcessedMonthlyClosure *ProcessedMonthlyClosure        `json:"-" bson:"-"`
}

func (p ProcessedMonthlyClosure) getCalendarPeriod(periodValue *ProcessedMonthlyClosurePeriod, Filled bool) *measures.DailyReadingClosureCalendarPeriod {

	return &measures.DailyReadingClosureCalendarPeriod{
		Filled:                   Filled,
		Values:                   periodValue.Values,
		InvalidationCodes:        periodValue.InvalidationCodes,
		ValidationStatus:         periodValue.ValidationStatus,
		ValidationStatusIsManual: periodValue.ValidationStatusIsManual,
	}
}

func (p ProcessedMonthlyClosure) ToDailyReadingClosure() measures.DailyReadingClosure {
	d := measures.DailyReadingClosure{
		ClosureType:       measures.Monthly,
		DistributorCode:   p.DistributorCode,
		DistributorID:     p.DistributorID,
		CUPS:              p.CUPS,
		InitDate:          p.StartDate,
		EndDate:           p.EndDate,
		MeterSerialNumber: p.MeterSerialNumber,
		GenerationDate:    p.GenerationDate,
		ReadingDate:       p.ReadingDate,
		ServiceType:       string(p.ServiceType),
		PointType:         string(p.PointType),
		Origin:            p.Origin,
		MeasurePointType:  p.MeasurePointType,
		ContractNumber:    p.ContractNumber,
		QualityCode:       p.QualityCode,
		ValidationStatus:  p.ValidationStatus,
		TariffId:          p.TariffId,
		CalendarCode:      p.CalendarCode,
		Coefficient:       p.Coefficient,
		P1Demand:          p.P1Demand,
		P2Demand:          p.P2Demand,
		P3Demand:          p.P3Demand,
		P4Demand:          p.P4Demand,
		P5Demand:          p.P5Demand,
		P6Demand:          p.P6Demand,
		Magnitudes:        p.Magnitudes,
	}
	if p.CalendarPeriods.P0 != nil {
		d.CalendarPeriods.P0 = p.getCalendarPeriod(p.CalendarPeriods.P0, p.CalendarPeriods.P0.Filled)
	}
	if p.CalendarPeriods.P1 != nil {
		d.CalendarPeriods.P1 = p.getCalendarPeriod(p.CalendarPeriods.P1, p.CalendarPeriods.P1.Filled)
	}
	if p.CalendarPeriods.P2 != nil {
		d.CalendarPeriods.P2 = p.getCalendarPeriod(p.CalendarPeriods.P2, p.CalendarPeriods.P2.Filled)
	}
	if p.CalendarPeriods.P3 != nil {
		d.CalendarPeriods.P3 = p.getCalendarPeriod(p.CalendarPeriods.P3, p.CalendarPeriods.P3.Filled)
	}
	if p.CalendarPeriods.P4 != nil {
		d.CalendarPeriods.P4 = p.getCalendarPeriod(p.CalendarPeriods.P4, p.CalendarPeriods.P4.Filled)
	}
	if p.CalendarPeriods.P5 != nil {
		d.CalendarPeriods.P5 = p.getCalendarPeriod(p.CalendarPeriods.P5, p.CalendarPeriods.P5.Filled)
	}
	if p.CalendarPeriods.P6 != nil {
		d.CalendarPeriods.P6 = p.getCalendarPeriod(p.CalendarPeriods.P6, p.CalendarPeriods.P6.Filled)
	}

	return d
}

func (p ProcessedMonthlyClosure) ToValidatable() validations.MeasureValidatable {

	m := validations.MeasureValidatable{
		Type:         measures.Incremental,
		EndDate:      p.EndDate,
		ReadingDate:  p.ReadingDate,
		ServiceType:  p.ServiceType,
		PointType:    string(p.PointType),
		CalendarCode: p.CalendarCode,
		PeriodsNumbers: utils.MapSlice(p.Periods, func(item measures.PeriodKey) string {
			return string(item)
		}),
		TechnologyType: string(p.MeterType),
		WhiteListKeys: map[string]struct{}{
			validations.StartDate:   {},
			validations.EndDate:     {},
			validations.MeasureDate: {},
			validations.AI:          {},
			validations.AE:          {},
			validations.R1:          {},
			validations.R2:          {},
			validations.R3:          {},
			validations.R4:          {},
			validations.Qualifier:   {},
		},
	}
	if p.CalendarPeriods.P0 != nil {
		m.P0 = validations.Periods{Values: p.CalendarPeriods.P0.Values}
		m.P0i = validations.Periods{Values: measures.Values{
			AI: p.CalendarPeriods.P0.AIi,
			AE: p.CalendarPeriods.P0.AEi,
			R1: p.CalendarPeriods.P0.R1i,
			R2: p.CalendarPeriods.P0.R2i,
			R3: p.CalendarPeriods.P0.R3i,
			R4: p.CalendarPeriods.P0.R4i,
		}}
	}
	if p.CalendarPeriods.P1 != nil {
		m.P1 = validations.Periods{Values: p.CalendarPeriods.P1.Values}
		m.P1i = validations.Periods{Values: measures.Values{
			AI: p.CalendarPeriods.P1.AIi,
			AE: p.CalendarPeriods.P1.AEi,
			R1: p.CalendarPeriods.P1.R1i,
			R2: p.CalendarPeriods.P1.R2i,
			R3: p.CalendarPeriods.P1.R3i,
			R4: p.CalendarPeriods.P1.R4i,
		}}
	}
	if p.CalendarPeriods.P2 != nil {
		m.P2 = validations.Periods{Values: p.CalendarPeriods.P2.Values}
		m.P2i = validations.Periods{Values: measures.Values{
			AI: p.CalendarPeriods.P2.AIi,
			AE: p.CalendarPeriods.P2.AEi,
			R1: p.CalendarPeriods.P2.R1i,
			R2: p.CalendarPeriods.P2.R2i,
			R3: p.CalendarPeriods.P2.R3i,
			R4: p.CalendarPeriods.P2.R4i,
		}}
	}
	if p.CalendarPeriods.P3 != nil {
		m.P3 = validations.Periods{Values: p.CalendarPeriods.P3.Values}
		m.P3i = validations.Periods{Values: measures.Values{
			AI: p.CalendarPeriods.P3.AIi,
			AE: p.CalendarPeriods.P3.AEi,
			R1: p.CalendarPeriods.P3.R1i,
			R2: p.CalendarPeriods.P3.R2i,
			R3: p.CalendarPeriods.P3.R3i,
			R4: p.CalendarPeriods.P3.R4i,
		}}
	}
	if p.CalendarPeriods.P4 != nil {
		m.P4 = validations.Periods{Values: p.CalendarPeriods.P4.Values}
		m.P4i = validations.Periods{Values: measures.Values{
			AI: p.CalendarPeriods.P4.AIi,
			AE: p.CalendarPeriods.P4.AEi,
			R1: p.CalendarPeriods.P4.R1i,
			R2: p.CalendarPeriods.P4.R2i,
			R3: p.CalendarPeriods.P4.R3i,
			R4: p.CalendarPeriods.P4.R4i,
		}}
	}
	if p.CalendarPeriods.P5 != nil {
		m.P5 = validations.Periods{Values: p.CalendarPeriods.P5.Values}
		m.P5i = validations.Periods{Values: measures.Values{
			AI: p.CalendarPeriods.P5.AIi,
			AE: p.CalendarPeriods.P5.AEi,
			R1: p.CalendarPeriods.P5.R1i,
			R2: p.CalendarPeriods.P5.R2i,
			R3: p.CalendarPeriods.P5.R3i,
			R4: p.CalendarPeriods.P5.R4i,
		}}
	}
	if p.CalendarPeriods.P6 != nil {
		m.P6 = validations.Periods{Values: p.CalendarPeriods.P6.Values}
		m.P6i = validations.Periods{Values: measures.Values{
			AI: p.CalendarPeriods.P6.AIi,
			AE: p.CalendarPeriods.P6.AEi,
			R1: p.CalendarPeriods.P6.R1i,
			R2: p.CalendarPeriods.P6.R2i,
			R3: p.CalendarPeriods.P6.R3i,
			R4: p.CalendarPeriods.P6.R4i,
		}}
	}

	return m
}

func (p *ProcessedMonthlyClosure) SetStatusMeasure(validation validations.ValidatorBase) {
	status := measures.StatusValue
	if !p.ValidationStatus.Compare(status, validation.Action) {
		p.ValidationStatus = validation.Action
	}

	if p.InvalidationCodes == nil {
		p.InvalidationCodes = make([]string, 0, 1)
	}

	for _, i := range p.InvalidationCodes {
		if i == validation.Code {
			return
		}
	}

	p.InvalidationCodes = append(p.InvalidationCodes, validation.Code)
}

func (p *ProcessedMonthlyClosure) GenerateID() {
	id := fmt.Sprintf("%s%s%s%s%s", p.DistributorID, p.EndDate.Format("2006-01-02_15:04"), p.MeterSerialNumber, p.CUPS, p.MeasurePointType)

	p.Id = fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func NewProcessedMonthlyClosure(meterConfig measures.MeterConfig, date time.Time) *ProcessedMonthlyClosure {
	p := ProcessedMonthlyClosure{
		DistributorCode:   meterConfig.DistributorCode,
		DistributorID:     meterConfig.DistributorID,
		EndDate:           date,
		CUPS:              meterConfig.Cups(),
		MeterSerialNumber: meterConfig.SerialNumber(),
		ServiceType:       meterConfig.ServiceType(),
		PointType:         meterConfig.PointType(),
		MeasurePointType:  meterConfig.MeasurePointType(),
		QualityCode:       measures.Complete,
		ValidationStatus:  measures.Valid,
		TariffId:          meterConfig.ContractualSituations.TariffID,
		CalendarCode:      meterConfig.CalendarID,
		MeterType:         meterConfig.MeterType(),
		P1Demand:          meterConfig.ContractualSituations.P1Demand,
		P2Demand:          meterConfig.ContractualSituations.P2Demand,
		P3Demand:          meterConfig.ContractualSituations.P3Demand,
		P4Demand:          meterConfig.ContractualSituations.P4Demand,
		P5Demand:          meterConfig.ContractualSituations.P5Demand,
		P6Demand:          meterConfig.ContractualSituations.P6Demand,
		Magnitudes:        meterConfig.GetMagnitudesActive(),
	}

	p.GenerateID()

	return &p
}

func (p *ProcessedMonthlyClosure) SetLastDailyClose(last ProcessedMonthlyClosure) {
	p.LastProcessedMonthlyClosure = &last
}
