package features

import "bitbucket.org/sercide/data-ingestion/internal/aggregations"

func featuresToResponse(features aggregations.Features) AggregationFeature {
	return AggregationFeature{
		Id: features.ID,
		AggregationFeaturesBase: AggregationFeaturesBase{
			Field: features.Field,
			Name:  features.Name,
		},
	}
}
