package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetServicePointService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewServicePointService(inventoryRepository inventory.RepositoryInventory) *GetServicePointService {
	return &GetServicePointService{inventoryRepository: inventoryRepository}
}

func (s GetServicePointService) Handler(ctx context.Context, dto inventory.ServicePointSchedulerDto) ([]inventory.ServicePointScheduler, error) {
	sp, err := s.inventoryRepository.GetServicePoints(ctx, dto)
	return sp, err
}
