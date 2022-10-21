package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"context"
)

type TaxMeasuresByCups struct {
	repository billing_measures.BillingMeasuresDashboardRepository
}

func NewTaxMeasuresByCups(
	repository billing_measures.BillingMeasuresDashboardRepository,

) *TaxMeasuresByCups {
	return &TaxMeasuresByCups{
		repository: repository,
	}
}

func (svc TaxMeasuresByCups) Handler(ctx context.Context, dto billing_measures.QueryBillingMeasuresTax) (billing_measures.BillingMeasuresTaxResult, error) {

	result, err := svc.repository.GetBillingMeasuresTax(ctx, dto)

	return result, err
}
