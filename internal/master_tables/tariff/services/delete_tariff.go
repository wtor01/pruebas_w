package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
	"context"
)

//DeleteTariffByCodeDto dto to get code
type DeleteTariffByCodeDto struct {
	Id string
}
type DeleteTariffService struct {
	TariffRepository tariff.RepositoryTariff
}

//NewDeleteTariffsService create delete service
func NewDeleteTariffsService(tariffRepository tariff.RepositoryTariff) DeleteTariffService {
	return DeleteTariffService{TariffRepository: tariffRepository}
}

//Handler handle delete
func (s DeleteTariffService) Handler(ctx context.Context, dto DeleteTariffByCodeDto) error {
	err := s.TariffRepository.DeleteTariff(ctx, dto.Id)

	return err
}
