package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ValidationsPostgresRepository struct {
	BaseRepository
	redisCache *redis.DataCacheRedis
}

func NewValidationMeasurePostgresRepository(db *gorm.DB, redisCache *redis.DataCacheRedis, dbTimeout time.Duration) *ValidationsPostgresRepository {
	return &ValidationsPostgresRepository{BaseRepository: BaseRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}, redisCache: redisCache}
}

func (r ValidationsPostgresRepository) SaveValidationMeasureConfig(ctx context.Context, v validations.ValidationMeasureConfig) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	k := validationRuleConfigToDB(v)

	tx := r.db.WithContext(ctxTimeout).Save(&k)
	keySearch := "validations*"

	if r.redisCache != nil {
		err := r.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}

	return tx.Error
}

func (r ValidationsPostgresRepository) GetDistributorValidationMeasureConfigForThisValidation(ctx context.Context, distributorId string, validationId string) (validations.ValidationMeasureConfig, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity ValidationRuleConfig

	tx := r.db.WithContext(ctxTimeout).Preload("ValidationRule").First(&entity, "distributor_id = ? AND validation_rule_id = ? ", distributorId, validationId)

	if tx.Error != nil {
		return validations.ValidationMeasureConfig{}, tx.Error
	}

	return entity.toDomain(), nil
}

func (r ValidationsPostgresRepository) GetDistributorValidationMeasureConfig(ctx context.Context, distributorId string, configId string) (validations.ValidationMeasureConfig, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity ValidationRuleConfig

	tx := r.db.WithContext(ctxTimeout).Preload("ValidationRule").First(&entity, "distributor_id = ? AND id = ? ", distributorId, configId)

	if tx.Error != nil {
		return validations.ValidationMeasureConfig{}, tx.Error
	}

	return entity.toDomain(), nil
}

func (r ValidationsPostgresRepository) ListDistributorValidationMeasureConfig(ctx context.Context, search validations.SearchValidationMeasureConfig) ([]validations.ValidationMeasureConfig, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var results []ValidationRule

	keyValidationRedis := fmt.Sprintf("validations:%v:%s", search.DistributorID, search.Type)
	if r.redisCache != nil {
		var resultsRedis []validations.ValidationMeasureConfig
		err := r.redisCache.Get(ctx, &resultsRedis, keyValidationRedis)

		if err == nil {
			return resultsRedis, nil
		}
	}

	query := r.db.WithContext(ctxTimeout).
		Preload("ValidationRuleConfig", func(db *gorm.DB) *gorm.DB {
			return db.Where("validation_rule_configs.distributor_id = ?", search.DistributorID)
		}).
		Joins("LEFT JOIN validation_rule_configs ON validation_rule_configs.validation_rule_id = validation_rules.id AND validation_rule_configs.distributor_id = ?", search.DistributorID).
		Where("validation_rules.enabled = true OR validation_rule_configs.id IS NOT NULL").
		Where("validation_rule_configs.deleted_at IS NULL").
		Order("validation_rule_configs.created_at DESC")

	if search.Type != "" {
		query.Where("validation_rules.type = ?", search.Type)
	}

	tx := query.Find(&results)

	if tx.Error != nil {
		return []validations.ValidationMeasureConfig{}, tx.Error
	}

	resultValidations := utils.MapSlice(results, func(item ValidationRule) validations.ValidationMeasureConfig {
		validation := item.toDomain()

		config := validations.ValidationMeasureConfig{
			Id:                "",
			DistributorID:     search.DistributorID,
			ValidationMeasure: validation,
			Action:            validation.Action,
			Enabled:           validation.Enabled,
			Params:            validation.Params,
		}

		if len(item.ValidationRuleConfig) == 1 {
			configDB := item.ValidationRuleConfig[0]
			configDB.ValidationRule = item
			config = configDB.toDomain()
		}

		return config
	})
	if r.redisCache != nil {
		r.redisCache.Set(ctx, keyValidationRedis, resultValidations)
	}

	return resultValidations, nil
}

func (r ValidationsPostgresRepository) DeleteDistributorValidationMeasureConfig(ctx context.Context, distributorId string, configId string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity ValidationRuleConfig

	tx := r.db.WithContext(ctxTimeout).Delete(&entity, "distributor_id = ? AND id = ? ", distributorId, configId)
	keySearch := fmt.Sprintf("validations:*")
	if r.redisCache != nil {
		err := r.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return tx.Error
}

func (r ValidationsPostgresRepository) SaveValidationMeasure(ctx context.Context, v validations.ValidationMeasure) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	k := validationRuleToDB(v)

	tx := r.db.WithContext(ctxTimeout).Save(&k)

	keySearch := fmt.Sprintf("validations:*")
	if r.redisCache != nil {
		err := r.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}
	return tx.Error
}

func (r ValidationsPostgresRepository) ListValidationMeasure(ctx context.Context, search validations.SearchValidationMeasure) ([]validations.ValidationMeasure, int, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	domainResult := make([]validations.ValidationMeasure, 0)

	equipments := make([]ValidationRule, 0)

	paginate := NewPaginate(search.Limit, search.Offset)

	query := r.db.WithContext(ctxTimeout).Scopes(
		WithPaginate(paginate),
		WithOrder(map[string]struct{}{
			"created_at": {},
		}, map[string]string{
			"created_at": "DESC",
		}),
	)

	if search.Type != "" {
		query.Where("validation_rules.type = ?", search.Type)
	}

	result, count := FetchAndCount(query.Find(&equipments))

	if result.Error != nil {
		return domainResult, 0, result.Error
	}

	for _, ds := range equipments {
		domainResult = append(domainResult, ds.toDomain())
	}
	return domainResult, int(count), nil

}

func (r ValidationsPostgresRepository) GetValidationMeasureByID(ctx context.Context, id string) (validations.ValidationMeasure, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity ValidationRule

	tx := r.db.WithContext(ctxTimeout).First(&entity, "id = ?", id)

	if tx.Error != nil {
		return validations.ValidationMeasure{}, tx.Error
	}

	return entity.toDomain(), nil

}

func (r ValidationsPostgresRepository) DeleteValidationMeasureByID(ctx context.Context, id string) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity ValidationRule

	tx := r.db.WithContext(ctxTimeout).Delete(&entity, "id = ?", id)

	keySearch := fmt.Sprintf("validations:*")
	if r.redisCache != nil {
		err := r.redisCache.Clean(ctx, keySearch)
		if err != nil {
			fmt.Println(err)
		}
	}

	return tx.Error
}
