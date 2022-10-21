package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type CalendarPeriodGenerateService struct {
	calendarRepository       calendar.RepositoryCalendar
	calendarPeriodRepository measures.CalendarPeriodRepository
	tracer                   trace.Tracer
}

func NewCalendarPeriodGenerateService(calendarRepository calendar.RepositoryCalendar, calendarPeriodRepository measures.CalendarPeriodRepository) *CalendarPeriodGenerateService {
	return &CalendarPeriodGenerateService{
		calendarRepository:       calendarRepository,
		calendarPeriodRepository: calendarPeriodRepository,
		tracer:                   telemetry.GetTracer(),
	}
}

func (s CalendarPeriodGenerateService) getDatesRange(startDate, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse("02-01-2006", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := time.Parse("02-01-2006", endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func (s CalendarPeriodGenerateService) addCalendarPeriod(periodMap map[string]*measures.CalendarPeriod, periodCalendar calendar.PeriodCalendar) {
	key := fmt.Sprintf("%s|%s|%s|%s",
		periodCalendar.CalendarCode,
		periodCalendar.GeographicCode,
		periodCalendar.StartDate,
		periodCalendar.EndDate,
	)

	if _, ok := periodMap[key]; !ok {

		startDate, endDate, err := s.getDatesRange(periodCalendar.StartDate, periodCalendar.EndDate)

		if err != nil {
			return
		}

		periodMap[key] = &measures.CalendarPeriod{
			CalendarCode: periodCalendar.CalendarCode,
			GeographicID: periodCalendar.GeographicCode,
			StartDate:    startDate,
			EndDate:      endDate,
		}
	}
	period := periodMap[key]
	for hour := periodCalendar.StartHour; hour < periodCalendar.EndHour+1; hour++ {
		period.SetHour(hour, measures.PeriodKey(periodCalendar.PeriodNumber))
	}

}

func (s CalendarPeriodGenerateService) Handler(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "CalendarPeriodsReGenerate - Handler")
	defer span.End()

	periodsCalendar, err := s.calendarRepository.GetPeriodsActive(ctx)
	if err != nil {
		return err
	}

	err = s.calendarPeriodRepository.DeleteCalendars(ctx)
	if err != nil {
		return err
	}

	calendarPeriodsMap := make(map[string]*measures.CalendarPeriod)
	festivesPeriodsMap := make(map[string]*measures.CalendarPeriod)

	for _, periodCalendar := range periodsCalendar {
		if periodCalendar.CalendarCode == "" || periodCalendar.GeographicCode == "" {
			continue
		}

		if periodCalendar.IsFestive() {
			s.addCalendarPeriod(festivesPeriodsMap, periodCalendar)
			continue
		}

		s.addCalendarPeriod(calendarPeriodsMap, periodCalendar)
	}

	calendarsMap := make(map[string][]measures.CalendarPeriod)

	for _, period := range calendarPeriodsMap {
		key := fmt.Sprintf("%s:%s", period.CalendarCode, period.GeographicID)
		calendars, ok := calendarsMap[key]
		if !ok {
			calendars = make([]measures.CalendarPeriod, 0, 20)
		}
		calendars = append(calendars, *period)

		calendarsMap[key] = calendars
	}

	for _, period := range festivesPeriodsMap {
		key := fmt.Sprintf("%s:%s", period.CalendarCode, period.GeographicID)
		calendars, ok := calendarsMap[key]
		if !ok {
			continue
		}

		for i, calendarPeriod := range calendars {
			if period.StartDate.After(calendarPeriod.StartDate) && period.EndDate.Before(calendarPeriod.EndDate) {
				continue
			}
			calendarPeriod.SetFestivesHours(period.CalendarPeriodHours)
			calendars[i] = calendarPeriod
		}
	}

	for _, periods := range calendarsMap {
		s.calendarPeriodRepository.SaveCalendarPeriod(ctx, periods)
	}

	return nil
}
