package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"time"
)

type SelfConsumptionServices struct {
	GetSelfConsumptionByCup               *GetSelfConsumptionByCup
	GetSelfConsumptionActiveByDistributor *GetSelfConsumptionActiveByDistributor
	GetBillingSelfConsumptionByCauService *GetBillingSelfConsumptionByCauService
	Location                              *time.Location
}

func NewSelfConsumptionServices(repository billing_measures.SelfConsumptionRepository, billingSelfConsumptionRepository billing_measures.BillingSelfConsumptionRepository, loc *time.Location) *SelfConsumptionServices {
	return &SelfConsumptionServices{
		Location:                              loc,
		GetSelfConsumptionByCup:               NewGetSelfConsumptionByCupService(repository),
		GetSelfConsumptionActiveByDistributor: NewGetSelfConsumptionActiveByDistributorService(repository),
		GetBillingSelfConsumptionByCauService: NewGetBillingSelfConsumptionByCauService(billingSelfConsumptionRepository, loc),
	}
}
