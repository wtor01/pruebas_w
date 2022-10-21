package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func New(ctx context.Context, cnf config.Config) (*mongo.Client, error) {
	uri := fmt.Sprintf(cnf.MeasureDBHost, cnf.MeasureDBUser, cnf.MeasureDBPass)
	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {

		return client, err
	}
	err = client.Ping(ctx, readpref.Primary())

	return client, err
}
