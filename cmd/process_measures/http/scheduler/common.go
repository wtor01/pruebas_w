package scheduler

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func schedulerToResponse(scheduler process_measures.Scheduler) ProcessMeasureScheduler {
	return ProcessMeasureScheduler{
		Id: scheduler.ID,
		ProcessMeasureSchedulerBase: ProcessMeasureSchedulerBase{
			Name: scheduler.Name,
			ProcessMeasureSchedulerUpdatable: ProcessMeasureSchedulerUpdatable{
				Description:   scheduler.Description,
				DistributorId: &scheduler.DistributorId,
				MeterType: utils.MapSlice(scheduler.MeterType, func(item string) ProcessMeasureSchedulerUpdatableMeterType {
					return ProcessMeasureSchedulerUpdatableMeterType(item)
				}),
				PointType:   ProcessMeasureSchedulerUpdatablePointType(scheduler.PointType),
				ReadingType: ProcessMeasureSchedulerUpdatableReadingType(scheduler.ReadingType),
				Scheduler:   scheduler.Format,
				ServiceType: ProcessMeasureSchedulerUpdatableServiceType(scheduler.ServiceType),
			},
		},
	}
}
