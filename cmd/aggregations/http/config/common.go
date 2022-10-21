package config

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func AggregationConfigToResponse(config aggregations.Config) AggregationConfig {
	aggregationConfig := AggregationConfig{
		Id: config.Id,
		AggregationConfigBase: AggregationConfigBase{
			Features:  utils.MapSlice(config.Features, AggregationConfigFeatureToResponse),
			Name:      config.Name,
			Scheduler: config.Scheduler,
			StartDate: config.StartDate,
		},
	}

	if !config.EndDate.IsZero() {
		aggregationConfig.EndDate = &config.EndDate
	}

	if config.Description != "" {
		aggregationConfig.Description = &config.Description
	}

	return aggregationConfig
}

func AggregationConfigFeatureToResponse(configFeature aggregations.Features) AggregationFeature {
	return AggregationFeature{
		Id: configFeature.ID,
		AggregationFeaturesBase: AggregationFeaturesBase{
			Field: configFeature.Field,
			Name:  configFeature.Name,
		},
	}
}
