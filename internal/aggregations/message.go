package aggregations

import (
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"time"
)

const (
	SchedulerEventType            = "AGGREGATIONS/INIT"
	SchedulerDistributorEventType = "AGGREGATIONS/PROCESS_DISTRIBUTOR"
)

type ConfigScheduler struct {
	DistributorId   string    `json:"distributor_id"`
	DistributorCDOS string    `json:"distributor_cdos"`
	Date            time.Time `json:"date"`
	Config
}

type SchedulerEvent = event.Message[ConfigScheduler]

func NewSchedulerEvent[P any](payload P) event.Message[P] {
	return event.Message[P]{
		Type:    SchedulerEventType,
		Payload: payload,
	}
}

func NewSchedulerDistributorEvent[P any](payload P) event.Message[P] {
	return event.Message[P]{
		Type:    SchedulerDistributorEventType,
		Payload: payload,
	}
}
