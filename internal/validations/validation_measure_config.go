package validations

import "bitbucket.org/sercide/data-ingestion/pkg/utils"

type Config struct {
	ID    string
	Extra map[string]string
}

type ValidationMeasureConfig struct {
	Id                string
	Name              string
	DistributorID     string
	ValidationMeasure ValidationMeasure
	Action            string
	Enabled           bool
	Params            ValidationMeasureParams
	CreatedByID       string
	UpdatedByID       string
}

func NewValidationMeasureConfig(id string, userId string, validation ValidationMeasure, distributorID, action string, enabled bool, config []Config) (ValidationMeasureConfig, error) {

	configMap := make(map[string]map[string]string)

	for _, c := range config {
		configMap[c.ID] = c.Extra
	}

	params, err := NewValidationMeasureParams(validation.Params.Type, utils.MapSlice[ValidationMeasureParamsValidations, Validation](validation.Params.Validations, func(item ValidationMeasureParamsValidations) Validation {

		newConfig, ok := configMap[item.Id]
		if !ok || newConfig == nil {
			newConfig = item.Config
		}

		return Validation{
			Id:       item.Id,
			Type:     item.Type,
			Keys:     item.Keys,
			Required: item.Required,
			Config:   newConfig,
		}
	}))

	if err != nil {
		return ValidationMeasureConfig{}, err
	}

	return ValidationMeasureConfig{
		Id:                id,
		DistributorID:     distributorID,
		ValidationMeasure: validation,
		Action:            action,
		Enabled:           enabled,
		Params:            params,
		CreatedByID:       userId,
	}, nil
}

func (v *ValidationMeasureConfig) UpdateValidationMeasureConfig(action string, userId string, enabled bool, config []Config) error {
	configMap := make(map[string]map[string]string)

	for _, c := range config {
		configMap[c.ID] = c.Extra
	}

	params, err := NewValidationMeasureParams(v.ValidationMeasure.Params.Type, utils.MapSlice[ValidationMeasureParamsValidations, Validation](v.ValidationMeasure.Params.Validations, func(item ValidationMeasureParamsValidations) Validation {

		newConfig, ok := configMap[item.Id]
		if !ok {
			newConfig = item.Config
		}

		return Validation{
			Id:       item.Id,
			Type:     item.Type,
			Keys:     item.Keys,
			Required: item.Required,
			Config:   newConfig,
		}
	}))

	if err != nil {
		return err
	}

	v.Action = action
	v.UpdatedByID = userId
	v.Enabled = enabled
	v.Params = params

	return nil
}
