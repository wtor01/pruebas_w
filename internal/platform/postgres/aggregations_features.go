package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type AggregationsFeatures struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	Name      string `gorm:"unique;column:name;not null"`
	Field     string `gorm:"unique;column:field;not null"`
}

func (f AggregationsFeatures) toDomain() aggregations.Features {
	return aggregations.Features{
		ID:    f.ID.String(),
		Name:  f.Name,
		Field: f.Field,
	}
}

func NewAggregationsFeaturesPostgres(db *gorm.DB) *AggregationsFeaturesRepository {
	return &AggregationsFeaturesRepository{db: db, dbTimeout: time.Second * 5}
}

type AggregationsFeaturesRepository struct {
	db        *gorm.DB
	dbTimeout time.Duration
}

func (r AggregationsFeaturesRepository) GetFeatures(ctx context.Context, id string) (aggregations.Features, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity AggregationsFeatures
	tx := r.db.WithContext(ctxTimeout).First(&entity, "id = ?", id)
	if tx.Error != nil {
		return aggregations.Features{}, tx.Error
	}
	return entity.toDomain(), nil
}

func (r AggregationsFeaturesRepository) validOrders() map[string]struct{} {
	return map[string]struct{}{
		"created_at": struct{}{},
	}
}

func (r AggregationsFeaturesRepository) ListFeatures(ctx context.Context, search db.Pagination) ([]aggregations.Features, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	items := make([]AggregationsFeatures, 0)

	paginate := NewPaginate(search.Limit, search.Offset)
	query := r.db.WithContext(ctxTimeout).Scopes(
		WithPaginate(paginate),
		WithOrder(r.validOrders(), map[string]string{"created_at": "desc"}),
	)
	tx, count := FetchAndCount(query.Find(&items))
	if tx.Error != nil {
		return []aggregations.Features{}, 0, tx.Error
	}

	return utils.MapSlice(items, func(item AggregationsFeatures) aggregations.Features {
		return item.toDomain()
	}), int(count), tx.Error
}

func (r AggregationsFeaturesRepository) DeleteFeatures(ctx context.Context, id string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	var entity AggregationsFeatures
	tx := r.db.WithContext(ctxTimeout).Delete(&entity, "id = ?", id)
	return tx.Error
}

func (r AggregationsFeaturesRepository) SaveFeatures(ctx context.Context, obj aggregations.Features) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	entity := r.toDB(obj)
	tx := r.db.WithContext(ctxTimeout).Save(&entity)
	return tx.Error
}

func (r AggregationsFeaturesRepository) SearchFeatures(ctx context.Context, featuresObj aggregations.SearchFeatures) ([]aggregations.Features, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var response []AggregationsFeatures

	tx := r.db.WithContext(ctxTimeout).Where("name = ?", featuresObj.Name)
	if featuresObj.Field != "" {
		tx.Where("field = ?", featuresObj.Field)
	}
	tx.Find(&response)

	return utils.MapSlice(response, func(item AggregationsFeatures) aggregations.Features {
		return item.toDomain()
	}), tx.Error
}

func (r AggregationsFeaturesRepository) toDB(featuresObj aggregations.Features) AggregationsFeatures {
	return AggregationsFeatures{
		ID:    uuid.FromStringOrNil(featuresObj.ID),
		Name:  featuresObj.Name,
		Field: featuresObj.Field,
	}
}

func (r AggregationsFeaturesRepository) GetFeaturesByIds(ctx context.Context, ids []string) ([]aggregations.Features, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var response []AggregationsFeatures

	tx := r.db.WithContext(ctxTimeout).Where("id IN ?", ids)

	tx.Find(&response)

	return utils.MapSlice(response, func(item AggregationsFeatures) aggregations.Features {
		return item.toDomain()
	}), tx.Error
}
