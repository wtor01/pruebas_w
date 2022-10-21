package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"time"
)

type SaveFestiveDayDto struct {
	Id           string
	Date         time.Time
	Description  string
	GeographicId string
	CreatedBy    string
	UpdatedBy    string
}

type SaveFestiveDay struct {
	repository festive_days.FestiveDayRepository
	publisher  event.PublisherCreator
	topic      string
}

func NewSaveFestiveDay(repository festive_days.FestiveDayRepository, publisher event.PublisherCreator, topic string) *SaveFestiveDay {
	return &SaveFestiveDay{
		repository: repository,
		publisher:  publisher,
		topic:      topic,
	}
}

func (s SaveFestiveDay) Handler(ctx context.Context, festive festive_days.FestiveDay) error {

	err := s.repository.SaveFestiveDay(ctx, festive)

	if err != nil {
		return err
	}

	err = event.PublishEvent(ctx, s.topic, s.publisher, master_tables.NewFestiveDaysGenerateEvent())
	return err
}
