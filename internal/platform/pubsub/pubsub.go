package pubsub

import (
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"cloud.google.com/go/pubsub"
	"context"
)

func PublisherCreatorGcp(projectID string) event.PublisherCreator {
	return func(ctx context.Context) (event.Publisher, error) {
		client, err := pubsub.NewClient(ctx, projectID)

		return PubSubGcp{
			projectID: projectID,
			client:    client,
		}, err
	}
}

type PubSubGcp struct {
	projectID string
	client    *pubsub.Client
}

func (p PubSubGcp) Close() error {
	return p.client.Close()
}

func NewPubSubGcp(projectID string) *PubSubGcp {
	return &PubSubGcp{projectID: projectID}
}

func (p PubSubGcp) Publish(ctx context.Context, topic string, message []byte, attributes map[string]string) error {

	t := p.client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data:       message,
		Attributes: attributes,
	})

	_, err := result.Get(ctx)

	return err
}
