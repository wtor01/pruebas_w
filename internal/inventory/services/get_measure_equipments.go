package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetMeasureEquipmentDto struct {
	Q             string
	Limit         int
	Offset        *int
	Sort          map[string]string
	DistributorID string
}

type GetMeasureEquipmentService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetMeasureEquipmentService(inventoryRepository inventory.RepositoryInventory) *GetMeasureEquipmentService {
	return &GetMeasureEquipmentService{inventoryRepository: inventoryRepository}
}

func (s GetMeasureEquipmentService) Handler(ctx context.Context, dto GetMeasureEquipmentDto) ([]inventory.MeasureEquipment, int, error) {
	ds, count, err := s.inventoryRepository.ListMeters(ctx, inventory.Search{
		Q:             dto.Q,
		Limit:         dto.Limit,
		Offset:        dto.Offset,
		Sort:          dto.Sort,
		DistributorID: dto.DistributorID,
	})

	return ds, count, err
}
