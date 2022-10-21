package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"time"
)

type UpdateFestiveDayDto struct {
	Id           string
	Date         time.Time
	Description  string
	GeographicId string
}

type UpdateFestiveDay struct {
	repository festive_days.FestiveDayRepository
	publisher  event.PublisherCreator
	topic      string
}

func NewUpdateFestiveDay(repository festive_days.FestiveDayRepository, publisher event.PublisherCreator, topic string) *UpdateFestiveDay {
	return &UpdateFestiveDay{
		repository: repository,
		publisher:  publisher,
		topic:      topic,
	}
}

func (s UpdateFestiveDay) Handler(ctx context.Context, festiveDayId string, dto UpdateFestiveDayDto) error {

	err := s.repository.UpdateFestiveDay(ctx, festiveDayId, festive_days.FestiveDay{
		Date:         dto.Date,
		Description:  dto.Description,
		GeographicId: dto.GeographicId,
	})

	if err != nil {
		return err
	}

	err = event.PublishEvent(ctx, s.topic, s.publisher, master_tables.NewFestiveDaysGenerateEvent())
	return err
}
