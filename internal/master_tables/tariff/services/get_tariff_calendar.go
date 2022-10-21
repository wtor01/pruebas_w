package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
	"context"
)

type GetTariffCalendarService struct {
	TariffRepository tariff.RepositoryTariff
}

// NewGetTariffCalendarService Generate Get distributor service
func NewGetTariffCalendarService(tariffRepository tariff.RepositoryTariff) GetTariffCalendarService {
	return GetTariffCalendarService{TariffRepository: tariffRepository}
}

// Handler Get data, count and errors from db
func (s GetTariffCalendarService) Handler(ctx context.Context, code string, dto GetAllDto) ([]tariff.TariffsCalendar, int, error) {
	ds, count, err := s.TariffRepository.GetAllTariffCalendar(ctx, code, tariff.Search{
		Q:           dto.Q,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
		Sort:        dto.Sort,
		CurrentUser: dto.CurrentUser,
	})

	return ds, count, err
}
