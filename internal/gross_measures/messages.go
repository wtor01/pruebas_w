package gross_measures

import (
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"time"
)

const (
	InsertMeasureCurveEventType      = "MEASURES/INSERT_MEASURE_CURVE"
	InsertMeasureCloseEventType      = "MEASURES/INSERT_MEASURE_CLOSE"
	RequestToRecoverMeasureEventType = "PROCESS/TO_RECOVER_MEASURES"
)

var MaxMeasuresInEvent = 100

type InsertMeasureCurveEvent = event.Message[[]MeasureCurveWrite]

func NewInsertMeasureCurveEvent(measures []MeasureCurveWrite) event.Message[[]MeasureCurveWrite] {
	return event.Message[[]MeasureCurveWrite]{
		Type:    InsertMeasureCurveEventType,
		Payload: measures,
	}
}

type InsertMeasureCloseEvent = event.Message[[]MeasureCloseWrite]

func NewInsertMeasureCloseEvent(measures []MeasureCloseWrite) InsertMeasureCloseEvent {
	return InsertMeasureCloseEvent{
		Type:    InsertMeasureCloseEventType,
		Payload: measures,
	}
}

type RequestToRecoverMeasurePayload struct {
	DistributorId     string    `json:"distributor_id"`
	MeterSerialNumber string    `json:"meter_serial_number"`
	Time              time.Time `json:"time"`
	Type              string    `json:"type"`
}

type RequestToRecoverMeasureEvent = event.Message[RequestToRecoverMeasurePayload]

func NewRequestToRecoverMeasureEvent(distributorId string, serialNumber string, time time.Time, Type string) RequestToRecoverMeasureEvent {
	return RequestToRecoverMeasureEvent{
		Type: RequestToRecoverMeasureEventType,
		Payload: RequestToRecoverMeasurePayload{
			DistributorId:     distributorId,
			MeterSerialNumber: serialNumber,
			Time:              time,
			Type:              Type,
		},
	}
}
