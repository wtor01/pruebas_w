package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UpdateFeatures struct {
	repository aggregations.AggregationsFeaturesRepository
	tracer     trace.Tracer
}

func NewUpdateFeaturesService(repository aggregations.AggregationsFeaturesRepository) *UpdateFeatures {
	return &UpdateFeatures{
		repository: repository,
		tracer:     telemetry.GetTracer(),
	}
}

type UpdateFeaturesDTO struct {
	ID    string
	Field string
	Name  string
}

func (f UpdateFeatures) Handler(ctx context.Context, dto UpdateFeaturesDTO) (aggregations.Features, error) {
	ctx, span := f.tracer.Start(ctx, "Update - Handler")
	defer span.End()

	features, err := f.repository.GetFeatures(ctx, dto.ID)
	if err != nil {
		return aggregations.Features{}, err
	}

	err = features.Update(
		dto.Name,
		dto.Field,
	)
	if err != nil {
		return aggregations.Features{}, err
	}

	savedFeatures, err := f.repository.SearchFeatures(ctx, aggregations.SearchFeatures{
		Name:  features.Name,
		Field: features.Field,
	})

	if err != nil {
		return aggregations.Features{}, err
	}
	if (len(savedFeatures) == 1 && savedFeatures[0].ID != features.ID) || len(savedFeatures) > 1 {
		return aggregations.Features{}, errors.New("features already exists")
	}

	err = f.repository.SaveFeatures(ctx, features)
	span.SetAttributes(attribute.Bool("response", true))

	return features, err

}
