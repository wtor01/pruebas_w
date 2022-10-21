package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func validationMeasureToResponse(v validations.ValidationMeasure) ValidationMeasure {

	return ValidationMeasure{
		Id: v.Id,
		ValidationMeasureBase: ValidationMeasureBase{
			Action:      ValidationMeasureBaseAction(v.Action),
			Code:        v.Code,
			Description: &v.Description,
			Enabled:     v.Enabled,
			MeasureType: ValidationMeasureBaseMeasureType(v.MeasureType),
			Message:     v.Message,
			Name:        v.Name,
			Type:        ValidationMeasureBaseType(v.Type),
			Params: Params{
				Type: ParamsType(v.Params.Type),
				Validations: utils.MapSlice[validations.ValidationMeasureParamsValidations, Validation](v.Params.Validations, func(item validations.ValidationMeasureParamsValidations) Validation {
					return Validation{
						Config: &Validation_Config{AdditionalProperties: item.Config},
						Id:     item.Id,
						Name:   item.Name,
						Keys: utils.MapSlice[string, ValidationKeys](item.Keys, func(item string) ValidationKeys {
							return ValidationKeys(item)
						}),
						Required: item.Required,
						Type:     ValidationType(item.Type),
					}
				}),
			},
		},
	}
}

func validationMeasureConfigToResponse(v validations.ValidationMeasureConfig) ValidationMeasureConfig {
	extraConfig := utils.MapSlice(v.Params.Validations, func(item validations.ValidationMeasureParamsValidations) ExtraConfig {
		return ExtraConfig{Id: item.Id, AdditionalProperties: item.Config, Name: item.Name}
	})

	return ValidationMeasureConfig{
		Id: v.Id,
		ValidationMeasureConfigBase: ValidationMeasureConfigBase{
			Action:            ValidationMeasureConfigBaseAction(v.Action),
			Enabled:           v.Enabled,
			ExtraConfig:       &extraConfig,
			ValidationMeasure: validationMeasureToResponse(v.ValidationMeasure),
		},
	}
}
