package aggregations

import (
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"time"
)

type FeatureValue struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Field string `json:"field"`
	Value string `json:"value"`
}

func NewAggregation(parameters []FeatureValue, config ConfigScheduler, s []ServicePoint) (Aggregation, error) {
	id, err := utils.GenerateId()

	if err != nil {
		return Aggregation{}, err
	}

	return Aggregation{
		Id:             id,
		TypeId:         config.Id,
		TypeName:       config.Name,
		GenerationDate: time.Now().UTC(),
		Parameters:     parameters,
		ServicePoints:  s,
	}, nil
}

type ServicePoint struct {
	Cups         string `bson:"cups"`
	ServicePoint string `bson:"serial_number"`
}

type Aggregation struct {
	Id             string         `bson:"_id"`
	TypeId         string         `bson:"type_id"`
	TypeName       string         `bson:"type_name"`
	GenerationDate time.Time      `bson:"generation_date"`
	Parameters     []FeatureValue `bson:"parameters"`
	ServicePoints  []ServicePoint `bson:"service_points"`
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=AggregationRepository
type AggregationRepository interface {
	GenerateAggregation(ctx context.Context, query ConfigScheduler) ([]Aggregation, error)
	SaveAllAggregations(ctx context.Context, agg []Aggregation) error
}

type GetAggregationsDto struct {
	Offset              *int
	Limit               int
	AggregationConfigId string
	StartDate           time.Time
	EndDate             time.Time
}

type GetAggregationDto struct {
	AggregationConfigId string
}

type CupsStateAggregation struct {
	Type string `bson:"type"`
	CUPS string `bson:"cups"`
}

type AggregationPrevious struct {
	Id                      string                 `bson:"_id"`
	TypeId                  string                 `bson:"type_id"`
	TypeName                string                 `bson:"type_name"`
	GenerationDate          time.Time              `bson:"generation_date"`
	Parameters              []FeatureValue         `bson:"parameters"`
	ServicePoints           []ServicePoint         `bson:"service_points"`
	PreviousAggregation     Aggregation            `bson:"previous_aggregation"`
	CurrentAggregationCups  []CupsStateAggregation `bson:"current_cups_aggregation_type"`
	PreviousAggregationCups []CupsStateAggregation `bson:"previous_cups_aggregation_type"`
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=AggregationMongoRepository
type AggregationMongoRepository interface {
	GetAggregations(ctx context.Context, params GetAggregationsDto) ([]Aggregation, int64, error)
	GetAggregation(ctx context.Context, params GetAggregationDto) (Aggregation, error)
	GetPreviousAggregation(ctx context.Context, params GetAggregationDto) (AggregationPrevious, error)
}
