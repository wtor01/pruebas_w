package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"context"
)

type GetSelfConsumptionActiveByDistributor struct {
	repository billing_measures.SelfConsumptionRepository
}

func NewGetSelfConsumptionActiveByDistributorService(repository billing_measures.SelfConsumptionRepository) *GetSelfConsumptionActiveByDistributor {
	return &GetSelfConsumptionActiveByDistributor{repository: repository}
}

func (sc GetSelfConsumptionActiveByDistributor) Handler(ctx context.Context, dto billing_measures.GetSelfConsumptionByDistributortDto) ([]billing_measures.SelfConsumption, int, error) {
	selfCons, count, err := sc.repository.GetSelfConsumptionActiveByDistributor(ctx, billing_measures.GetSelfConsumptionByDistributortDto{
		Limit:         dto.Limit,
		Date:          dto.Date,
		DistributorId: dto.DistributorId,
		Offset:        dto.Offset,
	})

	return selfCons, count, err

}
