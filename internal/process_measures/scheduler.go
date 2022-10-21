package process_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrSchedulerExist              = errors.New("scheduler exist")
	ErrSchedulerInvalidFormat      = errors.New("scheduler invalid format")
	ErrSchedulerInvalidName        = errors.New("scheduler invalid name [a-z|A-Z||\\d|_|-]{1,500}")
	ErrSchedulerInvalidDescription = errors.New("scheduler invalid description {0,500}")
	ErrSchedulerIdFormat           = errors.New("invalid scheduler id")
)

type Scheduler struct {
	ID            string               `json:"id"`
	DistributorId string               `json:"distributor_id"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	SchedulerId   string               `json:"-"`
	ServiceType   string               `json:"service_type"`
	PointType     string               `json:"point_type"`
	MeterType     []string             `json:"meter_type"`
	ReadingType   measures.ReadingType `json:"reading_type"`
	Format        string               `json:"format"`
}

func (s Scheduler) GetName() string {
	return s.Name
}

func (s Scheduler) GetDescription() string {
	return s.Description
}

func (s Scheduler) GetScheduler() string {
	return s.Format
}

func (s Scheduler) Marshall() ([]byte, error) {
	return json.Marshal(s.Event())
}

func (s Scheduler) Event() measures.SchedulerEvent {
	return measures.NewSchedulerEvent(SchedulerEventType, measures.SchedulerEventPayload{
		ID:            s.ID,
		DistributorId: s.DistributorId,
		Name:          s.Name,
		Description:   s.Description,
		ServiceType:   measures.ServiceType(s.ServiceType),
		PointType:     s.PointType,
		MeterType:     s.MeterType,
		ReadingType:   s.ReadingType,
		Format:        s.Format,
	})
}

func (s Scheduler) EventType() string {
	return SchedulerEventType
}

func (s *Scheduler) SetSchedulerId(id string) {
	s.SchedulerId = id
}

func (s *Scheduler) GetSchedulerID() string {
	return s.SchedulerId
}

func isValidSchedulerFormat(str string) bool {
	steps := strings.Split(str, " ")
	if len(steps) != 5 {
		return false
	}

	for i, s := range steps {
		if s == "*" {
			continue
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return false
		}

		if i == 0 && (v > 59 || v < 0) {
			return false
		}
		if i == 1 && (v > 23 || v < 0) {
			return false
		}
		if i == 2 && (v > 31 || v < 1) {
			return false
		}
		if i == 3 && (v > 12 || v < 1) {
			return false
		}
		if i == 4 && (v > 6 || v < 0) {
			return false
		}

	}

	return true
}

func isValidSchedulerName(name string) bool {
	if len(name) == 0 {
		return false
	}
	re := regexp.MustCompile(`[a-z|A-Z|\d|_|-]{1,500}`)
	str := re.FindString(name)

	return str == name
}

func isValidSchedulerDescription(description string) bool {
	return len(description) <= 500
}

func NewScheduler(
	ID string,
	Name string,
	Description string,
	DistributorId string,
	SchedulerId string,
	ServiceType string,
	PointType string,
	MeterType []string,
	ReadingType string,
	Format string,
) (Scheduler, error) {

	if !isValidSchedulerFormat(Format) {
		return Scheduler{}, ErrSchedulerInvalidFormat
	}
	if !isValidSchedulerName(Name) {
		return Scheduler{}, ErrSchedulerInvalidName
	}
	if !isValidSchedulerDescription(Description) {
		return Scheduler{}, ErrSchedulerInvalidDescription
	}

	return Scheduler{
		ID:            ID,
		Name:          Name,
		Description:   Description,
		DistributorId: DistributorId,
		SchedulerId:   SchedulerId,
		ServiceType:   ServiceType,
		PointType:     PointType,
		MeterType:     MeterType,
		ReadingType:   measures.ReadingType(ReadingType),
		Format:        Format,
	}, nil
}

func (s *Scheduler) Update(
	Description string,
	DistributorId string,
	ServiceType string,
	PointType string,
	MeterType []string,
	ReadingType string,
	Format string,
) error {
	if !isValidSchedulerFormat(Format) {
		return ErrSchedulerInvalidFormat
	}
	if !isValidSchedulerDescription(Description) {
		return ErrSchedulerInvalidDescription
	}
	s.Description = Description
	s.DistributorId = DistributorId
	s.ServiceType = ServiceType
	s.PointType = PointType
	s.MeterType = MeterType
	s.ReadingType = measures.ReadingType(ReadingType)
	s.Format = Format

	return nil
}

func (s Scheduler) Clone() Scheduler {
	meterTypes := make([]string, 0, cap(s.MeterType))

	sc := Scheduler{
		ID:            s.ID,
		DistributorId: s.DistributorId,
		Name:          s.Name,
		Description:   s.Description,
		SchedulerId:   s.SchedulerId,
		ServiceType:   s.ServiceType,
		PointType:     s.PointType,
		MeterType:     append(meterTypes, s.MeterType...),
		ReadingType:   s.ReadingType,
		Format:        s.Format,
	}

	return sc
}

type SearchScheduler struct {
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ReadingType   string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=SchedulerRepository
type SchedulerRepository interface {
	SaveScheduler(ctx context.Context, s Scheduler) error
	SearchScheduler(ctx context.Context, search SearchScheduler) ([]Scheduler, error)
	DeleteScheduler(ctx context.Context, id string) error
	GetScheduler(ctx context.Context, id string) (Scheduler, error)
	ListScheduler(ctx context.Context, query db.Pagination) ([]Scheduler, int, error)
}
