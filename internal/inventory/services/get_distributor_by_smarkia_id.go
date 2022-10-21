package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetDistributorBySmarkiaDto struct {
	Id string
}

type GetDistributorBySmarkiaIdService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetDistributorBySmarkiaIdService(inventoryRepository inventory.RepositoryInventory) *GetDistributorBySmarkiaIdService {
	return &GetDistributorBySmarkiaIdService{inventoryRepository: inventoryRepository}
}

func (s GetDistributorBySmarkiaIdService) Handler(ctx context.Context, dto GetDistributorBySmarkiaDto) (inventory.Distributor, error) {
	ds, err := s.inventoryRepository.GetDistributorBySmarkiaID(ctx, dto.Id)

	return ds, err
}
