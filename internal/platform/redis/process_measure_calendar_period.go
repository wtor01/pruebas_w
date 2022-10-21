package redis

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"time"
)

type ProcessMeasureFestiveDays struct {
	RedisCache *DataCacheRedis
}

type CalendarPeriodHours struct {
	Hour1  measures.PeriodKey
	Hour2  measures.PeriodKey
	Hour3  measures.PeriodKey
	Hour4  measures.PeriodKey
	Hour5  measures.PeriodKey
	Hour6  measures.PeriodKey
	Hour7  measures.PeriodKey
	Hour8  measures.PeriodKey
	Hour9  measures.PeriodKey
	Hour10 measures.PeriodKey
	Hour11 measures.PeriodKey
	Hour12 measures.PeriodKey
	Hour13 measures.PeriodKey
	Hour14 measures.PeriodKey
	Hour15 measures.PeriodKey
	Hour16 measures.PeriodKey
	Hour17 measures.PeriodKey
	Hour18 measures.PeriodKey
	Hour19 measures.PeriodKey
	Hour20 measures.PeriodKey
	Hour21 measures.PeriodKey
	Hour22 measures.PeriodKey
	Hour23 measures.PeriodKey
	Hour00 measures.PeriodKey
}
type CalendarPeriod struct {
	CalendarCode string
	GeographicID string
	CalendarPeriodHours
	Festive   CalendarPeriodHours
	StartDate time.Time
	EndDate   time.Time
}

func (c CalendarPeriod) toDomain(day time.Time, isFestive bool) measures.CalendarPeriod {
	return measures.CalendarPeriod{
		CalendarCode: c.CalendarCode,
		GeographicID: c.GeographicID,
		CalendarPeriodHours: measures.CalendarPeriodHours{
			Hour1:  c.Hour1,
			Hour2:  c.Hour2,
			Hour3:  c.Hour3,
			Hour4:  c.Hour4,
			Hour5:  c.Hour5,
			Hour6:  c.Hour6,
			Hour7:  c.Hour7,
			Hour8:  c.Hour8,
			Hour9:  c.Hour9,
			Hour10: c.Hour10,
			Hour11: c.Hour11,
			Hour12: c.Hour12,
			Hour13: c.Hour13,
			Hour14: c.Hour14,
			Hour15: c.Hour15,
			Hour16: c.Hour16,
			Hour17: c.Hour17,
			Hour18: c.Hour18,
			Hour19: c.Hour19,
			Hour20: c.Hour20,
			Hour21: c.Hour21,
			Hour22: c.Hour22,
			Hour23: c.Hour23,
			Hour00: c.Hour00,
		},
		Festive: measures.CalendarPeriodHours{
			Hour1:  c.Festive.Hour1,
			Hour2:  c.Festive.Hour2,
			Hour3:  c.Festive.Hour3,
			Hour4:  c.Festive.Hour4,
			Hour5:  c.Festive.Hour5,
			Hour6:  c.Festive.Hour6,
			Hour7:  c.Festive.Hour7,
			Hour8:  c.Festive.Hour8,
			Hour9:  c.Festive.Hour9,
			Hour10: c.Festive.Hour10,
			Hour11: c.Festive.Hour11,
			Hour12: c.Festive.Hour12,
			Hour13: c.Festive.Hour13,
			Hour14: c.Festive.Hour14,
			Hour15: c.Festive.Hour15,
			Hour16: c.Festive.Hour16,
			Hour17: c.Festive.Hour17,
			Hour18: c.Festive.Hour18,
			Hour19: c.Festive.Hour19,
			Hour20: c.Festive.Hour20,
			Hour21: c.Festive.Hour21,
			Hour22: c.Festive.Hour22,
			Hour23: c.Festive.Hour23,
			Hour00: c.Festive.Hour00,
		},
		StartDate: c.StartDate,
		EndDate:   c.EndDate,
		Day:       day,
		IsFestive: isFestive,
	}
}

func NewProcessMeasureFestiveDays(client *redis.Client) *ProcessMeasureFestiveDays {
	return &ProcessMeasureFestiveDays{RedisCache: NewDataCacheRedis(client)}
}

