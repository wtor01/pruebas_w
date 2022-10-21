package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type DeleteValidationMeasureConfigDto struct {
	ID            string
	DistributorID string
}

type DeleteValidationMeasureConfigService struct {
	repository validations.ValidationsRepository
}

func NewDeleteValidationMeasureConfigService(repository validations.ValidationsRepository) *DeleteValidationMeasureConfigService {
	return &DeleteValidationMeasureConfigService{repository: repository}
}

func (s DeleteValidationMeasureConfigService) Handle(ctx context.Context, dto DeleteValidationMeasureConfigDto) error {

	err := s.repository.DeleteDistributorValidationMeasureConfig(ctx, dto.DistributorID, dto.ID)

	return err
}
