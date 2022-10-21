package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"time"
)

type ClosureServices struct {
	GetClosure    *GetClosure
	CreateClosure *CreateClosure
	UpdateClosure *UpdateClosure
}

func NewClosureServices(
	mongoRepository process_measures.ProcessMeasureClosureRepository,
	inventoryRepository measures.InventoryRepository,

	masterTablesClient clients.MasterTables,
	generatorDate func() time.Time,
) *ClosureServices {
	getClosureService := NewGetClosureService(mongoRepository)
	createClosure := NewCreateClosure(mongoRepository, inventoryRepository, masterTablesClient, generatorDate)
	updateClosureService := NewUpdateClosureService(mongoRepository)
	return &ClosureServices{
		GetClosure:    getClosureService,
		CreateClosure: createClosure,
		UpdateClosure: updateClosureService,
	}
}