func (repository ProcessMeasureFestiveDays) toDb(calendarPeriod measures.CalendarPeriod) CalendarPeriod {
	return CalendarPeriod{
		CalendarCode: calendarPeriod.CalendarCode,
		GeographicID: calendarPeriod.GeographicID,
		CalendarPeriodHours: CalendarPeriodHours{
			Hour1:  calendarPeriod.Hour1,
			Hour2:  calendarPeriod.Hour2,
			Hour3:  calendarPeriod.Hour3,
			Hour4:  calendarPeriod.Hour4,
			Hour5:  calendarPeriod.Hour5,
			Hour6:  calendarPeriod.Hour6,
			Hour7:  calendarPeriod.Hour7,
			Hour8:  calendarPeriod.Hour8,
			Hour9:  calendarPeriod.Hour9,
			Hour10: calendarPeriod.Hour10,
			Hour11: calendarPeriod.Hour11,
			Hour12: calendarPeriod.Hour12,
			Hour13: calendarPeriod.Hour13,
			Hour14: calendarPeriod.Hour14,
			Hour15: calendarPeriod.Hour15,
			Hour16: calendarPeriod.Hour16,
			Hour17: calendarPeriod.Hour17,
			Hour18: calendarPeriod.Hour18,
			Hour19: calendarPeriod.Hour19,
			Hour20: calendarPeriod.Hour20,
			Hour21: calendarPeriod.Hour21,
			Hour22: calendarPeriod.Hour22,
			Hour23: calendarPeriod.Hour23,
			Hour00: calendarPeriod.Hour00,
		},
		Festive: CalendarPeriodHours{
			Hour1:  calendarPeriod.Festive.Hour1,
			Hour2:  calendarPeriod.Festive.Hour2,
			Hour3:  calendarPeriod.Festive.Hour3,
			Hour4:  calendarPeriod.Festive.Hour4,
			Hour5:  calendarPeriod.Festive.Hour5,
			Hour6:  calendarPeriod.Festive.Hour6,
			Hour7:  calendarPeriod.Festive.Hour7,
			Hour8:  calendarPeriod.Festive.Hour8,
			Hour9:  calendarPeriod.Festive.Hour9,
			Hour10: calendarPeriod.Festive.Hour10,
			Hour11: calendarPeriod.Festive.Hour11,
			Hour12: calendarPeriod.Festive.Hour12,
			Hour13: calendarPeriod.Festive.Hour13,
			Hour14: calendarPeriod.Festive.Hour14,
			Hour15: calendarPeriod.Festive.Hour15,
			Hour16: calendarPeriod.Festive.Hour16,
			Hour17: calendarPeriod.Festive.Hour17,
			Hour18: calendarPeriod.Festive.Hour18,
			Hour19: calendarPeriod.Festive.Hour19,
			Hour20: calendarPeriod.Festive.Hour20,
			Hour21: calendarPeriod.Festive.Hour21,
			Hour22: calendarPeriod.Festive.Hour22,
			Hour23: calendarPeriod.Festive.Hour23,
			Hour00: calendarPeriod.Festive.Hour00,
		},
		StartDate: calendarPeriod.StartDate,
		EndDate:   calendarPeriod.EndDate,
	}
}

func (repository ProcessMeasureFestiveDays) generateFestiveDayKey(date *time.Time, geographicID *string) string {
	var dateKey, geographicKey string

	if date != nil {
		dateKey = date.Format("2006-01-02")
	}

	if geographicID != nil {
		geographicKey = *geographicID
	}

	return fmt.Sprintf("festive_days:%s:%s", dateKey, geographicKey)
}

func (repository ProcessMeasureFestiveDays) GenerateCalendarPeriodKey(calendarCode, geographicID string) string {
	return fmt.Sprintf("calendars:%s:%s", calendarCode, geographicID)
}

func (repository ProcessMeasureFestiveDays) SaveFestiveDay(ctx context.Context, day measures.FestiveDay) error {
	key := repository.generateFestiveDayKey(&day.Date, &day.GeographicID)

	return repository.RedisCache.Set(ctx, key, day)
}

func (repository ProcessMeasureFestiveDays) SaveCalendarPeriod(ctx context.Context, c []measures.CalendarPeriod) error {
	if len(c) == 0 {
		return errors.New("calendar periods empty")
	}
	key := repository.GenerateCalendarPeriodKey(c[0].CalendarCode, c[0].GeographicID)

	calendarPeriods := utils.MapSlice(c, func(item measures.CalendarPeriod) CalendarPeriod {
		return repository.toDb(item)
	})

	return repository.RedisCache.Set(ctx, key, calendarPeriods)
}

func (repository ProcessMeasureFestiveDays) GetFestiveDay(ctx context.Context, search measures.SearchFestiveDay) (measures.FestiveDay, error) {
	key := repository.generateFestiveDayKey(&search.Date, &search.GeographicID)

	var result measures.FestiveDay

	err := repository.RedisCache.Get(ctx, &result, key)

	return result, err
}

func (repository ProcessMeasureFestiveDays) DeleteFestiveDays(ctx context.Context) error {
	key := repository.generateFestiveDayKey(nil, nil)
	err := repository.RedisCache.Clean(ctx, key)
	return err
}

func (repository ProcessMeasureFestiveDays) GetCalendarPeriod(ctx context.Context, search measures.SearchCalendarPeriod) (measures.CalendarPeriod, error) {
	_, err := repository.GetFestiveDay(ctx, measures.SearchFestiveDay{
		Date:         search.Day.In(search.Location),
		GeographicID: search.GeographicID,
	})

	var isFestive bool

	if err == nil {
		isFestive = true
	}

	key := repository.GenerateCalendarPeriodKey(search.CalendarCode, search.GeographicID)

	result := make([]CalendarPeriod, 0)

	err = repository.RedisCache.Get(ctx, &result, key)

	if err != nil {
		return measures.CalendarPeriod{}, err
	}

	calendarPeriods := make([]CalendarPeriod, 0, cap(result))
	date := time.Date(search.Day.Year(), search.Day.Month(), search.Day.Day(), 0, 0, 0, 0, time.UTC)

	for _, r := range result {
		if utils.InDateRange(r.StartDate, r.EndDate, date) {
			calendarPeriods = append(calendarPeriods, r)
		}
	}

	if len(calendarPeriods) == 0 {
		return measures.CalendarPeriod{}, errors.New("no calendar")
	}

	calendar := calendarPeriods[0].toDomain(search.Day, isFestive)

	return calendar, nil
}

func (repository ProcessMeasureFestiveDays) DeleteCalendars(ctx context.Context) error {
	key := "calendars:*:*:*"

	err := repository.RedisCache.Clean(ctx, key)
	return err
}
