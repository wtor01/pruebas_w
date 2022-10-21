package master_tables

import "bitbucket.org/sercide/data-ingestion/pkg/event"

const (
	CalendarPeriodGenerate = "CALENDAR_PERIODS/GENERATE"
	FestiveDaysGenerate    = "FESTIVE_DAYS/GENERATE"
)

type CalendarPeriodsPayload struct {
}

type CalendarPeriodsEvent = event.Message[CalendarPeriodsPayload]

func NewCalendarPeriodsEvent(t string, payload CalendarPeriodsPayload) CalendarPeriodsEvent {
	return CalendarPeriodsEvent{
		Type:    t,
		Payload: payload,
	}
}

func NewCalendarRedisGenerateEvent() CalendarPeriodsEvent {
	return NewCalendarPeriodsEvent(CalendarPeriodGenerate, CalendarPeriodsPayload{})
}

func NewFestiveDaysGenerateEvent() CalendarPeriodsEvent {
	return NewCalendarPeriodsEvent(FestiveDaysGenerate, CalendarPeriodsPayload{})
}
