package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetServicePointByMeterService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetServicePointByMeterService(inventoryRepository inventory.RepositoryInventory) *GetServicePointByMeterService {
	return &GetServicePointByMeterService{inventoryRepository: inventoryRepository}
}

func (s GetServicePointByMeterService) Handler(ctx context.Context, dto inventory.GetMeasureEquipmentByMeterQuery) (inventory.ServicePointScheduler, error) {
	sp, err := s.inventoryRepository.GetServicePointByMeter(ctx, dto)
	return sp, err
}
