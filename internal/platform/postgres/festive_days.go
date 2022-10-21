package postgres

import (
	master_tables "bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func (f FestiveDays) toDomain() master_tables.FestiveDay {
	return master_tables.FestiveDay{
		Id:           f.ID,
		Date:         f.Date,
		Description:  f.Description,
		GeographicId: f.Geographic.Code,
	}
}
func (f *FestiveDays) toDB(festiveDay master_tables.FestiveDay, gz GeographicZones) {
	f.Date = festiveDay.Date
	f.Description = festiveDay.Description
	f.Geographic = gz
	f.CreatedByID = festiveDay.CreatedBy
}

type FestiveDaysRepository struct {
	db         *gorm.DB
	redisCache *redis.DataCacheRedis
	dbTimeout  time.Duration
}

func NewFestiveDaysPostgres(db *gorm.DB, redisCache *redis.DataCacheRedis) *FestiveDaysRepository {
	return &FestiveDaysRepository{db: db, dbTimeout: time.Second * 5, redisCache: redisCache}
}

var validOrderColumnsFestiveDays = map[string]struct{}{
	"date":          struct{}{},
	"description":   struct{}{},
	"geographic_id": struct{}{},
}

func (r FestiveDaysRepository) GetListFestiveDays(ctx context.Context, search master_tables.Search) ([]master_tables.FestiveDay, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	festiveDays := make([]FestiveDays, 0)
	paginate := NewPaginate(search.Limit, search.Offset)
	query := r.db.WithContext(ctxTimeout).Preload("Geographic").
		Scopes(
			WithPaginate(paginate),
			WithOrder(validOrderColumnsFestiveDays, search.Sort),
		)

	result, count := FetchAndCount(query.Find(&festiveDays))

	if result.Error != nil {
		return []master_tables.FestiveDay{}, 0, result.Error
	}

	return utils.MapSlice(festiveDays, func(item FestiveDays) master_tables.FestiveDay {

		return item.toDomain()
	}), int(count), result.Error

}

func (r FestiveDaysRepository) GetFestiveDayById(ctx context.Context, festiveDayId string) (master_tables.FestiveDay, error) {
	var festiveDay FestiveDays
	keyFDRedis := fmt.Sprintf("festiveday:%v:%s", "get", festiveDayId)
	if r.redisCache != nil {
		err := r.redisCache.Get(ctx, &festiveDay, keyFDRedis)
		if err == nil {
			return festiveDay.toDomain(), err
		}
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	tx := r.db.WithContext(ctxTimeout).First(&festiveDay, "id = ?", festiveDayId)

	if tx.Error != nil {
		return master_tables.FestiveDay{}, tx.Error
	}

	var gz GeographicZones
	tx = r.db.Where("id = ?", festiveDay.GeographicID).WithContext(ctxTimeout).Find(&gz)
	if tx.Error != nil {
		return master_tables.FestiveDay{}, tx.Error
	}
	festiveDay.Geographic = gz
	if r.redisCache != nil {
		err := r.redisCache.Set(ctx, keyFDRedis, festiveDay)
		if err != nil {
			fmt.Println(err)
		}
	}
	return festiveDay.toDomain(), tx.Error
}

func (r FestiveDaysRepository) SaveFestiveDay(ctx context.Context, festiveDay master_tables.FestiveDay) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	var gz GeographicZones
	result := r.db.Where("code = ?", festiveDay.GeographicId).WithContext(ctxTimeout).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	var newFestiveDay FestiveDays
	newFestiveDay.toDB(festiveDay, gz)
	tx := r.db.WithContext(ctxTimeout).Create(&newFestiveDay)

	return tx.Error
}

func (r FestiveDaysRepository) UpdateFestiveDay(ctx context.Context, festiveDayId string, festiveDay master_tables.FestiveDay) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var fd FestiveDays

	tx := r.db.Where("id = ?", festiveDayId).WithContext(ctxTimeout).Find(&fd)

	if tx.Error != nil {
		return tx.Error
	}
	fd.Date = festiveDay.Date
	fd.Description = festiveDay.Description
	result := r.db.Where("code = ?", festiveDay.GeographicId).WithContext(ctxTimeout).Find(&fd.Geographic)
	if result.Error != nil {
		return nil
	}

	tx = r.db.WithContext(ctxTimeout).Save(&fd)

	if tx.Error != nil {
		return tx.Error
	}
	keySearch := fmt.Sprintf("festiveday:%v:%s", "get", festiveDayId)
	if r.redisCache != nil {
		err := r.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}

	return tx.Error
}

func (r FestiveDaysRepository) DeleteFestiveDay(ctx context.Context, festiveDayId string) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	tx := r.db.WithContext(ctxTimeout).Unscoped().Delete(&FestiveDays{}, "ID = ?", festiveDayId)
	keySearch := fmt.Sprintf("festiveday:%v:%s", "get", festiveDayId)
	if r.redisCache != nil {
		err := r.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return tx.Error

}
