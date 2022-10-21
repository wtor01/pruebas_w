package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type GetValidationMeasureConfigDto struct {
	ID            string
	DistributorID string
}

type GetValidationMeasureConfigService struct {
	repository validations.ValidationsRepository
}

func NewGetValidationMeasureConfigService(repository validations.ValidationsRepository) *GetValidationMeasureConfigService {
	return &GetValidationMeasureConfigService{repository: repository}
}

func (s GetValidationMeasureConfigService) Handle(ctx context.Context, dto GetValidationMeasureConfigDto) (validations.ValidationMeasureConfig, error) {

	config, err := s.repository.GetDistributorValidationMeasureConfig(ctx, dto.DistributorID, dto.ID)
	if err != nil {
		return validations.ValidationMeasureConfig{}, err
	}

	return config, err
}
