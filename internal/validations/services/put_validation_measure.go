package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
)

type PutValidationMeasure struct {
	repository process_measures.ProcessedMeasureRepository
}

func NewPutValidationMeasure(repository process_measures.ProcessedMeasureRepository) *PutValidationMeasure {
	return &PutValidationMeasure{repository: repository}
}

type PutValidationMeasureDTO struct {
	MeasureType      string
	Status           string
	InvalidationCode string
	ID               string
}

func (svc PutValidationMeasure) Handle(ctx context.Context, dto PutValidationMeasureDTO) error {
	switch dto.MeasureType {
	case "curve":
		{
			curve, err := svc.repository.GetLoadCurveByID(ctx, dto.ID)
			if err != nil {
				return err
			}
			if dto.InvalidationCode != "" {
				curve.InvalidationCodes = append(curve.InvalidationCodes, dto.InvalidationCode)
			}
			if dto.Status != "" {
				curve.ValidationStatus = measures.Status(dto.Status)
			}
			curves := make([]process_measures.ProcessedLoadCurve, 0)
			curves = append(curves, curve)
			err = svc.repository.SaveAllProcessedLoadCurve(ctx, curves)
			if err != nil {
				return err
			}

		}
	case "daily-closure":
		{
			dailyClosure, err := svc.repository.GetDailyClosureByID(ctx, dto.ID)
			if err != nil {
				return err
			}
			if dto.InvalidationCode != "" {
				dailyClosure.InvalidationCodes = append(dailyClosure.InvalidationCodes, dto.InvalidationCode)
			}
			if dto.Status != "" {
				dailyClosure.ValidationStatus = measures.Status(dto.Status)
			}

			err = svc.repository.SaveDailyClosure(ctx, dailyClosure)
			if err != nil {
				return err
			}

		}
	case "monthly-closure":
		{
			monthlyClosure, err := svc.repository.GetMonthlyClosureByID(ctx, dto.ID)
			if err != nil {
				return err
			}
			if dto.InvalidationCode != "" {
				monthlyClosure.InvalidationCodes = append(monthlyClosure.InvalidationCodes, dto.InvalidationCode)
			}
			if dto.Status != "" {
				monthlyClosure.ValidationStatus = measures.Status(dto.Status)
			}
			err = svc.repository.SaveMonthlyClosure(ctx, monthlyClosure)
			if err != nil {
				return err
			}

		}

	}
	return nil
}
