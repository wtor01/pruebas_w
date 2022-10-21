package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type DeleteValidationMeasureByIdDto struct {
	ID string
}

type DeleteValidationMeasureByIdService struct {
	repository validations.ValidationsRepository
}

func NewDeleteValidationMeasureByIdService(repository validations.ValidationsRepository) *DeleteValidationMeasureByIdService {
	return &DeleteValidationMeasureByIdService{repository: repository}
}

func (s DeleteValidationMeasureByIdService) Handle(ctx context.Context, dto DeleteValidationMeasureByIdDto) error {

	err := s.repository.DeleteValidationMeasureByID(ctx, dto.ID)

	return err
}
