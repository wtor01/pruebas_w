package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type BillingMeasuresScheduler struct {
	ModelEntity
	DistributorID *string        `gorm:"type:uuid;"`
	Distributor   Distributor    `gorm:"distributor" gorm:"foreignKey:DistributorID"`
	ServiceType   string         `gorm:"column:service_type;not null;type:service_point_service_type"`
	PointType     string         `gorm:"column:point_type;not null;type:service_point_point_type;"`
	MeterType     pq.StringArray `gorm:"column:meter_type;type:meters_type[];not null;"`
	ProcessType   string         `gorm:"column:process_type"`
	Name          string         `gorm:"column:name"`
	SchedulerId   string         `gorm:"column:scheduler_id"`
	Scheduler     string         `gorm:"column:scheduler"`
}

func (p BillingMeasuresScheduler) toDomain() billing_measures.Scheduler {
	var distributorID string
	if p.DistributorID != nil {
		distributorID = *p.DistributorID
	}
	return billing_measures.Scheduler{
		ID:            p.ID.String(),
		DistributorId: distributorID,
		Name:          p.Name,
		SchedulerId:   p.SchedulerId,
		ServiceType:   p.ServiceType,
		PointType:     p.PointType,
		MeterType:     p.MeterType,
		ProcessType:   p.ProcessType,
		Format:        p.Scheduler,
	}
}

func NewProcessBillingSchedulerPostgres(db *gorm.DB) *BillingSchedulerRepository {
	return &BillingSchedulerRepository{db: db, dbTimeout: time.Second * 5}
}

type BillingSchedulerRepository struct {
	db        *gorm.DB
	dbTimeout time.Duration
}

func (r BillingSchedulerRepository) DeleteScheduler(ctx context.Context, id string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	var entity BillingMeasuresScheduler
	tx := r.db.WithContext(ctxTimeout).Delete(&entity, "id = ? ", id)
	return tx.Error
}

func (r BillingSchedulerRepository) validOrders() map[string]struct{} {
	return map[string]struct{}{
		"created_at": struct{}{},
	}
}

func (r BillingSchedulerRepository) ListScheduler(ctx context.Context, search db.Pagination) ([]billing_measures.Scheduler, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	items := make([]BillingMeasuresScheduler, 0)

	paginate := NewPaginate(search.Limit, search.Offset)
	query := r.db.WithContext(ctxTimeout).
		Scopes(
			WithPaginate(paginate),
			WithOrder(r.validOrders(), map[string]string{"created_at": "desc"}),
		)
	tx, count := FetchAndCount(query.Find(&items))
	if tx.Error != nil {
		return []billing_measures.Scheduler{}, 0, tx.Error
	}

	return utils.MapSlice(items, func(item BillingMeasuresScheduler) billing_measures.Scheduler {
		return item.toDomain()
	}), int(count), tx.Error

}
func (r BillingSchedulerRepository) GetScheduler(ctx context.Context, id string) (billing_measures.Scheduler, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity BillingMeasuresScheduler
	tx := r.db.WithContext(ctxTimeout).First(&entity, "id = ?", id)
	if tx.Error != nil {
		return billing_measures.Scheduler{}, tx.Error
	}
	return entity.toDomain(), nil
}

func (r BillingSchedulerRepository) SearchScheduler(ctx context.Context, search billing_measures.SearchScheduler) ([]billing_measures.Scheduler, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var response []BillingMeasuresScheduler

	tx := r.db.WithContext(ctxTimeout).Debug().Where("distributor_id = ?", search.DistributorId)
	if search.PointType != "" {
		tx.Where("service_type = ?", search.ServiceType)
	}
	if search.ServiceType != "" {
		tx.Where("point_type = ?", search.PointType)
	}
	if search.ProcessType != "" {
		tx.Where("process_type = ?", search.ProcessType)
	}
	if len(search.MeterType) != 0 {
		tx.Where("meter_type && ?", pq.StringArray(search.MeterType))
	}
	tx.Find(&response)

	return utils.MapSlice(response, func(item BillingMeasuresScheduler) billing_measures.Scheduler {
		return item.toDomain()
	}), tx.Error
}

func (r BillingSchedulerRepository) toDb(s billing_measures.Scheduler) BillingMeasuresScheduler {
	var distributorID *string
	if s.DistributorId != "" {
		distributorID = &s.DistributorId
	}
	return BillingMeasuresScheduler{
		ModelEntity: ModelEntity{
			ID: uuid.FromStringOrNil(s.ID),
		},
		DistributorID: distributorID,
		Distributor:   Distributor{},
		ServiceType:   s.ServiceType,
		PointType:     s.PointType,
		MeterType:     s.MeterType,
		ProcessType:   s.ProcessType,
		Name:          s.Name,
		SchedulerId:   s.SchedulerId,
		Scheduler:     s.Format,
	}

}

func (r BillingSchedulerRepository) SaveScheduler(ctx context.Context, s billing_measures.Scheduler) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	entity := r.toDb(s)

	tx := r.db.WithContext(ctxTimeout).Save(&entity)

	return tx.Error
}
