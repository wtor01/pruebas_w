package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type CalendarPostgres struct {
	db         *gorm.DB
	dbTimeout  time.Duration
	redisCache *redis.DataCacheRedis
}

// GetAllCalendars get all callendars from db
func (c CalendarPostgres) GetAllCalendars(ctx context.Context, search calendar.Search) ([]calendar.Calendar, int, error) {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()
	//create list of calendars
	calendarsResult := make([]calendar.Calendar, 0)
	calendars := make([]Calendars, 0)
	//query for pagination
	paginate := NewPaginate(search.Limit, search.Offset)
	query := c.db.WithContext(ctxTimeout).Preload("Geographic").
		Scopes(
			WithPaginate(paginate),
			WithOrder(validOrderColumnsGeographic, search.Sort),
		)
	result, count := FetchAndCount(query.Find(&calendars))
	if result.Error != nil {
		return calendarsResult, 0, result.Error
	}
	//read results and convert to service model
	for _, ds := range calendars {
		calendarsResult = append(calendarsResult, ds.toDomain())
	}
	return calendarsResult, int(count), nil
}
func (c CalendarPostgres) InsertCalendars(ctx context.Context, calendar calendar.Calendar) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()
	//check geographic code
	var gz GeographicZones
	result := c.db.WithContext(ctxTimeout).Where("code = ?", calendar.GeographicCode).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	if gz.Code == "" {
		return errors.New("Geographic Code not exists")
	}
	err := checkCalendarCode(c, calendar.Id, ctxTimeout)
	if err != nil {
		return err
	}
	var cals Calendars
	cals.toDB(calendar, gz)
	result = c.db.WithContext(ctxTimeout).Create(&cals)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (c CalendarPostgres) GetCalendar(ctx context.Context, calendarID string) (calendar.Calendar, error) {
	var cal Calendars
	keyCalendarRedis := fmt.Sprintf("calendar:%v:%s", "get", calendarID)
	if c.redisCache != nil {
		err := c.redisCache.Get(ctx, &cal, keyCalendarRedis)
		if err == nil {
			return cal.toDomain(), err
		}
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	result := c.db.WithContext(ctxTimeout).Where("id = ?", calendarID).Find(&cal)
	if result.Error != nil {
		return cal.toDomain(), result.Error
	}

	var gz GeographicZones
	result = c.db.WithContext(ctxTimeout).Where("id = ?", cal.GeographicID).Find(&gz)
	if result.Error != nil {
		return calendar.Calendar{}, result.Error
	}
	cal.Geographic = gz
	if c.redisCache != nil {
		err := c.redisCache.Set(ctx, keyCalendarRedis, cal)
		if err != nil {
			fmt.Println(err)
		}
	}

	return cal.toDomain(), result.Error
}
func (c CalendarPostgres) ModifyCalendar(ctx context.Context, calendarID string, cal calendar.Calendar) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()
	//check geographic code
	var gz GeographicZones
	result := c.db.WithContext(ctxTimeout).Where("code = ?", cal.GeographicCode).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	if gz.Code == "" {
		return errors.New("Geographic Code not exists")
	}
	// UPDATE
	var cals Calendars
	result = c.db.WithContext(ctxTimeout).Where("id = ?", calendarID).Find(&cals)
	if result.Error != nil {
		return result.Error
	}
	cals.Description = cal.Description
	cals.UpdatedByID = &cal.UpdatedBy
	cals.UpdatedAt = cal.UpdatedAt
	cals.Periods = cal.Periods
	cals.Geographic = gz
	c.db.Save(&cals)
	keySearch := fmt.Sprintf("calendar:%v:%s", "get", calendarID)
	if c.redisCache != nil {
		err := c.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}

	return result.Error
}
func (c CalendarPostgres) DeleteCalendar(ctx context.Context, calendarID string) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()
	result := c.db.WithContext(ctxTimeout).Unscoped().Delete(&Calendars{}, "id = ?", calendarID)
	if result.Error != nil {
		return result.Error
	}
	keySearch := fmt.Sprintf("calendar:%v:%s", "get", calendarID)
	if c.redisCache != nil {
		err := c.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}

