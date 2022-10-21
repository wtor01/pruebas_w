package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"time"
)

type Types = string

const (
	SchedulerEventType       = "BILLING_MEASURES/INIT"
	TypeDistributorProcess   = "BILLING_MEASURES/DISTRIBUTOR"
	ProcessMVH               = "BILLING_MEASURES/PROCESS_MVH"
	SelfConsumptionEventType = "BILLING_MEASURES/SELF_CONSUMPTION"
)

type MessageDistributorProcess struct {
	Id            string    `json:"id"`
	DistributorId string    `json:"distributor_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ServiceType   string    `json:"service_type"`
	PointType     string    `json:"point_type"`
	MeterType     []string  `json:"meter_type"`
	ReadingType   string    `json:"reading_type"`
	Format        string    `json:"format"`
	Time          time.Time `json:"time"`
}

type OnSaveBillingMeasurePayload struct {
	InitDate         time.Time `json:"init_date"`
	EndDate          time.Time `json:"end_date"`
	CUPS             string    `json:"CUPS"`
	BillingMeasureId string    `json:"billing_measure_id"`
}

type OnSaveBillingMeasureEvent = event.Message[OnSaveBillingMeasurePayload]

func NewProcessByDistributorEvent(payload measures.SchedulerEventPayload) measures.SchedulerEvent {
	return measures.NewSchedulerEvent(TypeDistributorProcess, payload)
}

func NewProcessMvhEvent(date time.Time, meterConfig measures.MeterConfig) measures.ProcessMeasureEvent {
	return measures.NewProcessMeasureEvent(ProcessMVH, date, meterConfig)
}

func NewOnSaveBillingMeasureEvent(payload OnSaveBillingMeasurePayload) OnSaveBillingMeasureEvent {
	return OnSaveBillingMeasureEvent{
		Type:    SelfConsumptionEventType,
		Payload: payload,
	}
}
