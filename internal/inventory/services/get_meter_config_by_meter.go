package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
	"time"
)

type GetMeterConfigByMeterDto struct {
	MeterSerialNumber string
	Date              time.Time
}

type GetMeterConfigByMeterService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetMeterConfigByMeterService(repositoryInventory inventory.RepositoryInventory) *GetMeterConfigByMeterService {
	return &GetMeterConfigByMeterService{
		inventoryRepository: repositoryInventory,
	}
}

func (s GetMeterConfigByMeterService) Handler(ctx context.Context, dto GetMeterConfigByMeterDto) (inventory.MeterConfig, error) {
	d, err := s.inventoryRepository.GetMeterConfig(ctx, inventory.GetMeterConfigByMeterQuery{
		MeterSerialNumber: dto.MeterSerialNumber,
		Date:              dto.Date,
	})

	return d, err
}
