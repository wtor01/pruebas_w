package validations

import (
	"context"
)

type SearchValidationMeasure struct {
	Limit  int
	Offset *int
	Type   string
}

type SearchValidationMeasureConfig struct {
	Type          string
	DistributorID string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=ValidationsRepository
type ValidationsRepository interface {
	SaveValidationMeasure(ctx context.Context, v ValidationMeasure) error
	ListValidationMeasure(ctx context.Context, search SearchValidationMeasure) ([]ValidationMeasure, int, error)
	GetValidationMeasureByID(ctx context.Context, id string) (ValidationMeasure, error)
	DeleteValidationMeasureByID(ctx context.Context, id string) error

	SaveValidationMeasureConfig(ctx context.Context, v ValidationMeasureConfig) error
	GetDistributorValidationMeasureConfigForThisValidation(ctx context.Context, distributorId string, validationId string) (ValidationMeasureConfig, error)
	GetDistributorValidationMeasureConfig(ctx context.Context, distributorId string, configId string) (ValidationMeasureConfig, error)
	DeleteDistributorValidationMeasureConfig(ctx context.Context, distributorId string, configId string) error
	ListDistributorValidationMeasureConfig(ctx context.Context, search SearchValidationMeasureConfig) ([]ValidationMeasureConfig, error)
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=ValidationMongoRepository
type ValidationMongoRepository interface {
	GetMonthlyClosureByCup(ctx context.Context, q QueryClosedCupsMeasureOnDate) (ProcessedMonthlyClosure, error)
	GetLoadCurveByQuery(ctx context.Context, q QueryCurveCupsMeasureOnDate) ([]ProcessedLoadCurve, error)
}
