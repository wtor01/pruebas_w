package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"time"
)

type DashboardServiceDTO struct {
	DistributorID string
	StartDate     time.Time
	EndDate       time.Time
}

func NewDashboardServiceDTO(distributorID string, startDate time.Time, endDate time.Time) DashboardServiceDTO {
	startDate = startDate.UTC()
	endDate = endDate.AddDate(0, 0, 1).UTC()
	return DashboardServiceDTO{DistributorID: distributorID, StartDate: startDate, EndDate: endDate}
}

type DashboardService struct {
	dashboardRepository gross_measures.DashboardRepository
	inventoryClient     clients.Inventory
	Location            *time.Location
}

func NewDashboardService(dashboardRepository gross_measures.DashboardRepository, inventoryClient clients.Inventory, loc *time.Location) *DashboardService {
	return &DashboardService{dashboardRepository: dashboardRepository, inventoryClient: inventoryClient, Location: loc}
}

func (svc DashboardService) getMeasureShouldBeInfo(ctx context.Context, dto DashboardServiceDTO) (measures.MeasureShouldBe, error) {
	equipments, err := svc.inventoryClient.GroupMetersByType(ctx, dto.DistributorID)

	if err != nil {
		return measures.MeasureShouldBe{}, err
	}

	countEquimentSTG := 0
	countEquimentSTM := 0
	countEquimentOthers := 0

	for t, eq := range equipments {
		if t == clients.EquipmentTelemedidaType {
			countEquimentSTM = eq
		}
		if t == clients.EquipmentTelegestionType {
			countEquimentSTG = eq
		}
		if t == clients.EquipmentOtherType {
			countEquimentOthers = eq
		}
	}

	diff := dto.EndDate.Sub(dto.StartDate)

	days := diff.Hours() / 24

	result := measures.MeasureShouldBe{
		Telegestion: struct {
			Curva   measures.MeasureShouldBeValue
			Closing measures.MeasureShouldBeValue
			Resumen measures.MeasureShouldBeValue
		}{
			Curva: measures.MeasureShouldBeValue{
				Daily: countEquimentSTG * 24 * 60 / 15,
				Total: countEquimentSTG * int(diff.Hours()) * 60 / 15,
			},
			Closing: measures.MeasureShouldBeValue{
				Daily: 0,
				Total: countEquimentSTG,
			},
			Resumen: measures.MeasureShouldBeValue{
				Daily: countEquimentSTG,
				Total: countEquimentSTG * int(days),
			},
		},
		Telemedida: struct {
			Curva   measures.MeasureShouldBeValue
			Closing measures.MeasureShouldBeValue
		}{
			Curva: measures.MeasureShouldBeValue{
				Daily: countEquimentSTM * 24 * 60 / 15,
				Total: countEquimentSTM * int(diff.Hours()) * 60 / 15,
			},
			Closing: measures.MeasureShouldBeValue{
				Daily: 0,
				Total: countEquimentSTM,
			},
		},
		Others: struct {
			Closing measures.MeasureShouldBeValue
		}{
			Closing: measures.MeasureShouldBeValue{
				Daily: 0,
				Total: countEquimentOthers,
			},
		},
	}

	return result, nil
}

func (svc DashboardService) Handle(ctx context.Context, dto DashboardServiceDTO) (measures.DashboardResult, error) {

	dashboardMeasures, err := svc.dashboardRepository.GetDashboard(ctx, gross_measures.GetDashboardQuery{
		DistributorId: dto.DistributorID,
		StartDate:     dto.StartDate,
		EndDate:       dto.EndDate,
	})

	if err != nil {
		return measures.DashboardResult{}, err
	}

	measureShouldBeInfo, err := svc.getMeasureShouldBeInfo(ctx, dto)

	if err != nil {
		return measures.DashboardResult{}, err
	}

	datesMeasures := measures.GetDashboardResultDailyData(dto.StartDate.In(svc.Location), dto.EndDate.In(svc.Location), measureShouldBeInfo)

	result := measures.GetDashboardResult(dashboardMeasures, measureShouldBeInfo, datesMeasures)

	return result, nil
}
