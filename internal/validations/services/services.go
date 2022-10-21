package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
)

type Services struct {
	CreateValidationMeasureService       *CreateValidationMeasureService
	ListValidationMeasureService         *ListValidationMeasureService
	GetValidationMeasureByIdService      *GetValidationMeasureByIdService
	UpdateValidationMeasureService       *UpdateValidationMeasureService
	DeleteValidationMeasureByIdService   *DeleteValidationMeasureByIdService
	CreateValidationMeasureConfigService *CreateValidationMeasureConfigService
	GetValidationMeasureConfigService    *GetValidationMeasureConfigService
	DeleteValidationMeasureConfigService *DeleteValidationMeasureConfigService
	ListValidationMeasureConfigService   *ListValidationMeasureConfigService
	PutMeasureValidation                 *PutValidationMeasure
}

func NewServices(repository validations.ValidationsRepository, repositoryMeasure process_measures.ProcessedMeasureRepository) *Services {
	return &Services{
		CreateValidationMeasureService:       NewCreateValidationMeasureService(repository),
		ListValidationMeasureService:         NewListValidationMeasureService(repository),
		GetValidationMeasureByIdService:      NewGetValidationMeasureByIdService(repository),
		UpdateValidationMeasureService:       NewUpdateValidationMeasureService(repository),
		DeleteValidationMeasureByIdService:   NewDeleteValidationMeasureByIdService(repository),
		CreateValidationMeasureConfigService: NewCreateValidationMeasureConfigService(repository),
		GetValidationMeasureConfigService:    NewGetValidationMeasureConfigService(repository),
		DeleteValidationMeasureConfigService: NewDeleteValidationMeasureConfigService(repository),
		ListValidationMeasureConfigService:   NewListValidationMeasureConfigService(repository),
		PutMeasureValidation:                 NewPutValidationMeasure(repositoryMeasure),
	}
}
