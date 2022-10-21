package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"time"
)

type GetMeterConfigByCupsDto struct {
	Cups        string
	Distributor string
	Date        time.Time
}
type GetMeterConfigByCupsService struct {
	mongoRepository measures.InventoryRepository
}

func NewGetMeterConfigByCupsService(mongoRepository measures.InventoryRepository) *GetMeterConfigByCupsService {
	return &GetMeterConfigByCupsService{
		mongoRepository: mongoRepository,
	}
}
func (s GetMeterConfigByCupsService) Handler(ctx context.Context, dto GetMeterConfigByCupsDto) (measures.MeterConfig, error) {
	return s.mongoRepository.GetMeterConfigByCupsAPI(
		ctx,
		measures.GetMeterConfigByCupsQuery{
			CUPS:        dto.Cups,
			Time:        dto.Date,
			Distributor: dto.Distributor,
		},
	)
}
