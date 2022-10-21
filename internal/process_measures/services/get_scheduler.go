package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
)

type GetSchedulerById struct {
	repository process_measures.SchedulerRepository
}

func NewGetSchedulerByIdService(repository process_measures.SchedulerRepository) *GetSchedulerById {
	return &GetSchedulerById{repository: repository}
}

type GetSchedulerDTO struct {
	ID string
}

func (s GetSchedulerById) Handler(ctx context.Context, dto GetSchedulerDTO) (process_measures.Scheduler, error) {
	if dto.ID == "" {
		return process_measures.Scheduler{}, process_measures.ErrSchedulerIdFormat
	}

	sc, err := s.repository.GetScheduler(ctx, dto.ID)

	if err != nil {
		return process_measures.Scheduler{}, err
	}
	return sc, nil

}
