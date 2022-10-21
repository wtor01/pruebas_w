package calendar

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
)

func CalendarToResponse(d calendar.Calendar) CalendarWithId {
	return CalendarWithId{
		Id:             d.Id,
		Code:           d.Id,
		Description:    d.Description,
		Periods:        d.Periods,
		GeographicCode: d.GeographicCode,
	}
}

func PCtoResponse(calendar calendar.PeriodCalendar) CalendarPeriod {
	return CalendarPeriod{
		Id:           calendar.ID.String(),
		CalendarCode: calendar.CalendarCode,
		Description:  calendar.Description,
		DayType:      CalendarPeriodDayType(calendar.DayType),
		StartDate:    calendar.StartDate,
		EndDate:      calendar.EndDate,
		PeriodNumber: CalendarPeriodPeriodNumber(calendar.PeriodNumber),
		Year:         calendar.Year,
		StartHour:    calendar.StartHour,
		EndHour:      calendar.EndHour,
		Energy:       calendar.Energy,
		Power:        calendar.Power,
	}
}
