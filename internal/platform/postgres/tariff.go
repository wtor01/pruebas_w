package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
	"bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TariffPostgres struct {
	db         *gorm.DB
	redisCache *redis.DataCacheRedis
	dbTimeout  time.Duration
}

func (t TariffPostgres) GetAllTariffs(ctx context.Context, search tariff.Search) ([]tariff.Tariffs, int, error) {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, t.dbTimeout)
	defer cancel()
	// create list of gz models
	tariffResult := make([]tariff.Tariffs, 0)

	tariffs := make([]Tariff, 0)
	//query for pagination
	paginate := NewPaginate(search.Limit, search.Offset)
	query := t.db.WithContext(ctxTimeout).Preload("Geographic").Preload("Calendar").
		Scopes(
			WithPaginate(paginate),
			WithOrder(validOrderColumnsGeographic, search.Sort),
		)
	result, count := FetchAndCount(query.Find(&tariffs))

	if result.Error != nil {
		return tariffResult, 0, result.Error
	}
	//read results and convert to service model
	for _, ds := range tariffs {
		tariffResult = append(tariffResult, ds.toDomain())
	}
	return tariffResult, int(count), nil
}
func (t TariffPostgres) InsertTariff(ctx context.Context, tf tariff.Tariffs) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, t.dbTimeout)

	var cal Calendars
	resultCal := t.db.WithContext(ctxTimeout).Where("id = ?", tf.CalendarId).Find(&cal)
	if resultCal.Error != nil {
		return resultCal.Error
	}

	defer cancel()
	var gz GeographicZones
	result := t.db.WithContext(ctxTimeout).Where("code = ?", tf.GeographicId).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	if gz.Code == "" {
		return errors.New("Geographic Code not exists")
	}
	err := checkTariffCode(t, tf.Id, ctxTimeout)
	if err != nil {
		return err
	}

	tf.CalendarId = cal.ID
	tf.GeographicId = gz.ID.String()
	var tfs Tariff
	tfs.toDb(tf)
	result = t.db.WithContext(ctxTimeout).Create(&tfs)

	return result.Error
}
func (t TariffPostgres) ModifyTariff(ctx context.Context, tariffId string, tf tariff.Tariffs) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, t.dbTimeout)
	defer cancel()
	var cal Calendars
	resultCal := t.db.WithContext(ctxTimeout).Where("id = ?", tf.CalendarId).Find(&cal)
	if resultCal.Error != nil {
		return resultCal.Error
	}
	var tfs Tariff
	result := t.db.WithContext(ctxTimeout).Where("ID = ?", tariffId).Find(&tfs)
	if result.Error != nil {
		return result.Error
	}
	var gz GeographicZones
	result = t.db.WithContext(ctxTimeout).Where("code = ?", tf.GeographicId).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	tfs.Description = tf.Description
	tfs.CodeOne = tf.CodeOne
	tfs.CodeOdos = tf.CodeOdos
	tfs.TensionLevel = TariffsTensionLevel(tf.TensionLevel)
	tfs.Geographic = gz
	tfs.Periods = tf.Periods
	tfs.Calendar = cal
	tfs.Coef = TariffsCoef(tf.Coef)
	tfs.UpdatedByID = &tf.UpdatedBy

	t.db.Save(&tfs)
	keySearch := fmt.Sprintf("tariff:%v:%s", "get", tariffId)

	if t.redisCache != nil {
		err := t.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}
func (t TariffPostgres) DeleteTariff(ctx context.Context, tariffId string) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, t.dbTimeout)
	defer cancel()
	var contract ContractualSituation
	var exists bool
	err := t.db.WithContext(ctxTimeout).Model(contract).Select("count(*) > 0").
		Where("tariff_id = ?", tariffId).
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Error")
	}

	resultTariffCal := t.db.Unscoped().WithContext(ctxTimeout).Delete(&TariffCalendar{}, "tariff_id = ?", tariffId)
	if resultTariffCal.Error != nil {
		return resultTariffCal.Error
	}

	result := t.db.Unscoped().WithContext(ctxTimeout).Delete(&Tariff{}, "ID = ?", tariffId)
	if result.Error != nil {
		return result.Error
	}
	keySearch := fmt.Sprintf("tariff:%v:%s", "get", tariffId)

	if t.redisCache != nil {
		err := t.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Error
}
func (t TariffPostgres) GetTariff(ctx context.Context, tariffId string) (tariff.Tariffs, error) {
	var tf Tariff
	keyTariffRedis := fmt.Sprintf("tariff:%v:%s", "get", tariffId)
	if t.redisCache != nil {
		err := t.redisCache.Get(ctx, &tf, keyTariffRedis)
		if err == nil {
			return tf.toDomain(), err
		}
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, t.dbTimeout)
	defer cancel()

	result := t.db.WithContext(ctxTimeout).Where("id = ?", tariffId).Preload("Geographic").Preload("Calendar").Find(&tf)
	if result.Error != nil {
		return tariff.Tariffs{}, result.Error
	}
	if t.redisCache != nil {
		err := t.redisCache.Set(ctx, keyTariffRedis, tf)
		if err != nil {
			fmt.Println(err)
		}
	}

	return tf.toDomain(), result.Error
}

