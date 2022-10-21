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

type MeasureCurveWrite struct {
	Id                string                           `bson:"_id"`
	EndDate           time.Time                        `json:"end_date" bson:"end_date"`
	ReadingDate       time.Time                        `json:"reading_date" bson:"reading_date"`
	GenerationDate    time.Time                        `json:"generation_date" bson:"generation_date"`
	Type              measures.Type                    `json:"type" bson:"type"`
	Status            measures.Status                  `json:"status" bson:"status"`
	ReadingType       measures.ReadingType             `json:"reading_type" bson:"reading_type"`
	CurveType         measures.MeasureCurveReadingType `json:"curve_type" bson:"curve_type"`
	Contract          string                           `json:"contract" bson:"contract"`
	MeterSerialNumber string                           `json:"meter_serial_number" bson:"meter_serial_number"`
	ConcentratorID    string                           `json:"concentrator_id" bson:"concentrator_id"`
	File              string                           `json:"file" bson:"file"`
	DistributorID     string                           `json:"distributor_id" bson:"distributor_id"`
	DistributorCDOS   string                           `json:"distributor_cdos" bson:"distributor_cdos"`
	Origin            measures.OriginType              `json:"origin" bson:"origin"`
	Qualifier         string                           `json:"qualifier" bson:"qualifier"`
	AI                float64                          `json:"AI" bson:"AI"`
	AE                float64                          `json:"AE" bson:"AE"`
	R1                float64                          `json:"R1" bson:"R1"`
	R2                float64                          `json:"R2" bson:"R2"`
	R3                float64                          `json:"R3" bson:"R3"`
	R4                float64                          `json:"R4" bson:"R4"`
	Invalidations     []string                         `json:"invalidations" bson:"invalidations"`
}

type MeasureCurveMeterSerialNumber struct {
	Day               int    `json:"day" bson:"day"`
	Month             int    `json:"month" bson:"month"`
	Year              int    `json:"year" bson:"year"`
	MeterSerialNumber string `json:"meter_serial_number" bson:"meter_serial_number"`
}

func (c *MeasureCurveWrite) GenerateID() {
	id := fmt.Sprintf("%s%s%s%s%s%s", c.DistributorID, c.EndDate.Format("2006-01-02_15:04:05"), c.MeterSerialNumber, c.File, c.Type, c.CurveType)

	c.Id = fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func (c *MeasureCurveWrite) GetID() string {
	return c.Id
}

func (c *MeasureCurveWrite) GetMeterSerialNumber() string {
	return c.MeterSerialNumber
}

func (c *MeasureCurveWrite) GetEndDate() time.Time {
	return c.EndDate
}

func (c *MeasureCurveWrite) GetDistributorCDOS() string {
	return c.DistributorCDOS
}

func (c *MeasureCurveWrite) GetDistributorID() string {
	return c.DistributorID
}

func (c *MeasureCurveWrite) GetReadingType() measures.ReadingType {
	return c.ReadingType
}

func (c *MeasureCurveWrite) SetDistributorID(distributorId string) {
	c.DistributorID = distributorId
}

func (c *MeasureCurveWrite) SetStatusMeasure(validation validations.ValidatorBase) {
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

func (c *MeasureCurveWrite) SetStatus(status measures.Status) {
	c.Status = status
}

func (c *MeasureCurveWrite) ToValidatable() []validations.MeasureValidatable {
	return []validations.MeasureValidatable{
		{
			Type:        c.Type,
			EndDate:     c.EndDate,
			ReadingDate: c.ReadingDate,
			AI:          c.AI,
			AE:          c.AE,
			R1:          c.R1,
			R2:          c.R2,
			R3:          c.R3,
			R4:          c.R4,
			Qualifier:   c.Qualifier,
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
		},
	}
}

func (c *MeasureCurveWrite) GetOriginFile() (originFile string) {
	return path.Base(c.File)
}

func ListMeasureCurveWriteToEvents(measuresCurve []MeasureCurveWrite, maxInEvent int) []InsertMeasureCurveEvent {
	measuresCurvGrouped := make([]MeasureCurveWrite, 0, maxInEvent)
	numberEvents := int64(math.Ceil(float64(len(measuresCurve)) / float64(maxInEvent)))
	events := make([]InsertMeasureCurveEvent, 0, numberEvents)

	for _, m := range measuresCurve {
		measuresCurvGrouped = append(measuresCurvGrouped, m)
		if len(measuresCurvGrouped) < maxInEvent {
			continue
		}
		events = append(events, NewInsertMeasureCurveEvent(measuresCurvGrouped))
		measuresCurvGrouped = make([]MeasureCurveWrite, 0, maxInEvent)
	}

	if len(measuresCurvGrouped) != 0 {
		events = append(events, NewInsertMeasureCurveEvent(measuresCurvGrouped))
		measuresCurvGrouped = make([]MeasureCurveWrite, 0, maxInEvent)
	}

	return events
}
