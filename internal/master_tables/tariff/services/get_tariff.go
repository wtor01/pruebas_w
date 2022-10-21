package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
	"context"
)

type GetAllDto struct {
	Q           string
	Limit       int
	Offset      *int
	Sort        map[string]string
	CurrentUser auth.User
}

type GetTariffsService struct {
	TariffRepository tariff.RepositoryTariff
}

// NewGetTariffRepositoryService Generate Get distributor service
func NewGetTariffRepositoryService(tariffRepository tariff.RepositoryTariff) GetTariffsService {
	return GetTariffsService{TariffRepository: tariffRepository}
}

// Handler Get data, count and errors from db
func (s GetTariffsService) Handler(ctx context.Context, dto GetAllDto) ([]tariff.Tariffs, int, error) {
	ds, count, err := s.TariffRepository.GetAllTariffs(ctx, tariff.Search{
		Q:           dto.Q,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
		Sort:        dto.Sort,
		CurrentUser: dto.CurrentUser,
	})

	return ds, count, err
}
