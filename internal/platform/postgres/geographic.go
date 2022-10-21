package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type GeographicPostgres struct {
	db         *gorm.DB
	redisCache *redis.DataCacheRedis
	dbTimeout  time.Duration
}

func NewGeographicPostgres(db *gorm.DB, redisCache *redis.DataCacheRedis) *GeographicPostgres {
	return &GeographicPostgres{db: db, dbTimeout: time.Second * 5, redisCache: redisCache}
}

// var used to order columns used for get all
var validOrderColumnsGeographic = map[string]struct{}{
	"code":       struct{}{},
	"id":         struct{}{},
	"created_at": struct{}{},
	"created_by": struct{}{},
}

// toDomain convert postgres struct to service struct
func (d GeographicZones) toDomain() geographic.GeographicZones {
	return geographic.GeographicZones{
		ID:          d.ID,
		Code:        d.Code,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		CreatedBy:   d.CreatedByID,
		UpdatedAt:   d.UpdatedAt,
	}
}

// GetAllGeographicZones get all gz from gorm models
func (d GeographicPostgres) GetAllGeographicZones(ctx context.Context, search geographic.Search) ([]geographic.GeographicZones, int, error) {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	// create list of gz models
	geographicsResult := make([]geographic.GeographicZones, 0)

	geographics := make([]GeographicZones, 0)
	//query for pagination
	paginate := NewPaginate(search.Limit, search.Offset)
	query := d.db.WithContext(ctxTimeout).
		Scopes(
			WithPaginate(paginate),
			WithOrder(validOrderColumnsGeographic, search.Sort),
		)
	result, count := FetchAndCount(query.Find(&geographics))

	if result.Error != nil {
		return geographicsResult, 0, result.Error
	}
	//read results and convert to service model
	for _, ds := range geographics {
		geographicsResult = append(geographicsResult, ds.toDomain())
	}
	return geographicsResult, int(count), nil
}

// InsertGeographicZone insert gz in gorm model
func (d GeographicPostgres) InsertGeographicZone(ctx context.Context, zones geographic.GeographicZones) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	if zones.Code == "" {
		return errors.New("Code is empty")
	}
	err := checkGeographicCode(d, zones.Code, ctxTimeout)
	if err != nil {
		return err
	}
	var gz GeographicZones
	gz.toDB(zones)
	result := d.db.WithContext(ctxTimeout).Create(&gz)
	return result.Error
}

// GetGeographicZone get only one gz
func (d GeographicPostgres) GetGeographicZone(ctx context.Context, geographicId string) (geographic.GeographicZones, error) {
	var gz GeographicZones

	keyGeographicRedis := fmt.Sprintf("geographic:%v:%s", "get", geographicId)
	if d.redisCache != nil {
		err := d.redisCache.Get(ctx, &gz, keyGeographicRedis)
		if err == nil {
			return gz.toDomain(), err
		}
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	result := d.db.Where("id = ?", geographicId).WithContext(ctxTimeout).Find(&gz)

	if d.redisCache != nil {
		err := d.redisCache.Set(ctx, keyGeographicRedis, gz)
		if err != nil {
			fmt.Println(err)
		}
	}
	return gz.toDomain(), result.Error

}

// ModifyGeographicZone modify gorm gz
func (d GeographicPostgres) ModifyGeographicZone(ctx context.Context, geographicId string, zones geographic.GeographicZones) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	var gz GeographicZones

	result := d.db.Where("id = ?", geographicId).WithContext(ctxTimeout).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	if gz.Code != zones.Code {
		return errors.New("Code cannot be modified")
	}

	gz.Description = zones.Description
	gz.UpdatedByID = &zones.UpdatedBy
	gz.UpdatedAt = zones.UpdatedAt
	result = d.db.Save(&gz)

	keySearch := fmt.Sprintf("geographic:%v:%s", "get", geographicId)
	if d.redisCache != nil {
		err := d.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}

	return result.Error

}

// DeleteGeographicZone delete gorm model
func (d GeographicPostgres) DeleteGeographicZone(ctx context.Context, geographicId string) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()
	result := d.db.WithContext(ctxTimeout).Unscoped().Delete(&GeographicZones{}, "id = ?", geographicId)
	if result.Error != nil {
		return result.Error
	}
	keySearch := fmt.Sprintf("geographic:%v:%s", "get", geographicId)
	if d.redisCache != nil {
		err := d.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}

// toDB convert gorm model to service model
func (d *GeographicZones) toDB(zones geographic.GeographicZones) {
	d.Code = zones.Code
	d.Description = zones.Description
	d.CreatedByID = zones.CreatedBy
}

// checkGeographicCode check if code alredy exists
func checkGeographicCode(d GeographicPostgres, code string, ctxTimeout context.Context) error {
	var geographics GeographicZones
	var exists bool
	err := d.db.Model(geographics).
		Select("count(*) > 0").
		Where("code = ?", code).
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	result := d.db.WithContext(ctxTimeout).Find(&geographics)
	if exists {
		return errors.New("code is alredy created")
	}
	return result.Error
}
