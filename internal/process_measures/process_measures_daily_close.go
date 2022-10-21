package process_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"crypto/sha256"
	"fmt"
	"time"
)

type ProcessedDailyClosurePeriod struct {
	measures.Values          `bson:",inline"`
	InvalidationCodes        []string        `json:"invalidation_codes" bson:"invalidation_codes"`
	ValidationStatus         measures.Status `json:"validation_status" bson:"validation_status"`
	ValidationStatusIsManual bool            `json:"validation_status_is_manual" bson:"validation_status_is_manual"`
	Filled                   bool            `json:"filled" bson:"filled"`
}

type ProcessedDailyClosureCalendar struct {
	P0 *ProcessedDailyClosurePeriod `json:"P0" bson:"P0"`
	P1 *ProcessedDailyClosurePeriod `json:"P1" bson:"P1"`
	P2 *ProcessedDailyClosurePeriod `json:"P2" bson:"P2"`
	P3 *ProcessedDailyClosurePeriod `json:"P3" bson:"P3"`
	P4 *ProcessedDailyClosurePeriod `json:"P4" bson:"P4"`
	P5 *ProcessedDailyClosurePeriod `json:"P5" bson:"P5"`
	P6 *ProcessedDailyClosurePeriod `json:"P6" bson:"P6"`
}

type ProcessedDailyClosure struct {
	Id                string                        `bson:"_id"`
	DistributorCode   string                        `json:"distributor_code" bson:"distributor_code"`
	DistributorID     string                        `json:"distributor_id" bson:"distributor_id"`
	CUPS              string                        `json:"cups" bson:"cups"`
	EndDate           time.Time                     `json:"end_date" bson:"end_date"`
	MeterSerialNumber string                        `json:"meter_serial_number" bson:"meter_serial_number"`
	GenerationDate    time.Time                     `json:"generation_date" bson:"generation_date"`
	ReadingDate       time.Time                     `json:"reading_date" bson:"reading_date"`
	ServiceType       measures.ServiceType          `json:"service_type" bson:"service_type"`
	PointType         measures.PointType            `json:"point_type" bson:"point_type"`
	Origin            measures.OriginType           `json:"origin" bson:"origin"`
	MeasurePointType  measures.MeasurePointType     `json:"measure_point_type" bson:"measure_point_type"`
	ContractNumber    string                        `json:"contract_number" bson:"contract_number"`
	CalendarPeriods   ProcessedDailyClosureCalendar `json:"calendar_periods" bson:"calendar_periods"`
	QualityCode       measures.QualityCode          `json:"quality_code" bson:"quality_code"`
	ValidationStatus  measures.Status               `json:"validation_status" bson:"validation_status"`
	TariffId          string                        `json:"tariff_id" bson:"tariff_id"`
	CalendarCode      string                        `json:"calendar_code" bson:"calendar_code"`
	MeterType         measures.MeterType            `json:"meter_type"  bson:"meter_type"`
	Coefficient       string                        `json:"coefficient" bson:"coefficient"`
	P1Demand          float64                       `json:"p1_demand" bson:"p1_demand"`
	P2Demand          float64                       `json:"p2_demand" bson:"p2_demand"`
	P3Demand          float64                       `json:"p3_demand" bson:"p3_demand"`
	P4Demand          float64                       `json:"p4_demand" bson:"p4_demand"`
	P5Demand          float64                       `json:"p5_demand" bson:"p5_demand"`
	P6Demand          float64                       `json:"p6_demand" bson:"p6_demand"`
	Magnitudes        []measures.Magnitude          `json:"magnitudes" bson:"magnitudes"`
	InvalidationCodes []string                      `json:"invalidation_codes" bson:"invalidation_codes"`
	LastDailyClose    *ProcessedDailyClosure        `json:"-" bson:"-"`
}

func (p *ProcessedDailyClosure) SetLastDailyClose(last ProcessedDailyClosure) {
	p.LastDailyClose = &last
}

func (p ProcessedDailyClosure) getCalendarPeriod(periodValue *ProcessedDailyClosurePeriod, Filled bool) *measures.DailyReadingClosureCalendarPeriod {

	return &measures.DailyReadingClosureCalendarPeriod{
		Filled:                   Filled,
		Values:                   periodValue.Values,
		InvalidationCodes:        periodValue.InvalidationCodes,
		ValidationStatus:         periodValue.ValidationStatus,
		ValidationStatusIsManual: periodValue.ValidationStatusIsManual,
	}
}

func (p ProcessedDailyClosure) ToDailyReadingClosure() measures.DailyReadingClosure {
	d := measures.DailyReadingClosure{
		ClosureType:       measures.Daily,
		DistributorCode:   p.DistributorCode,
		DistributorID:     p.DistributorID,
		CUPS:              p.CUPS,
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

func (p ProcessedDailyClosure) ToValidatable() validations.MeasureValidatable {
	m := validations.MeasureValidatable{
		Type:        measures.Incremental,
		EndDate:     p.EndDate,
		ReadingDate: p.ReadingDate,
		ServiceType: p.ServiceType,
		PointType:   string(p.PointType),
		P0:          validations.Periods{Values: p.CalendarPeriods.P0.Values},
		P1:          validations.Periods{Values: p.CalendarPeriods.P1.Values},
		P2:          validations.Periods{Values: p.CalendarPeriods.P2.Values},
		P3:          validations.Periods{Values: p.CalendarPeriods.P3.Values},
		P4:          validations.Periods{Values: p.CalendarPeriods.P4.Values},
		P5:          validations.Periods{Values: p.CalendarPeriods.P5.Values},
		P6:          validations.Periods{Values: p.CalendarPeriods.P6.Values},

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

	return m
}

func (p *ProcessedDailyClosure) SetStatusMeasure(validation validations.ValidatorBase) {
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

func (p *ProcessedDailyClosure) GenerateID() {
	id := fmt.Sprintf("%s%s%s%s%s", p.DistributorID, p.EndDate.Format("2006-01-02_15:04"), p.MeterSerialNumber, p.CUPS, p.MeasurePointType)

	p.Id = fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func NewProcessedDailyClosure(meterConfig measures.MeterConfig, date time.Time) *ProcessedDailyClosure {
	p := ProcessedDailyClosure{
		QualityCode:       measures.Complete,
		ValidationStatus:  measures.Valid,
		EndDate:           date.UTC(),
		Magnitudes:        meterConfig.GetMagnitudesActive(),
		DistributorCode:   meterConfig.DistributorCode,
		DistributorID:     meterConfig.DistributorID,
		CUPS:              meterConfig.Cups(),
		MeterSerialNumber: meterConfig.SerialNumber(),
		ServiceType:       meterConfig.ServiceType(),
		PointType:         meterConfig.PointType(),
		MeasurePointType:  meterConfig.MeasurePointType(),
		TariffId:          meterConfig.ContractualSituations.TariffID,
		CalendarCode:      meterConfig.CalendarID,
		MeterType:         meterConfig.MeterType(),
		P1Demand:          meterConfig.ContractualSituations.P1Demand,
		P2Demand:          meterConfig.ContractualSituations.P2Demand,
		P3Demand:          meterConfig.ContractualSituations.P3Demand,
		P4Demand:          meterConfig.ContractualSituations.P4Demand,
		P5Demand:          meterConfig.ContractualSituations.P5Demand,
		P6Demand:          meterConfig.ContractualSituations.P6Demand,
	}
	p.GenerateID()

	return &p
}
