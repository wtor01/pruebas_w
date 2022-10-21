package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/apperrors"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type ProcessSelfConsumption struct {
	selfConsumptionRepository        billing_measures.SelfConsumptionRepository
	billingMeasureRepository         billing_measures.BillingMeasureRepository
	billingSelfConsumptionRepository billing_measures.BillingSelfConsumptionRepository
	tracer                           trace.Tracer
	Location                         *time.Location
	consumptionCoefficientRepository billing_measures.ConsumCoefficientRepository
}

func NewProcessSelfConsumption(
	billingMeasureRepository billing_measures.BillingMeasureRepository,
	selfConsumptionRepository billing_measures.SelfConsumptionRepository,
	billingSelfConsumptionRepository billing_measures.BillingSelfConsumptionRepository,
	Location *time.Location,
	consumptionCoefficientRepository billing_measures.ConsumCoefficientRepository,
) *ProcessSelfConsumption {
	return &ProcessSelfConsumption{
		billingMeasureRepository:         billingMeasureRepository,
		selfConsumptionRepository:        selfConsumptionRepository,
		billingSelfConsumptionRepository: billingSelfConsumptionRepository,
		Location:                         Location,
		tracer:                           telemetry.GetTracer(),
		consumptionCoefficientRepository: consumptionCoefficientRepository,
	}
}

func (svc ProcessSelfConsumption) Handler(ctx context.Context, dto billing_measures.OnSaveBillingMeasurePayload) error {
	billingMeasure, err := svc.billingMeasureRepository.Find(ctx, dto.BillingMeasureId)

	if err != nil {
		return apperrors.NewAppError(err.Error(), false)
	}

	selfConsumption, err := svc.selfConsumptionRepository.GetActiveSelfConsumptionByCUP(
		ctx,
		billing_measures.QueryGetActiveSelfConsumptionByCUP{
			CUP:       billingMeasure.CUPS,
			StartDate: billingMeasure.InitDate,
			EndDate:   billingMeasure.EndDate,
		})

	if err != nil {
		return apperrors.NewAppError(err.Error(), false)
	}

	if selfConsumption.ID == "" {
		return nil
	}

	billingMeasure.Status = billing_measures.PendingCau

	err = svc.billingMeasureRepository.Save(ctx, billingMeasure)

	if err != nil {
		return apperrors.NewAppError(err.Error(), true)
	}

	billingSelfConsumption, _ := svc.billingSelfConsumptionRepository.GetBySelfConsumptionBetweenDates(
		ctx,
		billing_measures.QueryGetBillingSelfConsumption{
			SelfconsumptionId: selfConsumption.ID,
			StartDate:         billingMeasure.InitDate,
			EndDate:           billingMeasure.EndDate,
		})

	if billingSelfConsumption.Id == "" {
		billingSelfConsumption, err = billing_measures.NewBillingSelfConsumption(selfConsumption, billingMeasure)
	} else {
		billingSelfConsumption.AddBillingMeasure(billingMeasure)
	}

	if err != nil {
		return apperrors.NewAppError(err.Error(), true)
	}

	if !billingSelfConsumption.IsRedyToProcess() {
		billingSelfConsumption.SetStatusPendingPoints()
		err = svc.billingSelfConsumptionRepository.Save(ctx, billingSelfConsumption)

		if err != nil {
			return apperrors.NewAppError(err.Error(), true)
		}

		return nil
	}

	billingSelfConsumption.StartProcess()
	err = svc.billingSelfConsumptionRepository.Save(ctx, billingSelfConsumption)

	if err != nil {
		return apperrors.NewAppError(err.Error(), true)
	}

	graph := billing_measures.GenerateSelfConsumptionTree(&billingSelfConsumption, svc.consumptionCoefficientRepository)

	err = graph.Execute(ctx)

	billingSelfConsumption.GraphHistory = graph
	billingSelfConsumption.BeforeSave()

	if err != nil {
		billingSelfConsumption.SetStatusSupervision()

		if errSave := svc.billingSelfConsumptionRepository.Save(ctx, billingSelfConsumption); errSave != nil {
			return apperrors.NewAppError(err.Error(), true)
		}

		return nil
	}

	billingMeasuresToUpdate, err := billingSelfConsumption.UpdateBillingMeasureCurve(ctx, svc.billingMeasureRepository)

	if err != nil {
		return err
	}

	err = svc.billingMeasureRepository.SaveAll(ctx, billingMeasuresToUpdate)

	if err != nil {
		return err
	}

	billingSelfConsumption.FinishProcess()

	err = svc.billingSelfConsumptionRepository.Save(ctx, billingSelfConsumption)

	return err
}
