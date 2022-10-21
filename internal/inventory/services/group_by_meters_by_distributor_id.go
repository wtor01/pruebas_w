package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GroupByMetersByDistributorIdDto struct {
	DistributorID string
}

type GroupByMetersByDistributorIdService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGroupByMetersByDistributorIdService(inventoryRepository inventory.RepositoryInventory) *GroupByMetersByDistributorIdService {
	return &GroupByMetersByDistributorIdService{inventoryRepository: inventoryRepository}
}

func (s GroupByMetersByDistributorIdService) Handler(ctx context.Context, dto GroupByMetersByDistributorIdDto) (map[string]int, error) {
	result, err := s.inventoryRepository.GroupMetersByType(ctx, dto.DistributorID)

	return result, err
}
