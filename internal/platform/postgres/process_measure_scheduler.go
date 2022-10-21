package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type ProcessMeasureScheduler struct {
	ModelEntity
	DistributorID *string        `gorm:"type:uuid;"`
	Distributor   Distributor    `gorm:"distributor" gorm:"foreignKey:DistributorID"`
	ServiceType   string         `gorm:"column:service_type;not null;type:service_point_service_type"`
	PointType     string         `gorm:"column:point_type;not null;type:service_point_point_type;"`
	MeterType     pq.StringArray `gorm:"column:meter_type;type:meters_type[];not null;"`
	ReadingType   string         `gorm:"column:reading_type"`
	Name          string         `gorm:"column:name"`
	Description   string         `gorm:"column:description"`
	SchedulerId   string         `gorm:"column:scheduler_id"`
	Scheduler     string         `gorm:"column:scheduler"`
}

func (p ProcessMeasureScheduler) toDomain() process_measures.Scheduler {
	var distributorID string
	if p.DistributorID != nil {
		distributorID = *p.DistributorID
	}
	return process_measures.Scheduler{
		ID:            p.ID.String(),
		DistributorId: distributorID,
		Name:          p.Name,
		Description:   p.Description,
		SchedulerId:   p.SchedulerId,
		ServiceType:   p.ServiceType,
		PointType:     p.PointType,
		MeterType:     p.MeterType,
		ReadingType:   measures.ReadingType(p.ReadingType),
		Format:        p.Scheduler,
	}
}

func NewProcessMeasureSchedulerPostgres(db *gorm.DB) *SchedulerRepository {
	return &SchedulerRepository{db: db, dbTimeout: time.Second * 5}
}

type SchedulerRepository struct {
	db        *gorm.DB
	dbTimeout time.Duration
}

func (r SchedulerRepository) validOrders() map[string]struct{} {
	return map[string]struct{}{
		"created_at": struct{}{},
	}
}

func (r SchedulerRepository) ListScheduler(ctx context.Context, search db.Pagination) ([]process_measures.Scheduler, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	items := make([]ProcessMeasureScheduler, 0)

	paginate := NewPaginate(search.Limit, search.Offset)
	query := r.db.WithContext(ctxTimeout).
		Scopes(
			WithPaginate(paginate),
			WithOrder(r.validOrders(), map[string]string{"created_at": "desc"}),
		)

	tx, count := FetchAndCount(query.Find(&items))

	if tx.Error != nil {
		return []process_measures.Scheduler{}, 0, tx.Error
	}

	return utils.MapSlice(items, func(item ProcessMeasureScheduler) process_measures.Scheduler {
		return item.toDomain()
	}), int(count), tx.Error
}

func (r SchedulerRepository) DeleteScheduler(ctx context.Context, id string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	var entity ProcessMeasureScheduler

	tx := r.db.WithContext(ctxTimeout).Delete(&entity, "id = ? ", id)

	return tx.Error
}

func (r SchedulerRepository) GetScheduler(ctx context.Context, id string) (process_measures.Scheduler, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity ProcessMeasureScheduler

	tx := r.db.WithContext(ctxTimeout).First(&entity, "id = ?", id)

	if tx.Error != nil {
		return process_measures.Scheduler{}, tx.Error
	}

	return entity.toDomain(), nil
}

func (r SchedulerRepository) toDb(s process_measures.Scheduler) ProcessMeasureScheduler {
	var distributorID *string
	if s.DistributorId != "" {
		distributorID = &s.DistributorId
	}
	return ProcessMeasureScheduler{
		ModelEntity: ModelEntity{
			ID: uuid.FromStringOrNil(s.ID),
		},
		DistributorID: distributorID,
		ServiceType:   s.ServiceType,
		PointType:     s.PointType,
		MeterType:     s.MeterType,
		ReadingType:   string(s.ReadingType),
		Name:          s.Name,
		Description:   s.Description,
		SchedulerId:   s.SchedulerId,
		Scheduler:     s.Format,
	}
}

func (r SchedulerRepository) SaveScheduler(ctx context.Context, s process_measures.Scheduler) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	entity := r.toDb(s)

	tx := r.db.WithContext(ctxTimeout).Save(&entity)

	return tx.Error
}

func (r SchedulerRepository) SearchScheduler(ctx context.Context, search process_measures.SearchScheduler) ([]process_measures.Scheduler, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if search.DistributorId == "" {
		return []process_measures.Scheduler{}, nil
	}

	var response []ProcessMeasureScheduler

	tx := r.db.WithContext(ctxTimeout).Debug().Where("service_type = ?", search.ServiceType).Where("distributor_id = ?", search.DistributorId).Where("point_type = ?", search.PointType).Where("reading_type = ?", search.ReadingType).Where("meter_type && ?", pq.StringArray(search.MeterType)).
		Find(&response)

	return utils.MapSlice(response, func(item ProcessMeasureScheduler) process_measures.Scheduler {
		return item.toDomain()
	}), tx.Error
}

// si curva ayer si cierre hoy
