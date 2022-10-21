package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"context"
	"time"
)

type SearchConsumProfileDTO struct {
	StartDate time.Time
	EndDate   time.Time
}

type SearchConsumProfile struct {
	repository billing_measures.ConsumProfileRepository
}

func NewSearchConsumProfile(repository billing_measures.ConsumProfileRepository) *SearchConsumProfile {
	return &SearchConsumProfile{repository: repository}
}

func (svc SearchConsumProfile) Handler(ctx context.Context, dto SearchConsumProfileDTO) ([]billing_measures.ConsumProfile, error) {

	cps, err := svc.repository.Search(ctx, billing_measures.QueryConsumProfile{
		EndDate:   dto.EndDate,
		StartDate: dto.StartDate,
	})
	return cps, err
}
