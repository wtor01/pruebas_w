package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type ListValidationMeasureConfigDto struct {
	Type          string
	DistributorID string
}

type ListValidationMeasureConfigService struct {
	repository validations.ValidationsRepository
}

func NewListValidationMeasureConfigService(repository validations.ValidationsRepository) *ListValidationMeasureConfigService {
	return &ListValidationMeasureConfigService{repository: repository}
}

func (s ListValidationMeasureConfigService) Handle(ctx context.Context, dto ListValidationMeasureConfigDto) ([]validations.ValidationMeasureConfig, error, []validations.ValidationMeasureConfig) {

	config, err := s.repository.ListDistributorValidationMeasureConfig(ctx, validations.SearchValidationMeasureConfig{
		Type:          dto.Type,
		DistributorID: dto.DistributorID,
	})
	if err != nil {
		return []validations.ValidationMeasureConfig{}, err, nil
	}

	return config, err, nil
}
