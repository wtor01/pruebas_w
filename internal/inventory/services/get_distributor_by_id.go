package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetDistributorByIdDto struct {
	Id string
}

type GetDistributorByIdService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetDistributorByIdService(inventoryRepository inventory.RepositoryInventory) *GetDistributorByIdService {
	return &GetDistributorByIdService{inventoryRepository: inventoryRepository}
}

func (s GetDistributorByIdService) Handler(ctx context.Context, dto GetDistributorBySmarkiaDto) (inventory.Distributor, error) {
	ds, err := s.inventoryRepository.GetDistributorByID(ctx, dto.Id)

	return ds, err
}
