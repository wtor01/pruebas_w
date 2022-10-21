package process_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"time"
)

type Types = string

const (
	ProcessCurveEventType              = "PROCESS_MEASURES/CURVE"
	ProcessBillingClosureEventType     = "PROCESS_MEASURES/BILLING_CLOSURE"
	ProcessDailyClosureEventType       = "PROCESS_MEASURES/DAILY_CLOSURE"
	SchedulerEventType                 = "PROCESS_MEASURES/INIT"
	SchedulerReprocessingEventType     = "REPROCESSING_MEASURES/INIT"
	TypeDistributorProcess             = "PROCESS_MEASURES/DISTRIBUTOR"
	TypeReprocessingDistributorProcess = "REPROCESSING_MEASURES/DISTRIBUTOR"
	TypeReprocessingMeterProcess       = "REPROCESSING_MEASURES/METER"
)

type ReSchedulerEvent = event.Message[ReSchedulerEventPayload]
type ReSchedulerMeterEvent = event.Message[ReSchedulerMeterPayload]

func NewReSchedulerEventPayload(t string, payload ReSchedulerEventPayload) ReSchedulerEvent {
	return ReSchedulerEvent{
		Type:    t,
		Payload: payload,
	}
}
func NewReSchedulerMeterPayload(t string, payload ReSchedulerMeterPayload) ReSchedulerMeterEvent {
	return ReSchedulerMeterEvent{
		Type:    t,
		Payload: payload,
	}
}

type ReSchedulerMeterPayload struct {
	MeterSerialNumber string               `json:"meter_serial_number"`
	Date              time.Time            `json:"date"`
	ReadingType       measures.ReadingType `json:"reading_type"`
}
type ReSchedulerEventPayload struct {
	ID            string               `json:"id"`
	DistributorId string               `json:"distributor_id"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	ServiceType   measures.ServiceType `json:"service_type"`
	PointType     string               `json:"point_type"`
	MeterType     []string             `json:"meter_type"`
	ReadingType   measures.ReadingType `json:"reading_type"`
	Format        string               `json:"format"`
	Limit         int                  `json:"limit"`
	Date          time.Time            `json:"date"`
	Offset        int                  `json:"offset"`
	StartDate     time.Time            `json:"start_date"`
	EndDate       time.Time            `json:"end_date"`
}

func NewProcessCurveEvent(date time.Time, meterConfig measures.MeterConfig, curveType measures.MeasureCurveReadingType) measures.ProcessMeasureEvent {
	return measures.NewProcessMeasureCurveEvent(ProcessCurveEventType, date, meterConfig, curveType)
}

func NewProcessBillingClosureEvent(date time.Time, meterConfig measures.MeterConfig) measures.ProcessMeasureEvent {
	return measures.NewProcessMeasureEvent(ProcessBillingClosureEventType, date, meterConfig)
}

func NewProcessDailyClosureEvent(date time.Time, meterConfig measures.MeterConfig) measures.ProcessMeasureEvent {
	return measures.NewProcessMeasureEvent(ProcessDailyClosureEventType, date, meterConfig)
}

func NewProcessByDistributorEvent(payload measures.SchedulerEventPayload) measures.SchedulerEvent {
	return measures.NewSchedulerEvent(TypeDistributorProcess, payload)
}

// NewReprocessingProcessByDistributorEvent Create a payload for Rescheduler distributors
func NewReprocessingProcessByDistributorEvent(payload ReSchedulerEventPayload) ReSchedulerEvent {
	return NewReSchedulerEventPayload(TypeReprocessingDistributorProcess, payload)
}
func NewReprocessingMeterEvent(payload ReSchedulerMeterPayload) ReSchedulerMeterEvent {
	return NewReSchedulerMeterPayload(TypeReprocessingMeterProcess, payload)

}
