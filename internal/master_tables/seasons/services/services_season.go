package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
)

type SeasonsServices struct {
	GetSeasons     *GetSeasonsService
	GetSeasonById  *GetSeasonByIdService
	InsertSeason   *InsertSeasonService
	ModifySeason   *ModifySeasonService
	DeleteSeason   *DeleteSeasonService
	GetDayTypes    *GetDayTypesService
	GetDayTypeById *GetDayTypeByIdService
	InsertDayType  *InsertDayTypeService
	ModifyDayType  *ModifyDayTypeService
	DeleteDayType  *DeleteDayTypeService
}

func NewSeasonsServices(repositorySeasons seasons.RepositorySeasons) *SeasonsServices {
	getSeasons := NewGetSeasonsRepositoryService(repositorySeasons)
	getSeasonById := NewGetSeasonByIdService(repositorySeasons)
	insertSeason := NewInsertSeasonService(repositorySeasons)
	modifySeason := NewModifySeasonService(repositorySeasons)
	deleteSeason := NewDeleteSeasonService(repositorySeasons)
	getDayTypes := NewGetDayTypesService(repositorySeasons)
	getDayTypeById := NewGetDayTypeByIdService(repositorySeasons)
	insertDayType := NewInsertDayTypeService(repositorySeasons)
	modifyDayType := NewModifyDayTypeService(repositorySeasons)
	deleteDayType := NewDeleteDayTypeService(repositorySeasons)
	return &SeasonsServices{
		GetSeasons:     &getSeasons,
		GetSeasonById:  &getSeasonById,
		InsertSeason:   &insertSeason,
		ModifySeason:   &modifySeason,
		DeleteSeason:   &deleteSeason,
		GetDayTypes:    &getDayTypes,
		GetDayTypeById: &getDayTypeById,
		InsertDayType:  &insertDayType,
		ModifyDayType:  &modifyDayType,
		DeleteDayType:  &deleteDayType,
	}
}
