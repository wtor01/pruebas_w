package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type GetValidationMeasureByIdDto struct {
	ID string
}

type GetValidationMeasureByIdService struct {
	repository validations.ValidationsRepository
}

func NewGetValidationMeasureByIdService(repository validations.ValidationsRepository) *GetValidationMeasureByIdService {
	return &GetValidationMeasureByIdService{repository: repository}
}

func (s GetValidationMeasureByIdService) Handle(ctx context.Context, dto GetValidationMeasureByIdDto) (validations.ValidationMeasure, error) {

	result, err := s.repository.GetValidationMeasureByID(ctx, dto.ID)

	if err != nil {
		return validations.ValidationMeasure{}, err
	}

	return result, err
}
