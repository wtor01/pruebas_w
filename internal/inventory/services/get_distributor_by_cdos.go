package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetDistributorByCdosDto struct {
	Cdos string
}

type GetDistributorByCdosService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetDistributorByCdosService(repositoryInventory inventory.RepositoryInventory) *GetDistributorByCdosService {
	return &GetDistributorByCdosService{
		inventoryRepository: repositoryInventory,
	}
}

func (s GetDistributorByCdosService) Handler(ctx context.Context, dto GetDistributorByCdosDto) (inventory.Distributor, error) {
	d, err := s.inventoryRepository.GetDistributorByCdos(ctx, dto.Cdos)

	return d, err
}
