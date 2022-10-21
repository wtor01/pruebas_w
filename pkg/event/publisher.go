package event

import (
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"sync"
	"sync/atomic"
)

//go:generate mockery --case=snake --outpkg=mocks --output=./mocks --name=Publisher
type Publisher interface {
	Publish(ctx context.Context, topic string, message []byte, attributes map[string]string) error
	Close() error
}

type PublisherCreator = func(ctx context.Context) (Publisher, error)

func PublishAllEvents[P any](ctx context.Context, topic string, pubCreator PublisherCreator, events []Message[P]) error {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "PublishAllEvents")
	defer span.End()

	var wg sync.WaitGroup

	totalErrors := int64(len(events))
	wg.Add(len(events))
	publisher, err := pubCreator(ctx)

	if err != nil {
		return err
	}

	for _, ev := range events {
		go func(ctx context.Context, publisher Publisher, event Message[P], wg *sync.WaitGroup) {
			defer wg.Done()
			tr := telemetry.GetTracer()
			ctx, spanChild := tr.Start(ctx, "publish message "+event.Type)
			defer spanChild.End()

			msg, err := event.Marshal()
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				return
			}
			err = publisher.Publish(ctx, topic, msg, event.GetAttributes())
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				return
			}
			atomic.AddInt64(&totalErrors, -1)
		}(ctx, publisher, ev, &wg)
	}

	wg.Wait()

	defer publisher.Close()

	if totalErrors != 0 {
		return errors.New(fmt.Sprintf("number errors %d", totalErrors))
	}

	return nil
}

func PublishEvent[P any](ctx context.Context, topic string, pubCreator PublisherCreator, event Message[P]) error {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "PublishEvent")
	defer span.End()

	var err error

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	publisher, err := pubCreator(ctx)

	defer publisher.Close()

	if err != nil {
		return err
	}

	msg, err := event.Marshal()

	if err != nil {
		return err
	}
	err = publisher.Publish(ctx, topic, msg, event.GetAttributes())

	if err != nil {
		return err
	}

	return nil
}
