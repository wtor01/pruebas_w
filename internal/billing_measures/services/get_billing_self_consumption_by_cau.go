package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sort"
	"strings"
	"time"
)

type GetBillingSelfConsumptionDto struct {
	CauId         string    `json:"cau_id"`
	DistributorId string    `json:"distributor_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

func NewGetBillingSelfConsumptionDto(cauId, distributorId string, startDate, endDate time.Time) GetBillingSelfConsumptionDto {
	return GetBillingSelfConsumptionDto{
		CauId:         cauId,
		DistributorId: distributorId,
		StartDate:     startDate,
		EndDate:       endDate.AddDate(0, 0, 1),
	}
}

func (dto GetBillingSelfConsumptionDto) String() string {
	dtoJson, _ := json.Marshal(dto)
	return string(dtoJson)
}

type GetBillingSelfConsumptionByCauService struct {
	selfConsumptionRepository billing_measures.BillingSelfConsumptionRepository
	location                  *time.Location
	tracer                    trace.Tracer
}

func NewGetBillingSelfConsumptionByCauService(selfConsumptionRepository billing_measures.BillingSelfConsumptionRepository, loc *time.Location) *GetBillingSelfConsumptionByCauService {
	return &GetBillingSelfConsumptionByCauService{
		selfConsumptionRepository: selfConsumptionRepository,
		location:                  loc,
		tracer:                    telemetry.GetTracer(),
	}
}

func (s GetBillingSelfConsumptionByCauService) Handler(ctx context.Context, dto GetBillingSelfConsumptionDto) ([]billing_measures.BillingSelfConsumptionUnit, error) {
	ctx, span := s.tracer.Start(ctx, "GetBillingSelfConsumptionByCauService - Handler")
	defer span.End()

	span.SetAttributes(attribute.String("dto", dto.String()))
	billingSelfConsumption, err := s.selfConsumptionRepository.GetSelfConsumptionByCau(ctx, billing_measures.QueryGetBillingSelfConsumptionByCau{
		CauId:         dto.CauId,
		DistributorId: dto.DistributorId,
		StartDate:     dto.StartDate,
		EndDate:       dto.EndDate,
	})

	if err != nil {
		return []billing_measures.BillingSelfConsumptionUnit{}, err
	}

	return utils.MapSlice(billingSelfConsumption, s.toResponse), nil
}

func (s GetBillingSelfConsumptionByCauService) toResponse(consumption billing_measures.BillingSelfConsumption) billing_measures.BillingSelfConsumptionUnit {
	selfConsumptionUnit := billing_measures.NewBillingSelfConsumptionUnit(consumption)

	cupsInfoMap := make(map[string]billing_measures.BillingSelfConsumptionPoint)

	for _, point := range consumption.Points {
		s.fillCupsMapInfo(point, cupsInfoMap)
	}

	netGenerationMap := make(map[string]*billing_measures.BillingSelfConsumptionNetGeneration)
	unitConsumptionMap := make(map[string]*billing_measures.BillingSelfConsumptionUnitConsumption)
	calendarConsumptionMap := make(map[string]*billing_measures.BillingSelfConsumptionCalendarConsumption)
	cupsMap := make(map[string]*billing_measures.BillingSelfConsumptionCups)

	for _, curve := range consumption.Curve {
		date, hour := s.getFormattedDate(curve.EndDate)
		for _, curvePoint := range curve.Points {
			s.setCupsMap(curvePoint, cupsInfoMap[curvePoint.CUPS], cupsMap)
		}

		s.setNetGenerationMap(date, curve.EHGN, curve.EHEX, netGenerationMap)
		s.setUnitConsumptionMap(date, curve.EHCR, curve.EHAU, curve.EHSA, unitConsumptionMap)
		s.setCalendarConsumptionMap(date, hour, curve.EHEX, calendarConsumptionMap)
		selfConsumptionUnit.SumTotals(billing_measures.BillingSelfConsumptionTotals{
			GrossGeneration:    utils.Float64(curve.EHGN) + utils.Float64(curve.EHSA),
			NetGeneration:      utils.Float64(curve.EHGN),
			SelfConsumption:    utils.Float64(curve.EHAU),
			NetworkConsumption: utils.Float64(curve.EHCR),
			AuxConsumption:     utils.Float64(curve.EHSA),
		})
	}

	netGeneration := s.transformNetGenerationMap(netGenerationMap)
	selfConsumptionUnit.SetNetGeneration(netGeneration)

	unitConsumption := s.transformUnitConsumptionMap(unitConsumptionMap)
	selfConsumptionUnit.SetUnitConsumption(unitConsumption)

	calendarConsumption := s.transformCalendarConsumptionMap(calendarConsumptionMap)
	selfConsumptionUnit.SetCalendarConsumption(calendarConsumption)

	cupsList := s.transformCupsMap(cupsMap)
	selfConsumptionUnit.SetCups(cupsList)

	return selfConsumptionUnit
}

func (s GetBillingSelfConsumptionByCauService) fillCupsMapInfo(info billing_measures.BillingSelfConsumptionPoint, cupsMap map[string]billing_measures.BillingSelfConsumptionPoint) {
	cupsMap[info.CUPS] = info
}

func (s GetBillingSelfConsumptionByCauService) getFormattedDate(date time.Time) (string, string) {
	dateTime := date.In(s.location).Format("2006-01-02 15:04")

	dateSplitted := strings.Split(dateTime, " ")

	return dateSplitted[0], dateSplitted[1]
}

func (s GetBillingSelfConsumptionByCauService) setNetGenerationMap(date string, ehgn, ehex *float64, netGenerationMap map[string]*billing_measures.BillingSelfConsumptionNetGeneration) {
	mapValue, ok := netGenerationMap[date]
	var ehgnValue, ehexValue = utils.Float64(ehgn), utils.Float64(ehex)
	if !ok {
		netGenerationMap[date] = &billing_measures.BillingSelfConsumptionNetGeneration{
			Date:      date,
			Net:       ehgnValue,
			Excedente: ehexValue,
		}
		return
	}

	mapValue.Sum(ehgnValue, ehexValue)
}

func (s GetBillingSelfConsumptionByCauService) setUnitConsumptionMap(date string, ehcr, ehau, ehsa *float64, unitConsumptionMap map[string]*billing_measures.BillingSelfConsumptionUnitConsumption) {
	mapValue, ok := unitConsumptionMap[date]

	var ehcrValue, ehauValue, ehsaValue = utils.Float64(ehcr), utils.Float64(ehau), utils.Float64(ehsa)

	if !ok {
		unitConsumptionMap[date] = &billing_measures.BillingSelfConsumptionUnitConsumption{
			Date:    date,
			Network: ehcrValue,
			Self:    ehauValue,
			Aux:     ehsaValue,
		}
		return
	}

	mapValue.Sum(ehcrValue, ehauValue, ehsaValue)
}

func (s GetBillingSelfConsumptionByCauService) setCalendarConsumptionMap(date, hour string, ehex *float64, calendarConsumptionMap map[string]*billing_measures.BillingSelfConsumptionCalendarConsumption) {
	mapValue, ok := calendarConsumptionMap[date]

	var ehexValue = utils.Float64(ehex)
	if !ok {
		hourValues := make([]billing_measures.CalendarConsumptionValues, 0, 24)
		hourValues = append(hourValues, billing_measures.CalendarConsumptionValues{
			Hour:   hour,
			Energy: ehexValue,
		})
		calendarConsumptionMap[date] = &billing_measures.BillingSelfConsumptionCalendarConsumption{
			Date:   date,
			Energy: ehexValue,
			Values: hourValues,
		}
		return
	}

	mapValue.AddValue(billing_measures.CalendarConsumptionValues{
		Hour:   hour,
		Energy: ehexValue,
	})
	mapValue.Sum(ehexValue)

}

func (s GetBillingSelfConsumptionByCauService) setCupsMap(curvePoint billing_measures.BillingSelfConsumptionCurvePoint, pointInfo billing_measures.BillingSelfConsumptionPoint, cupsMap map[string]*billing_measures.BillingSelfConsumptionCups) {

	cups := curvePoint.CUPS
	cupsValue, ok := cupsMap[cups]

	layout := "02-01-2006"
	if !ok {
		cupsMap[cups] = &billing_measures.BillingSelfConsumptionCups{
			Cups:        cups,
			Type:        pointInfo.ServicePointType,
			Consumption: curvePoint.AI,
			Generation:  curvePoint.AE,
			StartDate:   pointInfo.InitDate.In(s.location).Format(layout),
			EndDate:     pointInfo.EndDate.In(s.location).Format(layout),
		}
		return
	}

	cupsValue.Sum(curvePoint.AI, curvePoint.AE)
}

func transformMap[T any](items map[string]*T) []T {
	arr := make([]T, 0, len(items))
	for _, value := range items {
		arr = append(arr, *value)
	}

	return arr
}

func (s GetBillingSelfConsumptionByCauService) transformNetGenerationMap(netGenerationMap map[string]*billing_measures.BillingSelfConsumptionNetGeneration) []billing_measures.BillingSelfConsumptionNetGeneration {
	arr := transformMap(netGenerationMap)

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Date < arr[j].Date
	})

	return arr
}

func (s GetBillingSelfConsumptionByCauService) transformUnitConsumptionMap(unitConsumptionMap map[string]*billing_measures.BillingSelfConsumptionUnitConsumption) []billing_measures.BillingSelfConsumptionUnitConsumption {
	arr := transformMap(unitConsumptionMap)

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Date < arr[j].Date
	})

	return arr
}

func (s GetBillingSelfConsumptionByCauService) transformCalendarConsumptionMap(calendarConsumptionMap map[string]*billing_measures.BillingSelfConsumptionCalendarConsumption) []billing_measures.BillingSelfConsumptionCalendarConsumption {
	arr := transformMap(calendarConsumptionMap)

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Date < arr[j].Date
	})

	return arr
}

func (s GetBillingSelfConsumptionByCauService) transformCupsMap(cupsMap map[string]*billing_measures.BillingSelfConsumptionCups) []billing_measures.BillingSelfConsumptionCups {
	arr := transformMap(cupsMap)

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Cups < arr[j].Cups
	})

	return arr
}
