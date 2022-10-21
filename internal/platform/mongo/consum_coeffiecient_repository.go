package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"context"
)

type ConsumCoefficientRepositoryMongo struct {
	collection string
	database   string
}

func (repository ConsumCoefficientRepositoryMongo) Search(ctx context.Context, q billing_measures.QueryConsumCoefficient) (float64, error) {
	return 0.33, nil
}
