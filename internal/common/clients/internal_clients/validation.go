package internal_clients

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	validations2 "bitbucket.org/sercide/data-ingestion/internal/validations"
	validations3 "bitbucket.org/sercide/data-ingestion/internal/validations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
)

type Validation struct {
	validationServices *validations3.Services
}

func NewValidation(validationServices *validations3.Services) *Validation {
	return &Validation{validationServices: validationServices}
}

func (client Validation) GetValidationConfigList(ctx context.Context, distributorID string, valType string) ([]clients.ValidationConfig, error) {
	v, err, _ := client.validationServices.ListValidationMeasureConfigService.Handle(ctx, validations3.ListValidationMeasureConfigDto{Type: valType, DistributorID: distributorID})

	if err != nil {
		return []clients.ValidationConfig{}, err
	}

	return utils.MapSlice(v, validationMeasureConfigToResponse), nil
}

func validationMeasureConfigToResponse(v validations2.ValidationMeasureConfig) clients.ValidationConfig {

	return clients.ValidationConfig{
		Id: v.Id,
		ValidationMeasure: struct {
			Id          string                         `json:"id"`
			Name        string                         `json:"name"`
			Action      string                         `json:"action"`
			Enabled     bool                           `json:"enabled"`
			MeasureType string                         `json:"measure_type"`
			Type        string                         `json:"type"`
			Code        string                         `json:"code"`
			Message     string                         `json:"message"`
			Description string                         `json:"description"`
			Params      clients.ValidationConfigParams `json:"validation_measure"`
		}{
			Id:          v.ValidationMeasure.Id,
			Name:        v.ValidationMeasure.Name,
			Action:      v.ValidationMeasure.Action,
			Enabled:     v.ValidationMeasure.Enabled,
			MeasureType: v.ValidationMeasure.MeasureType,
			Type:        v.ValidationMeasure.Type,
			Code:        v.ValidationMeasure.Code,
			Message:     v.ValidationMeasure.Message,
			Description: v.ValidationMeasure.Description,
			Params: clients.ValidationConfigParams{
				Type: v.Params.Type,
				Validations: utils.MapSlice(v.Params.Validations, func(item validations2.ValidationMeasureParamsValidations) clients.ValidationParams {
					return clients.ValidationParams{
						Config:   item.Config,
						Id:       item.Id,
						Name:     item.Name,
						Keys:     item.Keys,
						Required: item.Required,
						Type:     item.Type,
					}
				}),
			},
		},
		Action:  v.Action,
		Enabled: v.Enabled,
		ExtraConfig: utils.MapSlice(v.Params.Validations, func(item validations2.ValidationMeasureParamsValidations) clients.ValidationConfigExtraConfig {
			return clients.ValidationConfigExtraConfig{Id: item.Id, Config: item.Config, Name: item.Name}
		}),
	}
}
