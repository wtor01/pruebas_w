package aggregations

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func transformToResponseAggregations(agg aggregations.Aggregation) Aggregation {

	features := make([]Features, 0)

	for _, value := range agg.Parameters {
		features = append(features, Features{
			Feature: &AggregationFeature{
				Id: value.ID,
				AggregationFeaturesBase: AggregationFeaturesBase{
					Field: value.Field,
					Name:  value.Name,
				},
			},
			Value: &value.Value,
		})
	}
	response := Aggregation{
		Id: agg.Id,
		AggregationBase: AggregationBase{
			AggregationBefore:   nil,
			AggregationConfigId: agg.Id,
			Date:                agg.GenerationDate,
			Features:            features,
		},
	}

	return response
}

func transformToResponseAggregation(agg aggregations.AggregationPrevious) AggregationWithCUPS {

	cupsCurrent := utils.MapSlice(agg.CurrentAggregationCups, func(item aggregations.CupsStateAggregation) AggregationCUPSCurrent {
		typeCups := AggregationCUPSCurrentTypeNEUTRAL
		if item.Type == "IN" {
			typeCups = AggregationCUPSCurrentTypeIN
		}
		return AggregationCUPSCurrent{
			CUPS: &item.CUPS,
			Type: &typeCups,
		}
	})

	cupsPrevious := utils.MapSlice(agg.PreviousAggregationCups, func(item aggregations.CupsStateAggregation) AggregationCUPSPrevious {
		typeCups := AggregationCUPSPreviousTypeNEUTRAL
		if item.Type == "OUT" {
			typeCups = AggregationCUPSPreviousTypeOUT
		}
		return AggregationCUPSPrevious{
			CUPS: &item.CUPS,
			Type: &typeCups,
		}
	})

	return AggregationWithCUPS{
		ListCUPSCurrent:  &cupsCurrent,
		ListCUPSPrevious: &cupsPrevious,
		Aggregation: Aggregation{
			Id: agg.Id,
			AggregationBase: AggregationBase{
				AggregationBefore:   &agg.PreviousAggregation.Id,
				AggregationConfigId: agg.TypeId,
				Date:                agg.GenerationDate,
				Features: utils.MapSlice(agg.Parameters, func(item aggregations.FeatureValue) Features {
					return Features{
						Feature: &AggregationFeature{
							Id: item.ID,
							AggregationFeaturesBase: AggregationFeaturesBase{
								Field: item.Field,
								Name:  item.Name,
							},
						},
						Value: &item.Value,
					}
				}),
			},
		},
	}
}
