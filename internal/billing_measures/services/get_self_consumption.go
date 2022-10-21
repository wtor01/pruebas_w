package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"context"
)

type GetSelfConsumptionByCup struct {
	repository billing_measures.SelfConsumptionRepository
}

func NewGetSelfConsumptionByCupService(repository billing_measures.SelfConsumptionRepository) *GetSelfConsumptionByCup {
	return &GetSelfConsumptionByCup{repository: repository}
}

func (sc GetSelfConsumptionByCup) Handler(ctx context.Context, query billing_measures.GetSelfConsumptionByCUP) (billing_measures.SelfConsumption, error) {
	selfCons, err := sc.repository.GetSelfConsumptionByCUP(ctx, query)

	return selfCons, err

}
