package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

type CalendarPeriodsPubSubServices struct {
	CalendarPeriodGenerateService *CalendarPeriodGenerateService
	FestiveDaysGenerateService    *FestiveDaysGenerateService
}

func NewCalendarPeriodsPubSubServices(calendarRepository calendar.RepositoryCalendar, calendarPeriodRepository measures.CalendarPeriodRepository, festiveDaysRepository festive_days.FestiveDayRepository) *CalendarPeriodsPubSubServices {
	return &CalendarPeriodsPubSubServices{
		CalendarPeriodGenerateService: NewCalendarPeriodGenerateService(calendarRepository, calendarPeriodRepository),
		FestiveDaysGenerateService:    NewFestiveDaysGenerateService(festiveDaysRepository, calendarPeriodRepository),
	}
}
