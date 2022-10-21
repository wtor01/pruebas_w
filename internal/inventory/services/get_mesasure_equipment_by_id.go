package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetMeasureEquipmentByIdDto struct {
	Id            string
	DistributorId string
}

type GetMeasureEquipmentByIdService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetMeasureEquipmentByIdService(inventoryRepository inventory.RepositoryInventory) *GetMeasureEquipmentByIdService {
	return &GetMeasureEquipmentByIdService{inventoryRepository: inventoryRepository}
}

func (s GetMeasureEquipmentByIdService) Handler(ctx context.Context, dto GetMeasureEquipmentByIdDto) (inventory.MeasureEquipment, error) {
	equipment, err := s.inventoryRepository.GetMeasureEquipmentByID(ctx, dto.DistributorId, dto.Id)

	return equipment, err
}
