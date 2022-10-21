package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
)

type CreateValidationMeasureConfigDto struct {
	ID                  string
	ValidationMeasureID string
	DistributorID       string
	UserID              string
	Action              string
	Enabled             bool
	Params              []validations.Config
}

type CreateValidationMeasureConfigService struct {
	repository validations.ValidationsRepository
}

func NewCreateValidationMeasureConfigService(repository validations.ValidationsRepository) *CreateValidationMeasureConfigService {
	return &CreateValidationMeasureConfigService{repository: repository}
}

func (s CreateValidationMeasureConfigService) Handle(ctx context.Context, dto CreateValidationMeasureConfigDto) (validations.ValidationMeasureConfig, error) {

	validation, err := s.repository.GetValidationMeasureByID(ctx, dto.ValidationMeasureID)
	if err != nil {
		return validations.ValidationMeasureConfig{}, err
	}

	config, _ := s.repository.GetDistributorValidationMeasureConfigForThisValidation(ctx, dto.DistributorID, dto.ValidationMeasureID)

	if config.Id != "" {
		err = config.UpdateValidationMeasureConfig(dto.Action, dto.UserID, dto.Enabled, dto.Params)
	} else {
		config, err = validations.NewValidationMeasureConfig(dto.ID, dto.UserID, validation, dto.DistributorID, dto.Action, dto.Enabled, dto.Params)
	}

	if err != nil {
		return validations.ValidationMeasureConfig{}, err
	}

	err = s.repository.SaveValidationMeasureConfig(ctx, config)

	if err != nil {
		return validations.ValidationMeasureConfig{}, err
	}

	return config, err
}
