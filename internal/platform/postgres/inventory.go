package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

type InventoryPostgres struct {
	db         *gorm.DB
	redisCache *redis.DataCacheRedis
	dbTimeout  time.Duration
}

func servicePointConfigViewToDomain2(s ServicePointConfigView) measures.MeterConfig {
	return measures.MeterConfig{
		ID:               s.MeterConfigId,
		StartDate:        time.Time{},
		EndDate:          time.Time{},
		CurveType:        measures.RegisterType(s.CurveType),
		ReadingType:      s.ReadingType,
		PriorityContract: measures.PriorityContract(s.PriorityContract),
		CalendarID:       s.CalendarCode,
		AI:               s.AI,
		AE:               s.AE,
		R1:               s.R1,
		R2:               s.R2,
		R3:               s.R3,
		R4:               s.R4,
		M:                0,
		E:                0,
		RentingPrice:     0,
		DistributorID:    s.DistributorId,
		DistributorCode:  s.DistributorCode,
		TlgCode:          measures.MeterConfigTlgCode(s.TlgCode),
		Type:             measures.MeterType(s.TlgType),
		Meter: measures.Meter{
			ID:                s.MeterId,
			SerialNumber:      s.MeterSerialNumber,
			ActiveConstant:    1,
			ReactiveConstant:  1,
			MaximeterConstant: 1,
		},
		ServicePoint: measures.ServicePoint{
			Cups:      s.Cups,
			StartDate: time.Time{},
			PointType: measures.PointType(s.PointType),
		},
		MeasurePoint: measures.MeasurePoint{
			Type: measures.MeasurePointType(s.MeasurePointType),
		},
		ContractualSituations: measures.ContractualSituations{
			ID:           s.ContractualSituationId,
			TariffID:     s.TariffID,
			P1Demand:     s.P1Demand,
			P2Demand:     s.P2Demand,
			P3Demand:     s.P3Demand,
			P4Demand:     s.P4Demand,
			P5Demand:     s.P5Demand,
			P6Demand:     s.P6Demand,
			RetailerCode: 0,
			InitDate:     time.Time{},
			EndDate:      time.Time{},
			RetailerCdos: s.Coefficient,
			RetailerName: s.GeographicID,
		},
	}
}

func (d InventoryPostgres) ListMeterConfigByDate(ctx context.Context, query measures.ListMeterConfigByDateQuery) ([]measures.MeterConfig, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var a []ServicePointConfigView
	dateFormat := query.Date.UTC().Format("02/01/2006")

	tx := d.db.WithContext(ctxTimeout).Debug().Where("distributor_id = ?", query.DistributorID).Where("meters_type IN (?)", query.MeterType).Where("service_type = ?", query.ServiceType).Where("point_type = ?", query.PointType).Where("start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("end_date IS NULL or end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_end_date IS NULL or contractual_situation_end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat)

	if query.Limit > 0 {
		tx = tx.Limit(query.Limit).Offset(query.Offset)
	}

	tx.Order("meter_id")

	tx.Find(&a)

	if tx.Error != nil {
		return []measures.MeterConfig{}, tx.Error
	}

	return utils.MapSlice(a, servicePointConfigViewToDomain2), nil
}

func (d InventoryPostgres) CountMeterConfigByDate(ctx context.Context, query measures.ListMeterConfigByDateQuery) (int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var count int64
	dateFormat := query.Date.UTC().Format("02/01/2006")

	tx := d.db.WithContext(ctxTimeout).Model(&ServicePointConfigView{}).Debug().Where("distributor_id = ?", query.DistributorID).Where("meters_type IN (?)", query.MeterType).Where("service_type = ?", query.ServiceType).Where("point_type = ?", query.PointType).Where("start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("end_date IS NULL or end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_end_date IS NULL or contractual_situation_end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat).Count(&count)

	if tx.Error != nil {
		return 0, tx.Error
	}

	return int(count), nil
}

func (d InventoryPostgres) GetMeterConfigByMeter(ctx context.Context, query measures.GetMeterConfigByMeterQuery) (measures.MeterConfig, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var result []ServicePointConfigView
	dateFormat := query.Date.Format("02/01/2006")

	q := d.db.WithContext(ctxTimeout).Debug().Where("start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("end_date IS NULL or end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_start_date <= to_timestamp(?,'DD/MM/YYYY')", dateFormat).Where("contractual_situation_end_date IS NULL or contractual_situation_end_date > to_timestamp(?,'DD/MM/YYYY')", dateFormat)

	switch {
	case query.MeterSerialNumber != "":
		q.Where("meter_serial_number = ?", query.MeterSerialNumber)

	default:
		return measures.MeterConfig{}, nil
	}

	tx := q.Limit(1).Find(&result)

	if tx.Error != nil {
		return measures.MeterConfig{}, tx.Error
	}

	if len(result) == 0 {
		return measures.MeterConfig{}, nil
	}

	return servicePointConfigViewToDomain2(result[0]), nil
}

func NewInventoryPostgres(db *gorm.DB, redis *redis.DataCacheRedis) *InventoryPostgres {
	return &InventoryPostgres{db: db, dbTimeout: time.Second * 5, redisCache: redis}
}
