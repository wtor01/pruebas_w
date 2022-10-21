package inventory

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	servicesInventory "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	"bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Inventory struct {
	inventoryServices *servicesInventory.InventoryServices
	redisCache        *redis.DataCacheRedis
}

func (client Inventory) GetMeasureEquipmentBySmarkiaId(ctx context.Context, smarikiaId string) (clients.Equipment, error) {

	var result clients.Equipment
	keyEquipmentRedis := fmt.Sprintf("equipment:%v:%s", "smarkia", smarikiaId)

	if client.redisCache != nil {
		err := client.redisCache.Get(ctx, &result, keyEquipmentRedis)
		if err == nil && result.ID == "" {
			return result, errors.New("record not found " + smarikiaId)
		}
		if err == nil {
			return result, err
		}
	}

	eq, err := client.inventoryServices.GetMeasureEquipmentBySmarkiaIdService.Handler(ctx, servicesInventory.GetMeasureEquipmentBySmarkiaDto{
		Id: smarikiaId,
	})

	if err != nil {
		// guardamos vacio cuando hay error
		_ = client.redisCache.SetWithExpiration(ctx, keyEquipmentRedis, result, time.Hour*6)

		return clients.Equipment{}, err
	}

	if client.redisCache != nil {
		err = client.redisCache.Set(ctx, keyEquipmentRedis, eq)
	}

	resultData := clients.Equipment{
		ID:                eq.ID.String(),
		Type:              eq.Type,
		SerialNumber:      eq.SerialNumber,
		Technology:        eq.Technology,
		Brand:             eq.Brand,
		Model:             eq.Model,
		ActiveConstant:    eq.ActiveConstant,
		ReactiveConstant:  eq.ReactiveConstant,
		MaximeterConstant: eq.MaximeterConstant,
	}
	return resultData, nil
}

func (client Inventory) ListMeasureEquipmentByDistributorId(ctx context.Context, distributorID string) ([]clients.Equipment, error) {
	eq, _, err := client.inventoryServices.GetMeasureEquipmentsService.Handler(ctx, servicesInventory.GetMeasureEquipmentDto{
		Limit:         -1,
		DistributorID: distributorID,
	})

	if err != nil {
		return []clients.Equipment{}, err
	}
	result := make([]clients.Equipment, 0, cap(eq))

	for _, e := range eq {
		result = append(result, clients.Equipment{
			ID:                e.ID.String(),
			Type:              e.Type,
			SerialNumber:      e.SerialNumber,
			Technology:        e.Technology,
			Brand:             e.Brand,
			Model:             e.Model,
			ActiveConstant:    e.ActiveConstant,
			ReactiveConstant:  e.ReactiveConstant,
			MaximeterConstant: e.MaximeterConstant,
		})
	}

	return result, nil
}

func NewInventory(inventoryServices *servicesInventory.InventoryServices, redis *redis.DataCacheRedis) *Inventory {
	return &Inventory{inventoryServices: inventoryServices, redisCache: redis}
}

func (client Inventory) GetAllDistributors(ctx context.Context) ([]clients.Distributor, error) {

	n := 0

	distributors, _, err := client.inventoryServices.GetDistributorsService.Handler(ctx, servicesInventory.GetDistributorsDto{Limit: 1000, Q: "", Offset: &n, Sort: map[string]string{"name": "asc"}, CurrentUser: auth.User{}})

	if err != nil {
		return []clients.Distributor{}, err
	}
	ds := make([]clients.Distributor, 0)

	for _, d := range distributors {
		ds = append(ds, clients.Distributor{
			ID:        d.ID.String(),
			Name:      d.Name,
			SmarkiaID: d.SmarkiaId,
			CDOS:      d.CDOS,
		})
	}
	return ds, nil
}

func (client Inventory) ListDistributors(ctx context.Context, dto clients.ListDistributorsDto) ([]clients.Distributor, int, error) {

	distributors, count, err := client.inventoryServices.GetDistributorsService.Handler(ctx, servicesInventory.GetDistributorsDto{
		Limit:  dto.Limit,
		Offset: dto.Offset,
		Sort:   dto.Sort,
		Values: dto.Values,
	})

	if err != nil {
		return []clients.Distributor{}, count, err
	}
	ds := make([]clients.Distributor, 0)

	for _, d := range distributors {
		ds = append(ds, clients.Distributor{
			ID:        d.ID.String(),
			Name:      d.Name,
			SmarkiaID: d.SmarkiaId,
			CDOS:      d.CDOS,
		})
	}
	return ds, count, nil
}

func (client Inventory) GetDistributorBySmarkiaId(ctx context.Context, smarikiaId string) (clients.Distributor, error) {

	var resultRedis clients.Distributor
	keyDistributorRedis := fmt.Sprintf("distributor:%v:%s", "smarkia", smarikiaId)
	if client.redisCache != nil {
		err := client.redisCache.Get(ctx, &resultRedis, keyDistributorRedis)
		if err == nil {
			return resultRedis, err
		}
	}

	d, err := client.inventoryServices.GetDistributorBySmarkiaIdService.Handler(ctx, servicesInventory.GetDistributorBySmarkiaDto{
		Id: smarikiaId,
	})

	if err != nil {
		return clients.Distributor{}, err
	}

	result := clients.Distributor{
		ID:        d.ID.String(),
		Name:      d.Name,
		SmarkiaID: d.SmarkiaId,
	}

	if client.redisCache != nil {
		err = client.redisCache.Set(ctx, keyDistributorRedis, result)
		if err != nil {
			fmt.Println(err)
			return result, err
		}
	}
	return result, nil
}

