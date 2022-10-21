package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type DeleteFestiveDay struct {
	repository festive_days.FestiveDayRepository
	publisher  event.PublisherCreator
	topic      string
}

func NewDeleteFestiveDay(repository festive_days.FestiveDayRepository, publisher event.PublisherCreator, topic string) *DeleteFestiveDay {
	return &DeleteFestiveDay{
		repository: repository,
		publisher:  publisher,
		topic:      topic,
	}
}

func (d DeleteFestiveDay) Handler(ctx context.Context, festiveDayId string) error {
	err := d.repository.DeleteFestiveDay(ctx, festiveDayId)

	if err != nil {
		return err
	}

	err = event.PublishEvent(ctx, d.topic, d.publisher, master_tables.NewFestiveDaysGenerateEvent())
	return err
}
