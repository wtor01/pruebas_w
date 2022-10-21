package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

type InventoryServices struct {
	GetMeasureEquipmentsService           *GetMeasureEquipmentService
	GetDistributorsService                *GetDistributorsService
	GetDistributorBySmarkiaIdService      *GetDistributorBySmarkiaIdService
	GetDistributorByIdService             *GetDistributorByIdService
	GetMeasureEquipmentByIdService        *GetMeasureEquipmentByIdService
	GetMeasureEquipmentBySmarkiaIdService *GetMeasureEquipmentBySmarkiaIdService
	GroupByMetersByDistributorIdService   *GroupByMetersByDistributorIdService
	GetServicePointService                *GetServicePointService
	GetServicePointByMeterService         *GetServicePointByMeterService
	GetDistributorByCdosService           *GetDistributorByCdosService
	GetMeterConfigByMeterService          *GetMeterConfigByMeterService
	GetMeterConfigByCupsService           *GetMeterConfigByCupsService
}

func New(inventoryRepository inventory.RepositoryInventory, mongoRepository measures.InventoryRepository) *InventoryServices {
	getMeasureEquipmentsService := NewGetMeasureEquipmentService(inventoryRepository)
	getDistributorsService := NewGetDistributorsService(inventoryRepository)
	getDistributorBySmarkiaIdService := NewGetDistributorBySmarkiaIdService(inventoryRepository)
	getDistributorByIdService := NewGetDistributorByIdService(inventoryRepository)
	getMeasureEquipmentByIdService := NewGetMeasureEquipmentByIdService(inventoryRepository)
	getMeasureEquipmentBySmarkiaIdService := NewGetMeasureEquipmentBySmarkiaIdService(inventoryRepository)
	groupByMetersByDistributorIdService := NewGroupByMetersByDistributorIdService(inventoryRepository)
	getServicePointService := NewServicePointService(inventoryRepository)
	getServicePointByMeterService := NewGetServicePointByMeterService(inventoryRepository)
	getDistributorByCdosService := NewGetDistributorByCdosService(inventoryRepository)
	getMeterConfigByMeterService := NewGetMeterConfigByMeterService(inventoryRepository)
	getMeterConfigByCupsService := NewGetMeterConfigByCupsService(mongoRepository)

	return &InventoryServices{
		GetMeasureEquipmentsService:           getMeasureEquipmentsService,
		GetDistributorsService:                getDistributorsService,
		GetDistributorBySmarkiaIdService:      getDistributorBySmarkiaIdService,
		GetDistributorByIdService:             getDistributorByIdService,
		GetMeasureEquipmentByIdService:        getMeasureEquipmentByIdService,
		GetMeasureEquipmentBySmarkiaIdService: getMeasureEquipmentBySmarkiaIdService,
		GroupByMetersByDistributorIdService:   groupByMetersByDistributorIdService,
		GetServicePointService:                getServicePointService,
		GetServicePointByMeterService:         getServicePointByMeterService,
		GetDistributorByCdosService:           getDistributorByCdosService,
		GetMeterConfigByMeterService:          getMeterConfigByMeterService,
		GetMeterConfigByCupsService:           getMeterConfigByCupsService,
	}
}
