package measures

import (
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"time"
)

type SchedulerEvent = event.Message[SchedulerEventPayload]

func NewSchedulerEvent(t string, payload SchedulerEventPayload) SchedulerEvent {
	return SchedulerEvent{
		Type:    t,
		Payload: payload,
	}
}

type SchedulerEventPayload struct {
	ID            string      `json:"id"`
	DistributorId string      `json:"distributor_id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	ServiceType   ServiceType `json:"service_type"`
	PointType     string      `json:"point_type"`
	MeterType     []string    `json:"meter_type"`
	ReadingType   ReadingType `json:"reading_type"`
	Format        string      `json:"format"`
	Date          time.Time   `json:"date"`
	Limit         int         `json:"limit"`
	Offset        int         `json:"offset"`
}

type ProcessMeasurePayload struct {
	MeterConfig       MeterConfig             `json:"meter_config"`
	Date              time.Time               `json:"date"`
	RecoverMeasuresId string                  `json:"recover_measures_id"`
	CurveType         MeasureCurveReadingType `json:"curve_type"`
}

type ProcessMeasureEvent = event.Message[ProcessMeasurePayload]

func NewProcessMeasureEvent(t string, date time.Time, meterConfig MeterConfig) ProcessMeasureEvent {
	return ProcessMeasureEvent{
		Type: t,
		Payload: ProcessMeasurePayload{
			MeterConfig: meterConfig,
			Date:        date,
		},
	}
}

func NewProcessMeasureCurveEvent(t string, date time.Time, meterConfig MeterConfig, curveType MeasureCurveReadingType) ProcessMeasureEvent {
	return ProcessMeasureEvent{
		Type: t,
		Payload: ProcessMeasurePayload{
			MeterConfig: meterConfig,
			Date:        date,
			CurveType:   curveType,
		},
	}
}
