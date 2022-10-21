package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"context"
)

type UpdateScheduler struct {
	repository             billing_measures.BillingSchedulerRepository
	schedulerClientCreator scheduler.ClientCreator
	topic                  string
}

func NewUpdateScheduler(repository billing_measures.BillingSchedulerRepository, schedulerClientCreator scheduler.ClientCreator, topic string) *UpdateScheduler {
	return &UpdateScheduler{repository: repository, schedulerClientCreator: schedulerClientCreator, topic: topic}
}

type UpdateSchedulerDTO struct {
	ID            string
	Name          string
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ProcessType   string
	Scheduler     string
}

func (s UpdateScheduler) Handler(ctx context.Context, dto UpdateSchedulerDTO) (billing_measures.Scheduler, error) {

	sc, err := s.repository.GetScheduler(ctx, dto.ID)

	if err != nil {
		return billing_measures.Scheduler{}, err
	}
	oldScheduler := sc.Clone()

	err = sc.Update(
		dto.DistributorId,
		dto.ServiceType,
		dto.PointType,
		dto.MeterType,
		dto.ProcessType,
		dto.Scheduler,
	)

	if err != nil {
		return billing_measures.Scheduler{}, err
	}

	savedSchedulers, err := s.repository.SearchScheduler(ctx, billing_measures.SearchScheduler{
		DistributorId: sc.DistributorId,
		ServiceType:   sc.ServiceType,
		PointType:     sc.PointType,
		MeterType:     sc.MeterType,
		ProcessType:   sc.ProcessType,
	})

	if err != nil {
		return billing_measures.Scheduler{}, err
	}

	if (len(savedSchedulers) == 1 && savedSchedulers[0].ID != sc.ID) || len(savedSchedulers) > 1 {
		return billing_measures.Scheduler{}, billing_measures.ErrSchedulerExist
	}

	clientScheduler, err := s.schedulerClientCreator(ctx)

	defer clientScheduler.Close()

	err = clientScheduler.UpdateJob(ctx, &sc, s.topic)

	if err != nil {
		return billing_measures.Scheduler{}, err
	}

	err = s.repository.SaveScheduler(ctx, sc)

	if err != nil {
		_ = clientScheduler.UpdateJob(ctx, &oldScheduler, s.topic)
		return billing_measures.Scheduler{}, err
	}

	return sc, err
}
