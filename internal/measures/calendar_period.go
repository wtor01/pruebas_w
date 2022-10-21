package measures

import (
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"sort"
	"time"
)

type CalendarPeriodHours struct {
	Hour1  PeriodKey
	Hour2  PeriodKey
	Hour3  PeriodKey
	Hour4  PeriodKey
	Hour5  PeriodKey
	Hour6  PeriodKey
	Hour7  PeriodKey
	Hour8  PeriodKey
	Hour9  PeriodKey
	Hour10 PeriodKey
	Hour11 PeriodKey
	Hour12 PeriodKey
	Hour13 PeriodKey
	Hour14 PeriodKey
	Hour15 PeriodKey
	Hour16 PeriodKey
	Hour17 PeriodKey
	Hour18 PeriodKey
	Hour19 PeriodKey
	Hour20 PeriodKey
	Hour21 PeriodKey
	Hour22 PeriodKey
	Hour23 PeriodKey
	Hour00 PeriodKey
}

func (periodHours CalendarPeriodHours) IsAllHourFill() bool {
	for i := 0; i < 24; i++ {
		period := periodHours.GetHour(i)

		if period == "" {
			return false
		}
	}
	return true
}

func (periodHours *CalendarPeriodHours) SetHour(hour int, period PeriodKey) {
	switch hour {
	case 1:
		periodHours.Hour1 = period
	case 2:
		periodHours.Hour2 = period
	case 3:
		periodHours.Hour3 = period
	case 4:
		periodHours.Hour4 = period
	case 5:
		periodHours.Hour5 = period
	case 6:
		periodHours.Hour6 = period
	case 7:
		periodHours.Hour7 = period
	case 8:
		periodHours.Hour8 = period
	case 9:
		periodHours.Hour9 = period
	case 10:
		periodHours.Hour10 = period
	case 11:
		periodHours.Hour11 = period
	case 12:
		periodHours.Hour12 = period
	case 13:
		periodHours.Hour13 = period
	case 14:
		periodHours.Hour14 = period
	case 15:
		periodHours.Hour15 = period
	case 16:
		periodHours.Hour16 = period
	case 17:
		periodHours.Hour17 = period
	case 18:
		periodHours.Hour18 = period
	case 19:
		periodHours.Hour19 = period
	case 20:
		periodHours.Hour20 = period
	case 21:
		periodHours.Hour21 = period
	case 22:
		periodHours.Hour22 = period
	case 23:
		periodHours.Hour23 = period
	case 0:
		periodHours.Hour00 = period

	}
}

func (periodHours CalendarPeriodHours) GetHour(hour int) PeriodKey {
	switch hour {
	case 1:
		return periodHours.Hour1
	case 2:
		return periodHours.Hour2
	case 3:
		return periodHours.Hour3
	case 4:
		return periodHours.Hour4
	case 5:
		return periodHours.Hour5
	case 6:
		return periodHours.Hour6
	case 7:
		return periodHours.Hour7
	case 8:
		return periodHours.Hour8
	case 9:
		return periodHours.Hour9
	case 10:
		return periodHours.Hour10
	case 11:
		return periodHours.Hour11
	case 12:
		return periodHours.Hour12
	case 13:
		return periodHours.Hour13
	case 14:
		return periodHours.Hour14
	case 15:
		return periodHours.Hour15
	case 16:
		return periodHours.Hour16
	case 17:
		return periodHours.Hour17
	case 18:
		return periodHours.Hour18
	case 19:
		return periodHours.Hour19
	case 20:
		return periodHours.Hour20
	case 21:
		return periodHours.Hour21
	case 22:
		return periodHours.Hour22
	case 23:
		return periodHours.Hour23
	case 0:
		return periodHours.Hour00
	default:
		return ""
	}
}

func (periodHours CalendarPeriodHours) GetPeriods() []PeriodKey {
	periods := utils.NewSet(make([]PeriodKey, 0))
	for i := 0; i < 24; i++ {
		periodHour := periodHours.GetHour(i)
		if periodHour == "" {
			continue
		}
		periods.Add(periodHour)
	}

	return periods.Slice()
}

type CalendarPeriod struct {
	CalendarCode string
	GeographicID string
	CalendarPeriodHours
	Festive   CalendarPeriodHours
	StartDate time.Time
	EndDate   time.Time
	Day       time.Time
	IsFestive bool
}

// GetHourPeriod al estar los periodos definidos en fecha Española debemos pasar la hora de la medida a horario español
// y restarle una hora si el equipo lee horariamente o 15min ya que la lectura de las 22 en utc corresponde
// con el periodo de las 21
func (c CalendarPeriod) GetHourPeriod(d time.Time, curveType MeasureCurveReadingType, isFestive bool, loc *time.Location) PeriodKey {
	date := d.In(loc)
	if curveType == HourlyMeasureCurveReadingType {
		date = date.Add(-time.Hour)
	} else {
		date = date.Add(-time.Minute * 15)
	}

	periods := c.CalendarPeriodHours

	if isFestive {
		periods = c.Festive
	}

	return periods.GetHour(date.Hour())
}

func (c CalendarPeriod) IsFestiveDay() bool {
	if c.IsFestive || c.Day.Weekday() == time.Sunday || c.Day.Weekday() == time.Saturday {
		return true
	}

	return false
}

func (c CalendarPeriod) GetAllPeriods() []PeriodKey {

	dailyPeriods := c.GetPeriods()
	festivePeriods := c.Festive.GetPeriods()

	set := utils.NewSet(append(dailyPeriods, festivePeriods...))
	periods := set.Slice()

	sort.Slice(periods, func(i, j int) bool {
		return periods[i] < periods[j]
	})

	return periods
}

func (c CalendarPeriod) GetDailyPeriods() []PeriodKey {
	return c.GetPeriods()
}

func (c *CalendarPeriod) SetFestivesHours(hours CalendarPeriodHours) {
	c.Festive = hours
}

type FestiveDay struct {
	Date         time.Time
	GeographicID string
}

type SearchFestiveDay struct {
	Date         time.Time
	GeographicID string
}

type SearchCalendarPeriod struct {
	Day          time.Time
	GeographicID string
	CalendarCode string
	Location     *time.Location
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=CalendarPeriodRepository
type CalendarPeriodRepository interface {
	SaveFestiveDay(ctx context.Context, day FestiveDay) error
	GetFestiveDay(ctx context.Context, search SearchFestiveDay) (FestiveDay, error)
	GetCalendarPeriod(ctx context.Context, search SearchCalendarPeriod) (CalendarPeriod, error)
	SaveCalendarPeriod(ctx context.Context, calendarPeriods []CalendarPeriod) error
	DeleteCalendars(ctx context.Context) error
	DeleteFestiveDays(ctx context.Context) error
}