// GetPeriodsByDate get all periods froma calendar in a date
func (c CalendarPostgres) GetPeriodsByDate(ctx context.Context, calendarCode string, date time.Time) ([]calendar.PeriodCalendar, error) {
	periodsResult := make([]calendar.PeriodCalendar, 0)
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	periods := make([]CalendarPeriods, 0)
	result := c.db.WithContext(ctxTimeout).Where("calendar_id = ?", calendarCode).Where("start_date <= ?", date).Where("end_date >= ?", date).Find(&periods)
	if result.Error != nil {
		return periodsResult, result.Error
	}
	for _, ds := range periods {
		periodsResult = append(periodsResult, ds.toDomain())
	}

	return periodsResult, result.Error
}
func (c CalendarPostgres) GetAllPeriodCalendars(ctx context.Context, calendarId string, search calendar.Search) ([]calendar.PeriodCalendar, int, error) {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()
	//create list of periods
	periodsResult := make([]calendar.PeriodCalendar, 0)
	periods := make([]CalendarPeriods, 0)
	//query for pagination
	paginate := NewPaginate(search.Limit, search.Offset)
	query := c.db.WithContext(ctxTimeout).
		Scopes(
			WithPaginate(paginate),
			WithOrder(validOrderColumnsGeographic, search.Sort),
		).Where("calendar_id = ?", calendarId)
	result, count := FetchAndCount(query.Find(&periods))
	for _, ds := range periods {
		periodsResult = append(periodsResult, ds.toDomain())
	}
	if result.Error != nil {
		return periodsResult, 0, result.Error
	}
	return periodsResult, int(count), result.Error
}
func (c CalendarPostgres) InsertPeriod(ctx context.Context, calendarId string, period calendar.PeriodCalendar) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()
	var pc CalendarPeriods
	pc.toDB(period, calendarId)
	result := c.db.WithContext(ctxTimeout).Create(&pc)

	return result.Error
}
func (c CalendarPostgres) ModifyPeriodicCalendar(ctx context.Context, code string, period calendar.PeriodCalendar) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()
	//get calendar period
	var cp CalendarPeriods
	result := c.db.WithContext(ctxTimeout).Where("ID = ?", code).Find(&cp)
	// UPDATE
	if result.Error != nil {
		return result.Error
	}
	layout := "02-01-2006"
	sd, _ := time.Parse(layout, period.StartDate)
	ed, _ := time.Parse(layout, period.EndDate)
	cp.Description = period.Description
	cp.StartDate = sd
	cp.EndDate = ed
	cp.Year = period.Year
	cp.PeriodNumber = CalendarPeriodPeriodNumber(period.PeriodNumber)
	cp.StartHour = period.StartHour
	cp.EndHour = period.EndHour
	cp.DayType = CalendarPeriodDayType(period.DayType)
	cp.UpdatedByID = &period.UpdatedBy
	cp.UpdatedAt = period.UpdatedAt
	cp.Energy = period.Energy
	cp.Power = period.Power
	result = c.db.Save(&cp)
	keySearch := fmt.Sprintf("period:%v:%s", "get", code)
	if c.redisCache != nil {
		err := c.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}
func (c CalendarPostgres) DeletePeriod(ctx context.Context, code string) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()
	result := c.db.WithContext(ctxTimeout).Unscoped().Delete(&CalendarPeriods{}, "id = ?", code)
	if result.Error != nil {
		return result.Error
	}
	keyCalendarRedis := fmt.Sprintf("period:%v:%s", "get", code)
	if c.redisCache != nil {
		err := c.redisCache.Clean(ctx, keyCalendarRedis)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}
