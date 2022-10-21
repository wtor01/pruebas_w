package aggregations

import (
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"time"
)

type Config struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Scheduler   string     `json:"scheduler"`
	SchedulerId string     `json:"-"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Features    []Features `json:"features"`
}

func (c Config) GetName() string {
	return c.Name
}

func (c Config) GetDescription() string {
	return c.Description
}

func (c Config) GetScheduler() string {
	return c.Scheduler
}

func (c Config) GetSchedulerID() string {
	return c.SchedulerId
}

func (c Config) Marshall() ([]byte, error) {
	return c.Event().Marshal()
}

func (c Config) Event() event.Message[Config] {
	return NewSchedulerEvent(c)
}

func (c Config) EventType() string {
	return SchedulerEventType
}

func (c *Config) SetSchedulerId(id string) {
	c.SchedulerId = id
}

func (c Config) Clone() Config {
	features := make([]Features, 0, cap(c.Features))

	return Config{
		Id:          c.Id,
		Name:        c.Name,
		Scheduler:   c.Scheduler,
		SchedulerId: c.SchedulerId,
		StartDate:   c.StartDate,
		EndDate:     c.EndDate,
		Description: c.Description,
		Features:    append(features, c.Features...),
	}
}

func (c *Config) Update(schedulerFormat, description string, startDate, endDate time.Time, features []Features) error {
	if !scheduler.IsValidFormat(schedulerFormat) {
		return scheduler.ErrSchedulerInvalidFormat
	}

	if !isValidDate(endDate) {
		return scheduler.ErrInvalidScheduler
	}

	if !isValidFeatures(features) {
		return scheduler.ErrInvalidScheduler
	}

	c.Scheduler = schedulerFormat
	c.Description = description
	c.StartDate = startDate
	c.EndDate = endDate
	c.Features = features

	return nil
}

func isValidDate(endDate time.Time) bool {
	now := time.Now().UTC()
	return endDate.IsZero() || !endDate.Before(now)
}

func isValidFeatures(features []Features) bool {
	for _, feature := range features {
		if feature.ID == "" || feature.Name == "" || feature.Field == "" {
			return false
		}
	}
	return len(features) > 0
}

func NewConfig(name, description, schedulerFormat string, startDate, endDate time.Time, features []Features) (Config, error) {
	if !scheduler.IsValidName(name) {
		return Config{}, scheduler.ErrSchedulerInvalidName
	}

	if !scheduler.IsValidFormat(schedulerFormat) {
		return Config{}, scheduler.ErrSchedulerInvalidFormat
	}

	if !isValidDate(endDate) {
		return Config{}, scheduler.ErrInvalidScheduler
	}

	if !isValidFeatures(features) {
		return Config{}, scheduler.ErrInvalidScheduler
	}

	id, err := utils.GenerateId()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Id:          id,
		Name:        name,
		Scheduler:   schedulerFormat,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		Features:    features,
	}, nil
}

type ConfigFeatureDto struct {
	Id    string
	Name  string
	Field string
}

func NewConfigFeatureDto(id, name, field string) ConfigFeatureDto {
	return ConfigFeatureDto{
		Id:    id,
		Name:  name,
		Field: field,
	}
}

type GetConfigsQuery struct {
	Query  string
	Limit  int
	Offset *int
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=AggregationConfigRepository
type AggregationConfigRepository interface {
	GetAggregationConfigs(ctx context.Context, dto GetConfigsQuery) ([]Config, int, error)
	GetAggregationConfigById(ctx context.Context, aggregationConfigId string) (Config, error)
	SaveAggregationConfig(ctx context.Context, aggregationConfig Config) (Config, error)
	DeleteAggregationConfig(ctx context.Context, aggregationConfigId string) error
}
