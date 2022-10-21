package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type ListValidationMeasureDto struct {
	Limit  int
	Offset *int
}

type ListValidationMeasureService struct {
	repository validations.ValidationsRepository
}

func NewListValidationMeasureService(repository validations.ValidationsRepository) *ListValidationMeasureService {
	return &ListValidationMeasureService{repository: repository}
}

func (s ListValidationMeasureService) Handle(ctx context.Context, dto ListValidationMeasureDto) ([]validations.ValidationMeasure, int, error) {

	result, count, err := s.repository.ListValidationMeasure(ctx, validations.SearchValidationMeasure{
		Limit:  dto.Limit,
		Offset: dto.Offset,
	})

	if err != nil {
		return []validations.ValidationMeasure{}, 0, err
	}

	return result, count, err
}
