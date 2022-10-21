package process_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"crypto/sha256"
	"fmt"
	"time"
)

type ProcessedLoadCurve struct {
	Id                       string                           `bson:"_id"`
	DistributorCode          string                           `json:"distributor_code" bson:"distributor_code"`
	DistributorID            string                           `json:"distributor_id" bson:"distributor_id"`
	CUPS                     string                           `json:"cups" bson:"cups"`
	Period                   measures.PeriodKey               `json:"period" bson:"period"`
	EndDate                  time.Time                        `json:"end_date" bson:"end_date"`
	MeterSerialNumber        string                           `json:"meter_serial_number" bson:"meter_serial_number"`
	GenerationDate           time.Time                        `json:"generation_date" bson:"generation_date"`
	ReadingDate              time.Time                        `json:"reading_date" bson:"reading_date"`
	RegisterType             measures.RegisterType            `json:"type" bson:"type"`
	CurveType                measures.MeasureCurveReadingType `json:"curve_type" bson:"curve_type"`
	ServiceType              measures.ServiceType             `json:"service_type" bson:"service_type"`
	MeterType                measures.MeterType               `json:"meter_type" bson:"meter_type"`
	PointType                measures.PointType               `json:"point_type" bson:"point_type"`
	Origin                   measures.OriginType              `json:"origin" bson:"origin"`
	MeasurePointType         measures.MeasurePointType        `json:"measure_point_type" bson:"measure_point_type"`
	PriorityContract         measures.PriorityContract        `json:"priority_contract" bson:"priority_contract"`
	MeterConfigId            string                           `json:"meter_config_id" bson:"meter_config_id"`
	AI                       float64                          `json:"AI" bson:"AI"`
	AE                       float64                          `json:"AE" bson:"AE"`
	R1                       float64                          `json:"R1" bson:"R1"`
	R2                       float64                          `json:"R2" bson:"R2"`
	R3                       float64                          `json:"R3" bson:"R3"`
	R4                       float64                          `json:"R4" bson:"R4"`
	InvalidationCodes        []string                         `json:"invalidation_codes" bson:"invalidation_codes"`
	ValidationStatus         measures.Status                  `json:"validation_status" bson:"validation_status"`
	ValidationStatusIsManual bool                             `json:"validation_status_is_manual" bson:"validation_status_is_manual"`
	SeasonId                 string                           `json:"season_id" bson:"season_id"`
	DayTypeId                string                           `json:"day_type_id" bson:"day_type_id"`
	P1Demand                 int                              `json:"p1_demand" bson:"p1_demand"`
	P2Demand                 int                              `json:"p2_demand" bson:"p2_demand"`
	P3Demand                 int                              `json:"p3_demand" bson:"p3_demand"`
	P4Demand                 int                              `json:"p4_demand" bson:"p4_demand"`
	P5Demand                 int                              `json:"p5_demand" bson:"p5_demand"`
	P6Demand                 int                              `json:"p6_demand" bson:"p6_demand"`
	Magnitudes               []measures.Magnitude             `json:"magnitudes" bson:"magnitudes"`
}

