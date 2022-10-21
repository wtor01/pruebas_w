package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"context"
)

type GetSchedulerById struct {
	repository billing_measures.BillingSchedulerRepository
}

func NewGetSchedulerByIdService(repository billing_measures.BillingSchedulerRepository) *GetSchedulerById {
	return &GetSchedulerById{repository: repository}
}

type GetSchedulerDTO struct {
	ID string
}

func (s GetSchedulerById) Handler(ctx context.Context, dto GetSchedulerDTO) (billing_measures.Scheduler, error) {
	if dto.ID == "" {
		return billing_measures.Scheduler{}, billing_measures.ErrSchedulerIdFormat
	}
	sc, err := s.repository.GetScheduler(ctx, dto.ID)
	if err != nil {
		return billing_measures.Scheduler{}, err
	}
	return sc, nil
}