func (client Inventory) GetDistributorById(ctx context.Context, id string) (clients.Distributor, error) {

	d, err := client.inventoryServices.GetDistributorByIdService.Handler(ctx, servicesInventory.GetDistributorBySmarkiaDto{
		Id: id,
	})

	if err != nil {
		return clients.Distributor{}, err
	}

	return clients.Distributor{
		ID:        d.ID.String(),
		Name:      d.Name,
		CDOS:      d.CDOS,
		SmarkiaID: d.SmarkiaId,
	}, nil
}

func (client Inventory) GetDistributorByCdos(ctx context.Context, cdos string) (clients.Distributor, error) {
	span := trace.SpanFromContext(ctx)

	span.AddEvent("GetDistributorByCdos")
	var d clients.Distributor

	distributor, err := client.inventoryServices.GetDistributorByCdosService.Handler(ctx, servicesInventory.GetDistributorByCdosDto{
		Cdos: cdos,
	})

	if err != nil {
		return clients.Distributor{}, err
	}

	d = clients.Distributor{
		ID:        distributor.ID.String(),
		Name:      distributor.Name,
		SmarkiaID: distributor.SmarkiaId,
	}

	return d, nil
}

func (client Inventory) GroupMetersByType(ctx context.Context, id string) (map[string]int, error) {
	d, err := client.inventoryServices.GroupByMetersByDistributorIdService.Handler(ctx, servicesInventory.GroupByMetersByDistributorIdDto{DistributorID: id})

	if err != nil {
		return map[string]int{}, err
	}

	return d, nil
}

func (client Inventory) GetServicePoints(ctx context.Context, inventoryClient clients.ServicePointSchedulerDto) ([]clients.ServicePointScheduler, error) {
	s, err := client.inventoryServices.GetServicePointService.Handler(ctx, inventory.ServicePointSchedulerDto(inventoryClient))
	return utils.MapSlice(s, func(s inventory.ServicePointScheduler) clients.ServicePointScheduler {
		return clients.ServicePointScheduler{
			MeterConfigId:          s.MeterConfigId,
			MeterId:                s.MeterId,
			DistributorId:          s.DistributorId,
			DistributorCode:        s.DistributorCode,
			Cups:                   s.Cups,
			MeterSerialNumber:      s.MeterSerialNumber,
			ServiceType:            s.ServiceType,
			ReadingType:            s.ReadingType,
			PointType:              s.PointType,
			MeasurePointType:       s.MeasurePointType,
			PriorityContract:       s.PriorityContract,
			TlgCode:                s.TlgCode,
			TlgType:                s.TlgType,
			Calendar:               s.Calendar,
			Type:                   s.Type,
			StartDate:              s.StartDate,
			EndDate:                s.EndDate,
			CurveType:              s.CurveType,
			Technology:             s.Technology,
			ContractualSituationId: s.ContractualSituationId,
			TariffID:               s.TariffID,
			Coefficient:            s.Coefficient,
			CalendarCode:           s.CalendarCode,
			P1Demand:               s.P1Demand,
			P2Demand:               s.P2Demand,
			P3Demand:               s.P3Demand,
			P4Demand:               s.P4Demand,
			P5Demand:               s.P5Demand,
			P6Demand:               s.P6Demand,
			AI:                     s.AI,
			AE:                     s.AE,
			R1:                     s.R1,
			R2:                     s.R2,
			R3:                     s.R3,
			R4:                     s.R4,
			ContractStartDate:      s.ContractStartDate,
			ContractEndDate:        s.ContractEndDate,
			ServicePointID:         s.ServicePointID,
			GeographicID:           s.GeographicID,
		}
	}), err
}

func (client Inventory) GetServicePointConfigByMeter(ctx context.Context, query clients.GetServicePointConfigByMeterIdDto) (clients.ServicePointScheduler, error) {
	s, err := client.inventoryServices.GetServicePointByMeterService.Handler(ctx, inventory.GetMeasureEquipmentByMeterQuery{
		MeterID:           query.MeterID,
		MeterSerialNumber: query.MeterSerialNumber,
		Time:              query.Time,
	})

	return clients.ServicePointScheduler(s), err
}

func (client Inventory) GetMeterConfigByMeterService(ctx context.Context, query clients.GetMeterConfigByMeterDto) (clients.MeterConfig, error) {
	s, err := client.inventoryServices.GetMeterConfigByMeterService.Handler(ctx, servicesInventory.GetMeterConfigByMeterDto{
		MeterSerialNumber: query.MeterSerialNumber,
		Date:              query.Date,
	})

	return client.meterConfigToResponse(s), err
}

func (client Inventory) meterConfigToResponse(c inventory.MeterConfig) clients.MeterConfig {
	return clients.MeterConfig{
		ID:               c.ID,
		StartDate:        c.StartDate,
		EndDate:          c.EndDate,
		AI:               c.AI,
		AE:               c.AE,
		R1:               c.R1,
		R2:               c.R2,
		R3:               c.R3,
		R4:               c.R4,
		M:                c.M,
		E:                c.E,
		CurveType:        c.CurveType,
		ReadingType:      c.ReadingType,
		CalendarID:       c.CalendarID,
		PriorityContract: c.PriorityContract,
		MeterID:          c.MeterID,
		MeasurePointID:   c.MeasurePointID,
		ServicePointID:   c.ServicePointID,
		TlgCode:          string(c.TlgCode),
		Type:             c.Type,
		RentingPrice:     c.RentingPrice,
		SerialNumber:     c.SerialNumber,
		DistributorID:    c.DistributorID,
		CUPS:             c.CUPS,
		MeasurePointType: string(c.MeasurePointType),
	}
}
