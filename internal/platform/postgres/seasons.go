package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"context"
	"fmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type SeasonPostgres struct {
	db         *gorm.DB
	dbTimeout  time.Duration
	redisCache *redis.DataCacheRedis
}

func NewSeasonPostgres(db *gorm.DB, redisCache *redis.DataCacheRedis) *SeasonPostgres {
	return &SeasonPostgres{db: db, dbTimeout: time.Second * 5, redisCache: redisCache}
}

func (d SeasonPostgres) GetAllSeasons(ctx context.Context, search seasons.Search) ([]seasons.Seasons, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	seasonsResult := make([]seasons.Seasons, 0)

	seasons := make([]Seasons, 0)
	paginate := NewPaginate(search.Limit, search.Offset)
	query := d.db.WithContext(ctxTimeout).Preload("Geographic").
		Scopes(
			WithPaginate(paginate),
		)
	result, count := FetchAndCount(query.Find(&seasons))
	if result.Error != nil {
		return seasonsResult, 0, result.Error
	}
	//read results and convert to service model
	for _, ds := range seasons {

		seasonsResult = append(seasonsResult, ds.toDomain())
	}
	return seasonsResult, int(count), nil
}
func (d SeasonPostgres) GetSeason(ctx context.Context, seasonId string) (seasons.Seasons, error) {
	var sea Seasons
	keySeasonRedis := fmt.Sprintf("season:%v:%s", "get", seasonId)
	if d.redisCache != nil {
		err := d.redisCache.Get(ctx, &sea, keySeasonRedis)
		if err == nil {
			return sea.toDomain(), err
		}
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	result := d.db.WithContext(ctxTimeout).Where("id = ?", seasonId).Find(&sea)
	if result.Error != nil {

		return sea.toDomain(), result.Error
	}

	var gz GeographicZones
	result = d.db.WithContext(ctxTimeout).Where("id = ?", sea.GeographicID).Find(&gz)
	if result.Error != nil {
		return seasons.Seasons{}, result.Error
	}
	sea.Geographic = gz
	if d.redisCache != nil {
		err := d.redisCache.Set(ctx, keySeasonRedis, sea)
		if err != nil {
			fmt.Println(err)
		}
	}

	return sea.toDomain(), result.Error
}
func (d SeasonPostgres) InsertSeason(ctx context.Context, season seasons.Seasons) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	var gz GeographicZones
	result := d.db.WithContext(ctxTimeout).Where("code = ?", season.GeographicCode).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	if gz.Code == "" {
		return errors.New("Geographic Code not exists")
	}
	var sea Seasons
	result = d.db.WithContext(ctxTimeout).Where("name = ?", season.Name).Find(&sea)
	if result.Error != nil {
		return result.Error
	}
	if sea.Name != "" {
		return errors.New("Season exists")
	}
	var s Seasons
	s.toDB(season, gz)
	result = d.db.WithContext(ctxTimeout).Create(&s)
	return result.Error
}
func (d SeasonPostgres) ModifySeason(ctx context.Context, seasonId uuid.UUID, season seasons.Seasons) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var gz GeographicZones
	result := d.db.WithContext(ctxTimeout).Where("code = ?", season.GeographicCode).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	if gz.Code == "" {
		return errors.New("Geographic Code not exists")
	}

	var seas Seasons
	result = d.db.Where("code = ")
	if seasonId != season.ID {
		err := checkSeasonId(d, season.ID, ctxTimeout)
		if err != nil {
			return err
		}
	}
	result = d.db.WithContext(ctxTimeout).Where("id = ?", seasonId).Find(&seas)
	if result.Error != nil {
		return result.Error
	}
	seas.Name = season.Name
	seas.Description = season.Description
	seas.Geographic = gz
	seas.UpdatedByID = &season.UpdatedBy
	seas.UpdatedAt = season.UpdatedAt
	d.db.Save(&seas)
	keySearch := fmt.Sprintf("season:%v:%s", "get", seasonId)

	if d.redisCache != nil {
		err := d.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}
func (d SeasonPostgres) DeleteSeason(ctx context.Context, seasonId uuid.UUID) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	result := d.db.WithContext(ctxTimeout).Unscoped().Delete(&Seasons{}, "id = ?", seasonId)
	if result.Error != nil {
		return result.Error
	}
	keySearch := fmt.Sprintf("season:%v:%s", "get", seasonId)
	if d.redisCache != nil {
		err := d.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}