func (t TariffPostgres) GetAllTariffCalendar(ctx context.Context, tariffId string, search tariff.Search) ([]tariff.TariffsCalendar, int, error) {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, t.dbTimeout)
	defer cancel()
	// create list of gz models
	tariffResult := make([]tariff.TariffsCalendar, 0)

	tariffs := make([]TariffCalendar, 0)
	//query for pagination
	paginate := NewPaginate(search.Limit, search.Offset)
	query := t.db.WithContext(ctxTimeout).Preload("Geographic").Preload("Calendar").Preload("TariffId").Scopes(
		WithPaginate(paginate),
		WithOrder(validOrderColumnsGeographic, search.Sort),
	)
	result, count := FetchAndCount(query.Where("tariff_id = ?", tariffId).Find(&tariffs))

	if result.Error != nil {
		return tariffResult, 0, result.Error
	}
	return utils.MapSlice(tariffs, func(item TariffCalendar) tariff.TariffsCalendar {
		return item.toDomain()
	}), int(count), nil
}
func (t TariffPostgres) InsertTariffCalendar(ctx context.Context, calendarId string, tariffs tariff.Tariffs, createdBy string) error {
	//generate context timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, t.dbTimeout)
	defer cancel()
	tf := make([]TariffCalendar, 0)
	var cal Calendars
	result := t.db.Where("id = ?", calendarId).Find(&cal)
	if result.Error != nil {
		return result.Error
	}
	var tcal TariffCalendar
	var exists bool
	err := t.db.Model(tcal).
		Select("count(*) > 0").
		Where("calendar_id = ?", cal.ID).
		Where("tariff_id = ?", tariffs.Id).
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	var gz GeographicZones
	t.db.WithContext(ctxTimeout).Where("code = ?", tariffs.GeographicId).Find(&gz)
	if result.Error != nil {
		return result.Error
	}
	if len(tf) != 0 {
		return errors.New("relation exist")
	}
	var tariff Tariff
	result = t.db.Where("id = ?", tariffs.Id).Find(&tariff)
	if result.Error != nil {
		return result.Error
	}
	var tfc TariffCalendar
	tfc.toDB(tariff, gz, createdBy, cal)
	result = t.db.WithContext(ctxTimeout).Create(&tfc)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func NewTariffPostgres(db *gorm.DB, redisCache *redis.DataCacheRedis) *TariffPostgres {
	return &TariffPostgres{db: db, dbTimeout: time.Second * 5, redisCache: redisCache}
}

func (t Tariff) toDomain() tariff.Tariffs {
	return tariff.Tariffs{
		Id:           t.ID,
		CodeOdos:     t.CodeOdos,
		CodeOne:      t.CodeOne,
		Description:  t.Description,
		GeographicId: t.Geographic.Code,
		Periods:      t.Periods,
		TensionLevel: string(t.TensionLevel),
		CalendarId:   t.Calendar.ID,
		Coef:         string(t.Coef),
	}
}
func (t *Tariff) toDb(tf tariff.Tariffs) {
	t.ID = tf.Id
	t.CodeOdos = tf.CodeOdos
	t.CodeOne = tf.CodeOne
	t.Description = tf.Description
	t.GeographicID = &tf.GeographicId
	t.Periods = tf.Periods
	t.TensionLevel = TariffsTensionLevel(tf.TensionLevel)
	t.CalendarID = &tf.CalendarId
	t.Coef = TariffsCoef(tf.Coef)
	t.CreatedByID = tf.CreatedBy
}

func (t *TariffCalendar) toDomain() tariff.TariffsCalendar {
	return tariff.TariffsCalendar{
		CalendarId:     t.Calendar.ID,
		TariffId:       t.TariffId.ID,
		EndDate:        t.EndDate,
		StartDate:      t.StartDate,
		GeoGraphicCode: t.Geographic.Code,
	}
}
func (t *TariffCalendar) toDB(tariffs Tariff, gz GeographicZones, createdBy string, cal Calendars) {
	t.Calendar = cal
	t.TariffId = tariffs
	t.StartDate = time.Now()
	t.CreatedByID = createdBy
	t.Geographic = gz

}

func checkTariffCode(t TariffPostgres, tariffId string, ctxTimeout context.Context) error {
	var tar Tariff
	var exists bool
	err := t.db.Model(tar).Select("count(*) > 0").
		Where("id = ?", tariffId).
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	result := t.db.WithContext(ctxTimeout).Find(&tar)
	if exists {
		return errors.New("tariff code")
	}
	return result.Error
}
