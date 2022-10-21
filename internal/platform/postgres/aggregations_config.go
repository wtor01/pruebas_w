package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type AggregationsConfig struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
	Name        string                 `gorm:"column:name;not null;unique;<-:create"`
	Description *string                `gorm:"column:description"`
	Scheduler   string                 `gorm:"column:scheduler;not null"`
	SchedulerId string                 `gorm:"column:scheduler_id"`
	StartDate   time.Time              `gorm:"column:start_date;not null"`
	EndDate     *time.Time             `gorm:"column:end_date"`
	Features    []AggregationsFeatures `gorm:"many2many:aggregation_config_features;constraint:OnUpdate:CASCADE"`
}

func (a AggregationsConfig) toDomain() aggregations.Config {
	aggregationConfig := aggregations.Config{
		Id:          a.ID.String(),
		Name:        a.Name,
		Scheduler:   a.Scheduler,
		StartDate:   a.StartDate,
		SchedulerId: a.SchedulerId,
		Features: utils.MapSlice(a.Features, func(item AggregationsFeatures) aggregations.Features {
			return item.toDomain()
		}),
	}

	if a.EndDate != nil {
		aggregationConfig.EndDate = *a.EndDate
	}

	if a.Description != nil {
		aggregationConfig.Description = *a.Description
	}

	return aggregationConfig
}

type AggregationsConfigRepository struct {
	db        *gorm.DB
	dbTimeout time.Duration
}

func NewAggregationsConfigRepository(db *gorm.DB) *AggregationsConfigRepository {
	return &AggregationsConfigRepository{
		db:        db,
		dbTimeout: time.Second * 5,
	}
}

func (r AggregationsConfigRepository) GetAggregationConfigs(ctx context.Context, q aggregations.GetConfigsQuery) ([]aggregations.Config, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	result := make([]AggregationsConfig, 0)

	query := r.db.
		WithContext(ctxTimeout).
		Scopes(
			WithLike(q.Query),
			WithPaginate(NewPaginate(q.Limit, q.Offset))).
		Preload("Features")

	tx, count := FetchAndCount(query.Find(&result))
	if tx.Error != nil {
		return []aggregations.Config{}, 0, tx.Error
	}

	configs := utils.MapSlice(result, func(item AggregationsConfig) aggregations.Config {
		return item.toDomain()
	})

	return configs, int(count), nil
}

func (r AggregationsConfigRepository) GetAggregationConfigById(ctx context.Context, aggregationConfigId string) (aggregations.Config, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var result AggregationsConfig

	tx := r.db.
		WithContext(ctxTimeout).
		Where("id = ?", aggregationConfigId).
		Preload("Features").
		First(&result)

	if tx.Error != nil {
		return aggregations.Config{}, tx.Error
	}

	return result.toDomain(), nil
}

func (r AggregationsConfigRepository) SaveAggregationConfig(ctx context.Context, config aggregations.Config) (aggregations.Config, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	entity := r.toDb(config)

	tx := r.db.WithContext(ctxTimeout).Begin()

	if err := tx.Omit("Features").Save(&entity).Error; err != nil {
		tx.Rollback()
		return aggregations.Config{}, err
	}

	if err := tx.Model(&entity).Association("Features").Replace(entity.Features); err != nil {
		tx.Rollback()
		return aggregations.Config{}, err
	}

	tx = tx.Commit()

	if tx.Error != nil {
		tx.Rollback()
		return aggregations.Config{}, tx.Error
	}

	return entity.toDomain(), nil

}

func (r AggregationsConfigRepository) DeleteAggregationConfig(ctx context.Context, aggregationConfigId string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	aggregationConfig := r.toDb(aggregations.Config{Id: aggregationConfigId})
	tx := r.db.WithContext(ctxTimeout).Begin()

	if err := tx.Model(&aggregationConfig).Association("Features").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&aggregationConfig).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx = tx.Commit()

	return tx.Error
}

func (r AggregationsConfigRepository) toDb(config aggregations.Config) AggregationsConfig {
	aggregationConfig := AggregationsConfig{
		ID:          uuid.FromStringOrNil(config.Id),
		Name:        config.Name,
		Scheduler:   config.Scheduler,
		SchedulerId: config.SchedulerId,
		StartDate:   config.StartDate,
		Features: utils.MapSlice(config.Features, func(item aggregations.Features) AggregationsFeatures {
			return AggregationsFeatures{
				ID: uuid.FromStringOrNil(item.ID),
			}
		}),
	}

	if !config.EndDate.IsZero() {
		aggregationConfig.EndDate = &config.EndDate
	}

	if config.Description != "" {
		aggregationConfig.Description = &config.Description
	}

	return aggregationConfig
}
