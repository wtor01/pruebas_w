package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CreateFeatures struct {
	repository aggregations.AggregationsFeaturesRepository
	tracer     trace.Tracer
}

func NewCreateFeaturesService(repository aggregations.AggregationsFeaturesRepository) *CreateFeatures {
	return &CreateFeatures{
		repository: repository,
		tracer:     telemetry.GetTracer(),
	}
}

type CreateFeaturesDTO struct {
	ID    string
	Name  string
	Field string
}

func (f CreateFeatures) Handler(ctx context.Context, dto CreateFeaturesDTO) (aggregations.Features, error) {
	ctx, span := f.tracer.Start(ctx, "SaveFeatures - Handler")
	defer span.End()

	newFeatures, err := aggregations.NewFeatures(dto.ID, dto.Name, dto.Field)

	if err != nil {
		return aggregations.Features{}, err
	}

	featuresFound, err := f.repository.SearchFeatures(ctx, aggregations.SearchFeatures{
		Name:  newFeatures.Name,
		Field: newFeatures.Field,
	})

	if err != nil {
		return aggregations.Features{}, err
	}
	if len(featuresFound) != 0 {
		return aggregations.Features{}, errors.New("error; Features Already exists")
	}

	err = f.repository.SaveFeatures(ctx, newFeatures)

	if err != nil {
		return aggregations.Features{}, err
	}
	span.SetAttributes(attribute.Bool("response", true))

	return newFeatures, err
}
