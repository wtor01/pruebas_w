package aggregations

import (
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"context"
	"errors"
)

type Features struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Field string `json:"field"`
}

func (f *Features) Update(name, field string) error {
	if !isValidString(name) || !isValidString(field) {
		return errors.New("invalid input")
	}
	f.Name = name
	f.Field = field
	return nil
}

func NewFeatures(id, name, field string) (Features, error) {
	if !isValidString(name) || !isValidString(field) {
		return Features{}, errors.New("invalid input")
	}
	return Features{
		ID:    id,
		Name:  name,
		Field: field,
	}, nil
}
func isValidString(str string) bool {
	return !(str == "" || str == " ")
}

type SearchFeatures struct {
	Name  string
	Field string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=AggregationsFeaturesRepository
type AggregationsFeaturesRepository interface {
	GetFeatures(ctx context.Context, id string) (Features, error)
	ListFeatures(ctx context.Context, query db.Pagination) ([]Features, int, error)
	DeleteFeatures(ctx context.Context, id string) error
	SaveFeatures(ctx context.Context, obj Features) error
	SearchFeatures(ctx context.Context, features SearchFeatures) ([]Features, error)
	GetFeaturesByIds(ctx context.Context, ids []string) ([]Features, error)
}
