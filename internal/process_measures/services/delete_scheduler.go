package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	scheduler "bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"context"
)

type DeleteScheduler struct {
	repository             process_measures.SchedulerRepository
	schedulerClientCreator scheduler.ClientCreator
}

func NewDeleteScheduler(repository process_measures.SchedulerRepository, schedulerClientCreator scheduler.ClientCreator) *DeleteScheduler {
	return &DeleteScheduler{repository: repository, schedulerClientCreator: schedulerClientCreator}
}

type DeleteSchedulerDTO struct {
	ID string
}

func (s DeleteScheduler) Handler(ctx context.Context, dto DeleteSchedulerDTO) error {

	if dto.ID == "" {
		return process_measures.ErrSchedulerIdFormat
	}

	sc, err := s.repository.GetScheduler(ctx, dto.ID)

	if err != nil {
		return err
	}

	clientScheduler, err := s.schedulerClientCreator(ctx)

	if err != nil {
		return err
	}

	defer clientScheduler.Close()

	err = clientScheduler.DeleteJob(ctx, sc.SchedulerId)

	if err != nil {
		return err
	}

	err = s.repository.DeleteScheduler(ctx, sc.ID)

	return err
}
