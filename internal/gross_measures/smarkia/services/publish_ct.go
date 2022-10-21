package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"time"
)

type PublishCtDTO struct {
	ProcessName     string
	DistributorId   string
	DistributorCDOS string
	SmarkiaId       string
	Date            time.Time
}

type PublishCtService struct {
	api       smarkia.GetCTer
	publisher event.PublisherCreator
	topic     string
}

func NewPublishCtService(publisher event.PublisherCreator, getCTer smarkia.GetCTer, topic string) *PublishCtService {
	return &PublishCtService{api: getCTer, publisher: publisher, topic: topic}
}

func (svc PublishCtService) Handle(ctx context.Context, dto PublishCtDTO) error {
	cts, err := svc.api.GetCTs(ctx, dto.SmarkiaId)
	if err != nil {
		return err
	}

	messages := make([]smarkia.CtProcessEvent, 0, cap(cts))

	for _, ct := range cts {
		messages = append(messages, smarkia.NewCtProcessEvent(dto.DistributorId, dto.SmarkiaId, dto.ProcessName, dto.DistributorCDOS, ct.ID, dto.Date))
	}

	err = event.PublishAllEvents(ctx, svc.topic, svc.publisher, messages)

	return err
}
