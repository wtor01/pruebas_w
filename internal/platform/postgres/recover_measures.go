package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type RecoverMeasuresRepository struct {
	db        *gorm.DB
	dbTimeout time.Duration
}

func NewRecoverMeasuresPostgres(db *gorm.DB) *RecoverMeasuresRepository {
	return &RecoverMeasuresRepository{db: db, dbTimeout: time.Second * 5}
}

type RecoverMeasures struct {
	ModelEntity

	MeterId           string    `gorm:"meter_id" gorm:"index:idx_name,unique"`
	Date              time.Time `gorm:"date" gorm:"index:idx_name,unique"`
	DistributorID     string    `gorm:"distributor_id" gorm:"index:idx_name,unique"`
	DistributorCode   string    `gorm:"distributor_code" gorm:"index:idx_name,unique"`
	CUPS              string    `gorm:"cups" gorm:"index:idx_name,unique"`
	MeterSerialNumber string    `gorm:"meter_serial_number" gorm:"index:idx_name,unique"`
	MeasurePointType  string    `gorm:"measure_point_type" gorm:"index:idx_name,unique"`
	ContractNumber    string    `gorm:"contract_number" gorm:"index:idx_name,unique"`
	MeterConfigId     string    `gorm:"meter_config_id" gorm:"index:idx_name,unique"`
	ReadingType       string    `gorm:"reading_type" gorm:"index:idx_name,unique"`

	ServiceType string `gorm:"service_type" gorm:"index:idx_name,unique"`
	PointType   string `gorm:"point_type" gorm:"index:idx_name,unique"`
}

func (r RecoverMeasures) toDomain() process_measures.RecoverMeasures {
	return process_measures.RecoverMeasures{
		ID:                r.ModelEntity.ID.String(),
		MeterId:           r.MeterId,
		Date:              r.Date,
		DistributorID:     r.DistributorID,
		DistributorCode:   r.DistributorCode,
		CUPS:              r.CUPS,
		MeterSerialNumber: r.MeterSerialNumber,
		MeasurePointType:  r.MeasurePointType,
		ContractNumber:    r.ContractNumber,
		MeterConfigId:     r.MeterConfigId,
		ReadingType:       r.ReadingType,
		ServiceType:       r.ServiceType,
		PointType:         r.PointType,
	}
}

func (r RecoverMeasuresRepository) toDb(s process_measures.RecoverMeasures) RecoverMeasures {

	return RecoverMeasures{
		ModelEntity: ModelEntity{
			ID: uuid.FromStringOrNil(s.ID),
		},
		MeterId:           s.MeterId,
		Date:              s.Date,
		DistributorID:     s.DistributorID,
		DistributorCode:   s.DistributorCode,
		CUPS:              s.CUPS,
		MeterSerialNumber: s.MeterSerialNumber,
		MeasurePointType:  s.MeasurePointType,
		ContractNumber:    s.ContractNumber,
		MeterConfigId:     s.MeterConfigId,
		ReadingType:       s.ReadingType,
		ServiceType:       s.ServiceType,
		PointType:         s.PointType,
	}
}

func (r RecoverMeasuresRepository) AllRecoverMeasures(ctx context.Context, search process_measures.PaginationRecoverMeasures) ([]process_measures.RecoverMeasures, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	recoverMeasuresResult := make([]process_measures.RecoverMeasures, 0, search.Limit)

	recoverMeasures := make([]RecoverMeasures, 0)

	paginate := NewPaginate(search.Limit, search.Offset)
	query := r.db.WithContext(ctxTimeout).
		Scopes(
			WithPaginate(paginate),
		)

	result, count := FetchAndCount(query.Find(&recoverMeasures))

	if result.Error != nil {
		return recoverMeasuresResult, 0, result.Error
	}

	for _, rm := range recoverMeasures {
		recoverMeasuresResult = append(recoverMeasuresResult, rm.toDomain())
	}

	return recoverMeasuresResult, int(count), nil
}

func (r RecoverMeasuresRepository) SaveRecoverMeasures(ctx context.Context, s process_measures.RecoverMeasures) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	entity := r.toDb(s)

	tx := r.db.WithContext(ctxTimeout).Save(&entity)

	return tx.Error
}

func (r RecoverMeasuresRepository) SearchRecoverMeasures(ctx context.Context, search process_measures.RecoverMeasures) (process_measures.RecoverMeasures, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var result []RecoverMeasures

	q := r.db.WithContext(ctxTimeout).Debug().Where("date = ?", search.Date)
	q.Where("distributor_id = ?", search.DistributorID)
	q.Where("cups = ?", search.CUPS)
	q.Where("meter_serial_number = ?", search.MeterSerialNumber)
	q.Where("measure_point_type = ?", search.MeasurePointType)
	q.Where("contract_number = ?", search.ContractNumber)
	q.Where("meter_config_id = ?", search.MeterConfigId)
	q.Where("service_type = ?", search.ServiceType)
	q.Where("point_type = ?", search.PointType)
	q.Where("reading_type = ?", search.ReadingType)

	tx := q.Limit(1).Find(&result)

	if tx.Error != nil {
		return process_measures.RecoverMeasures{}, 0, tx.Error
	}

	if len(result) == 0 {
		return process_measures.RecoverMeasures{}, 0, nil
	}
	return result[0].toDomain(), len(result), nil
}

func (r RecoverMeasuresRepository) GetRecoverMeasures(ctx context.Context, id string) (process_measures.RecoverMeasures, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity RecoverMeasures

	tx := r.db.WithContext(ctxTimeout).First(&entity, "id = ?", id)

	if tx.Error != nil {
		return process_measures.RecoverMeasures{}, tx.Error
	}

	return entity.toDomain(), nil
}

func (r RecoverMeasuresRepository) DeleteRecoverMeasures(ctx context.Context, id string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	var entity RecoverMeasures

	tx := r.db.WithContext(ctxTimeout).Delete(&entity, "id = ? ", id)

	return tx.Error

}
