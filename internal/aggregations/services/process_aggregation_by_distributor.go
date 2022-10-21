package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"context"
)

type ProcessAggregationByDistributorService struct {
	aggregationRepository aggregations.AggregationRepository
}

func NewProcessAggregationByDistributorService(
	aggregationRepository aggregations.AggregationRepository,
) *ProcessAggregationByDistributorService {
	return &ProcessAggregationByDistributorService{
		aggregationRepository: aggregationRepository,
	}
}

func (svc ProcessAggregationByDistributorService) Handler(ctx context.Context, dto aggregations.ConfigScheduler) error {
	agg, err := svc.aggregationRepository.GenerateAggregation(ctx, dto)

	if err != nil {
		return err
	}

	if len(agg) == 0 {
		return nil
	}

	err = svc.aggregationRepository.SaveAllAggregations(ctx, agg)

	return err
}
