package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type UpdateValidationMeasureDto struct {
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

type UpdateValidationMeasureService struct {
	repository validations.ValidationsRepository
}

func NewUpdateValidationMeasureService(repository validations.ValidationsRepository) *UpdateValidationMeasureService {
	return &UpdateValidationMeasureService{repository: repository}
}

func (s UpdateValidationMeasureService) Handle(ctx context.Context, dto UpdateValidationMeasureDto) (validations.ValidationMeasure, error) {
	validation, err := s.repository.GetValidationMeasureByID(ctx, dto.Id)

	if err != nil {
		return validations.ValidationMeasure{}, err
	}

	v, err := validation.UpdateValidationMeasure(
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
