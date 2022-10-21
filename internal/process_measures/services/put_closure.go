package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"time"
)

type UpdateClosure struct {
	Location        *time.Location
	mongoRepository process_measures.ProcessMeasureClosureRepository
}

func NewUpdateClosureService(mongoRepository process_measures.ProcessMeasureClosureRepository) *UpdateClosure {
	return &UpdateClosure{mongoRepository: mongoRepository}
}

type UpdateClosureDto struct {
	Monthly process_measures.ProcessedMonthlyClosure
}

func (c UpdateClosure) Handler(ctx context.Context, dto UpdateClosureDto) error {

	monthly := dto.Monthly

	result, err := c.mongoRepository.GetClosureOne(ctx, monthly.Id)

	if err != nil {
		return err
	}

	result.CalendarPeriods = monthly.CalendarPeriods
	result.CUPS = monthly.CUPS

	result.EndDate = monthly.EndDate
	result.StartDate = monthly.StartDate
	result.MeterType = monthly.MeterType
	result.MeterSerialNumber = monthly.MeterSerialNumber
	result.Origin = monthly.Origin
	result.ValidationStatus = monthly.ValidationStatus

	err = c.mongoRepository.UpdateClosure(ctx, monthly.Id, result)

	return err

}
