package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
	"context"
)

//GetTariffByCodeDto struct send code for Get
type GetTariffByCodeDto struct {
	ID string
}
type GetOneTariffService struct {
	TariffRepository tariff.RepositoryTariff
}

//NewGetOneTariffService create get one struct
func NewGetOneTariffService(geographicRepository tariff.RepositoryTariff) GetOneTariffService {
	return GetOneTariffService{TariffRepository: geographicRepository}
}

//Handler handle get one geographic service
func (s GetOneTariffService) Handler(ctx context.Context, dto GetTariffByCodeDto) (tariff.Tariffs, error) {
	res, err := s.TariffRepository.GetTariff(ctx, dto.ID)

	return res, err
}
