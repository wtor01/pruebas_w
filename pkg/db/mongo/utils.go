package mongo

import (
	"bitbucket.org/sercide/data-ingestion/pkg/async"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type total struct {
	Total int `bson:"total"`
}

type ResultPaginate[T any] struct {
	Result []T `bson:"data"`
	Total  int `bson:"total"`
}

type PaginateQuery struct {
	ResultQuery func(...bson.D) bson.A
	CountQuery  func(...bson.D) bson.A
	Paginate    db.Pagination
}

func AppendAggregate(aggregate bson.A) func(ds ...bson.D) bson.A {
	return func(ds ...bson.D) bson.A {
		for _, d := range ds {
			aggregate = append(aggregate, d)
		}
		return aggregate

	}
}

func Paginate[T any](ctx context.Context, collection *mongo.Collection, query PaginateQuery) ([]T, int, error) {

	asyncResult := async.Exec(ctx, buildPaginate[T](collection, query))
	asyncCount := async.Exec(ctx, buildCount(collection, query))

	res := asyncResult.Await(ctx)
	count := asyncCount.Await(ctx)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	if count.Error != nil {
		return nil, 0, count.Error
	}

	return res.Result, count.Result, nil
}

func FindOneAggregate[T any](ctx context.Context, cursor *mongo.Cursor) (T, error) {
	var result T

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return result, err
		}
		break
	}

	if err := cursor.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func buildPaginate[T any](collection *mongo.Collection, query PaginateQuery) func(context.Context) ([]T, error) {

	return func(ctx context.Context) ([]T, error) {
		offset := 0

		if query.Paginate.Offset != nil {
			offset = *query.Paginate.Offset
		}
		var results []T
		aggregate := append(query.ResultQuery(bson.D{{"$skip", offset}}, bson.D{{"$limit", query.Paginate.Limit}}))
		cursor, err := collection.Aggregate(ctx, aggregate)

		if err != nil {
			return results, err
		}

		if err = cursor.All(ctx, &results); err != nil {
			return results, err
		}
		return results, nil
	}
}

func buildCount(collection *mongo.Collection, query PaginateQuery) func(context.Context) (int, error) {

	return func(ctx context.Context) (int, error) {
		aggregate := append(query.CountQuery(bson.D{{"$count", "total"}}))
		cursor, err := collection.Aggregate(ctx, aggregate)

		if err != nil {
			return 0, err
		}

		results, err := FindOneAggregate[total](ctx, cursor)

		if err != nil {
			return 0, err
		}
		return results.Total, nil
	}
}
