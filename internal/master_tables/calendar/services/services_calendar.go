package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
)

//CalendarServices struct all services
type CalendarServices struct {
	GetCalendars      *GetCalendarService
	InsertCalendar    *InsertCalendarService
	DeleteCalendar    *DeleteCalendarService
	GetOneCalendar    *GetOneCalendarService
	PutCalendar       *ModifyCalendarService
	PutPeriodCalendar *ModifyPeriodCalendarService
	DeletePeriod      *DeletePeriodService
	InsertPeriod      *InsertPeriodService
	GetPeriodCalendar *GetPeriodCalendarService
	GetPeriodById     *GetOnePeriodService
}

//NewCalendarServices create CalendarServices
func NewCalendarServices(repositoryCalendar calendar.RepositoryCalendar, publisher event.PublisherCreator, topic string) *CalendarServices {
	getCalendars := NewGetCalendarService(repositoryCalendar)
	insertCalendar := NewInsertCalendarService(repositoryCalendar)
	deleteCalendar := NewDeleteCalendarService(repositoryCalendar, publisher, topic)
	getOneCalendar := NewGetOneCalendarService(repositoryCalendar)
	putCalendar := NewModifyCalendarService(repositoryCalendar, publisher, topic)
	putPeriodCalendar := NewModifyPeriodCalendarService(repositoryCalendar, publisher, topic)
	deletePeriod := NewDeletePeriodService(repositoryCalendar, publisher, topic)
	insertPeriod := NewInsertPeriodService(repositoryCalendar, publisher, topic)
	getAllPeriodCalendar := NewGetPeriodCalendarService(repositoryCalendar)
	getPeriodById := NewGetOnePeriodService(repositoryCalendar)

	return &CalendarServices{
		GetCalendars:      &getCalendars,
		InsertCalendar:    &insertCalendar,
		DeleteCalendar:    &deleteCalendar,
		GetOneCalendar:    &getOneCalendar,
		PutCalendar:       &putCalendar,
		PutPeriodCalendar: &putPeriodCalendar,
		DeletePeriod:      &deletePeriod,
		InsertPeriod:      &insertPeriod,
		GetPeriodCalendar: &getAllPeriodCalendar,
		GetPeriodById:     &getPeriodById,
	}
}
