package scheduler

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func schedulerToResponse(scheduler billing_measures.Scheduler) BillingMeasuresScheduler {
	return BillingMeasuresScheduler{
		Id: scheduler.ID,
		BillingMeasuresSchedulerBase: BillingMeasuresSchedulerBase{
			Name: scheduler.Name,
			BillingMeasuresSchedulerUpdatable: BillingMeasuresSchedulerUpdatable{
				DistributorId: &scheduler.DistributorId,
				MeterType: utils.MapSlice(scheduler.MeterType, func(item string) BillingMeasuresSchedulerUpdatableMeterType {
					return BillingMeasuresSchedulerUpdatableMeterType(item)
				}),
				PointType:   BillingMeasuresSchedulerUpdatablePointType(scheduler.PointType),
				ProcessType: BillingMeasuresSchedulerUpdatableProcessType(scheduler.ProcessType),
				Scheduler:   scheduler.Format,
				ServiceType: BillingMeasuresSchedulerUpdatableServiceType(scheduler.ServiceType),
			},
		},
	}
}