func (c CalendarPostgres) GetPeriodById(ctx context.Context, periodId uuid.UUID) (calendar.PeriodCalendar, error) {
	var cal CalendarPeriods
	keyCalendarRedis := fmt.Sprintf("period:%v:%s", "get", periodId.String())
	if c.redisCache != nil {
		err := c.redisCache.Get(ctx, &cal, keyCalendarRedis)
		if err == nil {
			return cal.toDomain(), err
		}
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	result := c.db.WithContext(ctxTimeout).Where("id = ?", periodId).Find(&cal)
	if result.Error != nil {
		return cal.toDomain(), result.Error
	}
	if c.redisCache != nil {
		err := c.redisCache.Set(ctx, keyCalendarRedis, cal)
		if err != nil {
			fmt.Println(err)
		}
	}
	return cal.toDomain(), result.Error
}
func (c CalendarPostgres) GetPeriodsActive(ctx context.Context) ([]calendar.PeriodCalendar, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	var calendarPeriods []CalendarPeriods
	tx := c.db.WithContext(ctxTimeout).Preload("Calendar.Geographic").Where("energy is true").Order("end_date ASC, start_date DESC").Find(&calendarPeriods)

	if tx.Error != nil {
		return []calendar.PeriodCalendar{}, tx.Error
	}

	return utils.MapSlice(calendarPeriods, func(item CalendarPeriods) calendar.PeriodCalendar {
		return item.toDomain()
	}), nil
}
func NewCalendarPostgres(db *gorm.DB, redisCache *redis.DataCacheRedis) *CalendarPostgres {
	return &CalendarPostgres{db: db, dbTimeout: time.Second * 5, redisCache: redisCache}
}

func (c *Calendars) toDB(calendar calendar.Calendar, geographic GeographicZones) {
	c.ID = calendar.Id
	c.Description = calendar.Description
	c.Periods = calendar.Periods
	c.Geographic = geographic
	c.CreatedByID = calendar.CreatedBy
}
func (c Calendars) toDomain() calendar.Calendar {
	return calendar.Calendar{
		Id:             c.ID,
		Description:    c.Description,
		CreatedAt:      c.CreatedAt,
		GeographicCode: c.Geographic.Code,
		UpdatedAt:      c.UpdatedAt,
		Periods:        c.Periods,
	}
}
func (cp *CalendarPeriods) toDB(periods calendar.PeriodCalendar, calendarId string) {
	layout := "02-01-2006"
	sd, _ := time.Parse(layout, periods.StartDate)
	ed, _ := time.Parse(layout, periods.EndDate)
	cp.CalendarID = &calendarId
	cp.Description = periods.Description
	cp.DayType = CalendarPeriodDayType(periods.DayType)
	cp.PeriodNumber = CalendarPeriodPeriodNumber(periods.PeriodNumber)
	cp.Year = periods.Year
	cp.StartHour = periods.StartHour
	cp.EndHour = periods.EndHour
	cp.StartDate = sd
	cp.EndDate = ed
	cp.CreatedByID = periods.CreatedBy
	cp.Energy = periods.Energy
	cp.Power = periods.Power
}
func (cp *CalendarPeriods) toDomain() calendar.PeriodCalendar {
	layout := "02-01-2006"
	sd := cp.StartDate.Format(layout)
	ed := cp.EndDate.Format(layout)
	var calendarCode, geographicCode string
	if cp.CalendarID != nil {
		calendarCode = *cp.CalendarID
	}
	if cp.Calendar.Geographic.Code != "" {
		geographicCode = cp.Calendar.Geographic.Code
	}

	return calendar.PeriodCalendar{
		ID:             cp.ID,
		CalendarCode:   calendarCode,
		GeographicCode: geographicCode,
		Description:    cp.Description,
		DayType:        calendar.DayType(cp.DayType),
		PeriodNumber:   string(cp.PeriodNumber),
		Year:           cp.Year,
		StartHour:      cp.StartHour,
		EndHour:        cp.EndHour,
		StartDate:      sd,
		EndDate:        ed,
		Energy:         cp.Energy,
		Power:          cp.Power,
	}
}

func checkCalendarCode(c CalendarPostgres, id string, ctxTimeout context.Context) error {
	var calendar Calendars
	var exists bool
	err := c.db.Model(calendar).Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	result := c.db.WithContext(ctxTimeout).Find(&calendar)
	if exists {
		return errors.New("calendar code")
	}
	return result.Error
}