func (d SeasonPostgres) GetAllDayTypes(ctx context.Context, seasonId string, search seasons.Search) ([]seasons.DayTypes, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	dayTypesResult := make([]seasons.DayTypes, 0)
	dayTypes := make([]DayTypes, 0)
	paginate := NewPaginate(search.Limit, search.Offset)
	query := d.db.WithContext(ctxTimeout).
		Scopes(
			WithPaginate(paginate),
		).Where("seasons_id = ?", seasonId)
	result, count := FetchAndCount(query.Find(&dayTypes))
	for _, ds := range dayTypes {
		dayTypesResult = append(dayTypesResult, ds.toDomain())
	}
	if result.Error != nil {
		return dayTypesResult, 0, result.Error
	}
	return dayTypesResult, int(count), result.Error
}
func (d SeasonPostgres) GetDayType(ctx context.Context, dayTypeId string) (seasons.DayTypes, error) {
	var dt DayTypes
	keydtRedis := fmt.Sprintf("daytype:%v:%s", "get", dayTypeId)
	if d.redisCache != nil {
		err := d.redisCache.Get(ctx, &dt, keydtRedis)
		if err == nil {
			return dt.toDomain(), err
		}
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	result := d.db.WithContext(ctxTimeout).Where("id = ?", dayTypeId).Find(&dt)
	if result.Error != nil {
		return dt.toDomain(), result.Error
	}
	if d.redisCache != nil {
		err := d.redisCache.Set(ctx, keydtRedis, dt)
		if err != nil {
			fmt.Println(err)
		}
	}
	return dt.toDomain(), result.Error
}
func (d SeasonPostgres) InsertDayTypes(ctx context.Context, seasonId string, dayT seasons.DayTypes) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	var dt DayTypes
	result := d.db.WithContext(ctxTimeout).Where("name = ?", dayT.Name).Find(&dt)
	if result.Error != nil {
		errors.New("Day type exist")
		return result.Error
	}
	dt.toDB(dayT, seasonId)
	result = d.db.WithContext(ctxTimeout).Create(&dt)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (d SeasonPostgres) ModifyDayTypes(ctx context.Context, dayTypeId uuid.UUID, dayT seasons.DayTypes) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	//get calendar period
	var dt DayTypes
	if dayTypeId.String() != dayT.ID {
		err := checkDayTypeId(d, uuid.FromStringOrNil(dayT.ID), ctxTimeout)
		if err != nil {
			return err
		}
	}
	result := d.db.WithContext(ctxTimeout).Where("id = ?", dayTypeId).Find(&dt)
	if result.Error != nil {
		return result.Error
	}

	dt.Name = dayT.Name
	dt.Month = dayT.Month
	dt.IsFestive = dayT.IsFestive
	dt.UpdatedAt = dayT.UpdatedAt
	dt.UpdatedByID = &dayT.UpdatedBy
	result = d.db.WithContext(ctxTimeout).Save(&dt)
	keySearch := fmt.Sprintf("daytype:%v:%s", "get", dayTypeId)
	if d.redisCache != nil {
		err := d.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error

}
func (d SeasonPostgres) DeleteDayType(ctx context.Context, dayTypeId uuid.UUID) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	result := d.db.WithContext(ctxTimeout).Unscoped().Delete(&DayTypes{}, "id = ?", dayTypeId)
	if result.Error != nil {
		return result.Error
	}
	keySearch := fmt.Sprintf("daytype:%v:%s", "get", dayTypeId)
	if d.redisCache != nil {
		err := d.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}
func (d SeasonPostgres) GetDayTypeByMonth(ctx context.Context, month int, isFestive bool) (seasons.DayTypes, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	var dayType DayTypes

	tx := d.db.Unscoped().WithContext(ctxTimeout).Where("month = ?", month).Where("is_festive = ?", isFestive).First(&dayType)
	if tx.Error != nil {
		return dayType.toDomain(), tx.Error
	}

	return dayType.toDomain(), tx.Error
}

func (s *Seasons) toDB(season seasons.Seasons, geographic GeographicZones) {
	s.Name = season.Name
	s.Description = season.Description
	s.Geographic = geographic
	s.CreatedByID = season.CreatedBy

}
func (s Seasons) toDomain() seasons.Seasons {
	return seasons.Seasons{
		ID:             s.ID,
		Name:           s.Name,
		Description:    s.Description,
		GeographicCode: s.Geographic.Code,
		CreatedAt:      s.CreatedAt,
		CreatedBy:      s.CreatedByID,
		UpdatedAt:      s.UpdatedAt,
	}
}

func (dt *DayTypes) toDB(dayT seasons.DayTypes, seasonId string) {
	dt.Name = dayT.Name
	dt.Month = dayT.Month
	dt.IsFestive = dayT.IsFestive
	dt.SeasonsID = &seasonId
	dt.CreatedByID = dayT.CreatedBy
}
func (dt *DayTypes) toDomain() seasons.DayTypes {
	s := seasons.DayTypes{
		ID:        dt.ID.String(),
		Name:      dt.Name,
		Month:     dt.Month,
		IsFestive: dt.IsFestive,
	}

	if dt.SeasonsID != nil {
		s.SeasonsId = *dt.SeasonsID
	}

	return s
}

func checkSeasonId(d SeasonPostgres, id uuid.UUID, ctxTimeout context.Context) error {
	var season Seasons
	var exists bool
	err := d.db.Model(season).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).Error
	if err != nil {
		return err
	}
	result := d.db.WithContext(ctxTimeout).Find(&season)
	if exists {
		return errors.New("id is already created")
	}
	return result.Error
}

func checkDayTypeId(d SeasonPostgres, id uuid.UUID, ctxTimeout context.Context) error {
	var dayT DayTypes
	var exists bool
	err := d.db.Model(dayT).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).Error
	if err != nil {
		return err
	}
	result := d.db.WithContext(ctxTimeout).Find(&dayT)
	if exists {
		return errors.New("id is already created")
	}
	return result.Error
}
