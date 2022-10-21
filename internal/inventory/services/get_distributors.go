package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

type GetDistributorsDto struct {
	Q           string
	Limit       int
	Offset      *int
	Sort        map[string]string
	CurrentUser auth.User
	Values      map[string]string
}

type GetDistributorsService struct {
	inventoryRepository inventory.RepositoryInventory
}

func NewGetDistributorsService(inventoryRepository inventory.RepositoryInventory) *GetDistributorsService {
	return &GetDistributorsService{inventoryRepository: inventoryRepository}
}

func (s GetDistributorsService) Handler(ctx context.Context, dto GetDistributorsDto) ([]inventory.Distributor, int, error) {
	ds, count, err := s.inventoryRepository.ListDistributors(ctx, inventory.Search{
		Q:           dto.Q,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
		Sort:        dto.Sort,
		CurrentUser: dto.CurrentUser,
		Values:      dto.Values,
	})

	return ds, count, err
}