func (p *ProcessedLoadCurve) GenerateID() {
	id := fmt.Sprintf("%s%s%s%s%s", p.DistributorID, p.EndDate.Format("2006-01-02_15:04"), p.MeterSerialNumber, p.CUPS, p.CurveType)

	p.Id = fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func (p *ProcessedLoadCurve) GetMagnitude(magnitude measures.Magnitude) float64 {
	switch magnitude {
	case measures.AI:
		return p.AI
	case measures.AE:
		return p.AE
	case measures.R1:
		return p.R1
	case measures.R2:
		return p.R2
	case measures.R3:
		return p.R3
	case measures.R4:
		return p.R4
	default:
		return .0
	}
}
func NewProcessedLoadCurve(meterConfig measures.MeterConfig, endDate time.Time, curveType measures.MeasureCurveReadingType) *ProcessedLoadCurve {
	p := ProcessedLoadCurve{
		DistributorID:     meterConfig.DistributorID,
		DistributorCode:   meterConfig.DistributorCode,
		EndDate:           endDate,
		MeterSerialNumber: meterConfig.SerialNumber(),
		CUPS:              meterConfig.Cups(),
		RegisterType:      meterConfig.CurveType,
		CurveType:         curveType,
		ValidationStatus:  measures.Valid,
		ServiceType:       meterConfig.ServiceType(),
		PointType:         meterConfig.PointType(),
		MeterType:         meterConfig.MeterType(),
		MeasurePointType:  meterConfig.MeasurePointType(),
		PriorityContract:  meterConfig.PriorityContract,
		MeterConfigId:     meterConfig.ID,
	}
	p.GenerateID()

	return &p
}

type ProcessedDailyLoadCurve struct {
	Id                       string                           `bson:"_id"`
	DistributorCode          string                           `json:"distributor_code" bson:"distributor_code"`
	DistributorID            string                           `json:"distributor_id" bson:"distributor_id"`
	CUPS                     string                           `json:"cups" bson:"cups"`
	EndDate                  time.Time                        `json:"end_date" bson:"end_date"`
	MeterSerialNumber        string                           `json:"meter_serial_number" bson:"meter_serial_number"`
	GenerationDate           time.Time                        `json:"generation_date" bson:"generation_date"`
	ReadingDate              time.Time                        `json:"reading_date" bson:"reading_date"`
	RegisterType             measures.RegisterType            `json:"type" bson:"type"`
	CurveType                measures.MeasureCurveReadingType `json:"curve_type" bson:"curve_type"`
	ServiceType              measures.ServiceType             `json:"service_type" bson:"service_type"`
	MeterType                measures.MeterType               `json:"meter_type" bson:"meter_type"`
	PointType                measures.PointType               `json:"point_type" bson:"point_type"`
	InvalidationCodes        []string                         `json:"invalidation_codes" bson:"invalidation_codes"`
	ValidationStatus         measures.Status                  `json:"validation_status" bson:"validation_status"`
	ValidationStatusIsManual bool                             `json:"validation_status_is_manual" bson:"validation_status_is_manual"`
	QualityCode              measures.QualityCode             `json:"quality_code" bson:"quality_code"`
}

func (p *ProcessedDailyLoadCurve) GenerateID() {
	id := fmt.Sprintf("%s%s%s%s%s", p.DistributorID, p.EndDate.Format("2006-01-02_15:04"), p.MeterSerialNumber, p.CUPS, p.CurveType)

	p.Id = fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func (p ProcessedLoadCurve) ToValidatable() validations.MeasureValidatable {
	return validations.MeasureValidatable{
		Type:         measures.Incremental,
		EndDate:      p.EndDate,
		ReadingDate:  p.ReadingDate,
		RegisterType: p.RegisterType,
		AI:           p.AI,
		AE:           p.AE,
		R1:           p.R1,
		R2:           p.R2,
		R3:           p.R3,
		R4:           p.R4,
		ServiceType:  p.ServiceType,
		PointType:    string(p.PointType),
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
		Period:   p.Period,
		P1Demand: p.P1Demand,
		P2Demand: p.P2Demand,
		P3Demand: p.P3Demand,
		P4Demand: p.P4Demand,
		P5Demand: p.P5Demand,
		P6Demand: p.P6Demand,
	}
}

func (p *ProcessedLoadCurve) SetStatusMeasure(validation validations.ValidatorBase) {
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

type ListProcessedLoadCurve []ProcessedLoadCurve

func (ms ListProcessedLoadCurve) GetQualityCodeLoadCurve() measures.QualityCode {

	qualityCode := measures.Complete
	hoursFilled := 0

	for _, m := range ms {
		if m.Origin == measures.Filled {
			hoursFilled++
		}
	}

	if hoursFilled != 0 {
		qualityCode = measures.Incomplete
	}

	if len(ms) == hoursFilled {
		qualityCode = measures.Absent
	}

	return qualityCode
}

func (ms ListProcessedLoadCurve) MeasuresBase() []ProcessMeasureBase {
	measuresBase := make([]ProcessMeasureBase, 0, cap(ms))

	for i := range ms {
		measuresBase = append(measuresBase, &ms[i])
	}

	return measuresBase
}
