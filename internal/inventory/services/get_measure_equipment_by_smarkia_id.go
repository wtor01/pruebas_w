package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetMeasureEquipmentBySmarkiaDto struct {
	Id string
}

type GetMeasureEquipmentBySmarkiaIdService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetMeasureEquipmentBySmarkiaIdService(inventoryRepository inventory.RepositoryInventory) *GetMeasureEquipmentBySmarkiaIdService {
	return &GetMeasureEquipmentBySmarkiaIdService{inventoryRepository: inventoryRepository}
}

func (s GetMeasureEquipmentBySmarkiaIdService) Handler(ctx context.Context, dto GetMeasureEquipmentBySmarkiaDto) (inventory.MeasureEquipment, error) {
	ds, err := s.inventoryRepository.GetMeasureEquipmentBySmarkiaID(ctx, dto.Id)

	return ds, err
}
