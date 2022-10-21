package gross_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"crypto/sha256"
	"fmt"
	"math"
	"path"
	"time"
)

type MeasureClosePeriod struct {
	Period measures.PeriodKey `json:"period" bson:"period"`
	AI     float64            `json:"AI" bson:"AI"`
	AE     float64            `json:"AE" bson:"AE"`
	R1     float64            `json:"R1" bson:"R1"`
	R2     float64            `json:"R2" bson:"R2"`
	R3     float64            `json:"R3" bson:"R3"`
	R4     float64            `json:"R4" bson:"R4"`
	MX     float64            `json:"MX" bson:"MX"`
	E      float64            `json:"E" bson:"E"`
	FX     time.Time          `json:"FX" bson:"FX"`
}

type MeasureClose struct {
	StartDate         time.Time                   `json:"start_date" bson:"start_date"`
	EndDate           time.Time                   `json:"end_date" bson:"end_date"`
	ReadingDate       time.Time                   `json:"reading_date" bson:"reading_date"`
	Type              measures.Type               `json:"type" bson:"type"`
	Status            measures.Status             `json:"status" bson:"status"`
	ReadingType       measures.ReadingType        `json:"reading_type" bson:"reading_type"`
	Contract          string                      `json:"contract" bson:"contract"`
	RegisterType      measures.RegisterType       `json:"register_type" bson:"register_type"`
	MeterID           string                      `json:"meter_id" bson:"meter_id"`
	CUPS              string                      `json:"cups" bson:"cups"`
	MeterSerialNumber string                      `json:"meter_serial_number" bson:"meter_serial_number"`
	MeasurePointType  measures.MeasurePointType   `json:"measure_point_type" bson:"measure_point_type"`
	ConcentratorID    string                      `json:"concentrator_id" bson:"concentrator_id"`
	ReaderID          string                      `json:"reader_id" bson:"reader_id"`
	File              string                      `json:"file" bson:"file"`
	DistributorID     string                      `json:"distributor_id" bson:"distributor_id"`
	DistributorCDOS   string                      `json:"distributor_cdos" bson:"distributor_cdos"`
	Origin            measures.OriginType         `json:"origin" bson:"origin"`
	Qualifier         string                      `json:"qualifier" bson:"qualifier"`
	Periods           []MeasureClosePeriod        `json:"periods" bson:"periods"`
	Invalidations     []string                    `json:"invalidations" bson:"invalidations"`
	MeterType         measures.MeterType          `json:"meter_type"bson:"meter_type"`
	TlgCode           measures.MeterConfigTlgCode `json:"tlg_code"bson:"tlg_code"`
}

type MeasureCloseWrite struct {
	Id                string               `bson:"_id"`
	StartDate         time.Time            `json:"start_date" bson:"start_date"`
	EndDate           time.Time            `json:"end_date" bson:"end_date"`
	ReadingDate       time.Time            `json:"reading_date" bson:"reading_date"`
	GenerationDate    time.Time            `json:"generation_date" bson:"generation_date"`
	Type              measures.Type        `json:"type" bson:"type"`
	Status            measures.Status      `json:"status" bson:"status"`
	ReadingType       measures.ReadingType `json:"reading_type" bson:"reading_type"`
	Contract          string               `json:"contract" bson:"contract"`
	MeterSerialNumber string               `json:"meter_serial_number" bson:"meter_serial_number"`
	ConcentratorID    string               `json:"concentrator_id" bson:"concentrator_id"`
	File              string               `json:"file" bson:"file"`
	DistributorID     string               `json:"distributor_id" bson:"distributor_id"`
	DistributorCDOS   string               `json:"distributor_cdos" bson:"distributor_cdos"`
	Origin            measures.OriginType  `json:"origin" bson:"origin"`
	Qualifier         string               `json:"qualifier" bson:"qualifier"`
	Periods           []MeasureClosePeriod `json:"periods" bson:"periods"`
	Invalidations     []string             `json:"invalidations" bson:"invalidations"`
}

func (c *MeasureCloseWrite) GenerateID() {
	id := fmt.Sprintf("%s%s%s%s%s%s", c.DistributorID, c.EndDate.Format("2006-01-02_15:04:05"), c.MeterSerialNumber, c.File, c.Type, c.ReadingType)
	c.Id = fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func (c *MeasureCloseWrite) GetID() string {
	return c.Id
}

func (c *MeasureCloseWrite) GetMeterSerialNumber() string {
	return c.MeterSerialNumber
}

func (c *MeasureCloseWrite) GetEndDate() time.Time {
	return c.EndDate
}

func (c *MeasureCloseWrite) GetDistributorCDOS() string {
	return c.DistributorCDOS
}

func (c *MeasureCloseWrite) GetDistributorID() string {
	return c.DistributorID
}

func (c *MeasureCloseWrite) GetReadingType() measures.ReadingType {
	return c.ReadingType
}

func (c *MeasureCloseWrite) SetDistributorID(distributorId string) {
	c.DistributorID = distributorId
}

func (c *MeasureCloseWrite) SetStatus(status measures.Status) {
	c.Status = status
}

func (c *MeasureCloseWrite) GetOriginFile() (originFile string) {
	return path.Base(c.File)
}

func (c *MeasureCloseWrite) SetStatusMeasure(validation validations.ValidatorBase) {
	status := measures.StatusValue
	if !c.Status.Compare(status, validation.Action) {
		c.Status = validation.Action
	}

	if c.Invalidations == nil {
		c.Invalidations = make([]string, 0, 1)
	}

	for _, i := range c.Invalidations {
		if i == validation.Code {
			return
		}
	}

	c.Invalidations = append(c.Invalidations, validation.Code)
}

func (c *MeasureCloseWrite) ToValidatable() []validations.MeasureValidatable {
	result := make([]validations.MeasureValidatable, 0, len(c.Periods))

	for _, p := range c.Periods {
		result = append(result, validations.MeasureValidatable{
			Type:        c.Type,
			StartDate:   c.StartDate,
			EndDate:     c.EndDate,
			ReadingDate: c.ReadingDate,
			AI:          p.AI,
			AE:          p.AE,
			R1:          p.R1,
			R2:          p.R2,
			R3:          p.R3,
			R4:          p.R4,
			MX:          p.MX,
			E:           p.E,
			FX:          p.FX,
			Qualifier:   "",
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
				validations.MX:          {},
				validations.E:           {},
				validations.FX:          {},
			},
		})
	}

	return result
}

func (c MeasureClose) NumberPeriods() int {
	return 6
}

func ListMeasureCloseWriteToEvents(measuresClose []MeasureCloseWrite, maxInEvent int) []InsertMeasureCloseEvent {
	measuresCloseGrouped := make([]MeasureCloseWrite, 0, maxInEvent)
	numberEvents := int64(math.Ceil(float64(len(measuresClose)) / float64(maxInEvent)))
	events := make([]InsertMeasureCloseEvent, 0, numberEvents)

	for _, m := range measuresClose {
		measuresCloseGrouped = append(measuresCloseGrouped, m)
		if len(measuresCloseGrouped) < maxInEvent {
			continue
		}
		events = append(events, NewInsertMeasureCloseEvent(measuresCloseGrouped))
		measuresCloseGrouped = make([]MeasureCloseWrite, 0, maxInEvent)
	}

	if len(measuresCloseGrouped) != 0 {
		events = append(events, NewInsertMeasureCloseEvent(measuresCloseGrouped))
		measuresCloseGrouped = make([]MeasureCloseWrite, 0, maxInEvent)
	}

	return events
}
