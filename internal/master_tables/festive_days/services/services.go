package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
)

type FestiveDaysServices struct {
	GetListFestiveDays *GetListFestiveDays
	GetFestiveDayById  *GetFestiveDayById
	SaveFestiveDay     *SaveFestiveDay
	UpdateFestiveDay   *UpdateFestiveDay
	DeleteFestiveDay   *DeleteFestiveDay
}

func New(repository festive_days.FestiveDayRepository, publisher event.PublisherCreator, topic string) *FestiveDaysServices {
	getListFestiveDays := NewGetListFestiveDaysService(repository)
	getFestiveDayById := NewGetFestiveDayByIdService(repository)
	saveFestiveDay := NewSaveFestiveDay(repository, publisher, topic)
	updateFestiveDay := NewUpdateFestiveDay(repository, publisher, topic)
	deleteFestiveDay := NewDeleteFestiveDay(repository, publisher, topic)

	return &FestiveDaysServices{
		GetListFestiveDays: getListFestiveDays,
		GetFestiveDayById:  getFestiveDayById,
		SaveFestiveDay:     saveFestiveDay,
		UpdateFestiveDay:   updateFestiveDay,
		DeleteFestiveDay:   deleteFestiveDay,
	}
}
