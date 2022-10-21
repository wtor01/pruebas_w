package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type CreateValidationMeasureDto struct {
	Id          string
	UserID      string
	Name        string
	Action      string
	Enabled     bool
	MeasureType string
	Type        string
	Code        string
	Message     string
	Description string
	Params      validations.Params
}

type CreateValidationMeasureService struct {
	repository validations.ValidationsRepository
}

func NewCreateValidationMeasureService(repository validations.ValidationsRepository) *CreateValidationMeasureService {
	return &CreateValidationMeasureService{repository: repository}
}

func (s CreateValidationMeasureService) Handle(ctx context.Context, dto CreateValidationMeasureDto) (validations.ValidationMeasure, error) {
	v, err := validations.NewValidationMeasure(
		dto.Id,
		dto.UserID,
		dto.Name,
		dto.Action,
		dto.Enabled,
		dto.MeasureType,
		dto.Type,
		dto.Code,
		dto.Message,
		dto.Description,
		dto.Params,
	)

	if err != nil {
		return validations.ValidationMeasure{}, err
	}

	err = s.repository.SaveValidationMeasure(ctx, v)

	if err != nil {
		return validations.ValidationMeasure{}, err
	}

	return v, err
}
