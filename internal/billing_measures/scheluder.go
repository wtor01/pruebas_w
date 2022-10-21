package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrSchedulerExist         = errors.New("scheduler exist")
	ErrSchedulerInvalidFormat = errors.New("scheduler invalid format")
	ErrSchedulerInvalidName   = errors.New("scheduler invalid name [a-z|A-Z||\\d|_|-]{1,500}")
	ErrSchedulerIdFormat      = errors.New("invalid scheduler id")
)

type Scheduler struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	SchedulerId   string   `json:"-"`
	DistributorId string   `json:"distributor_id"`
	ServiceType   string   `json:"service_type"`
	PointType     string   `json:"point_type"`
	MeterType     []string `json:"meter_type"`
	ProcessType   string   `json:"process_type"`
	Format        string   `json:"format"`
}

func (s Scheduler) GetDescription() string {
	return "Billing Scheduler"
}

func (s Scheduler) GetName() string {
	return s.Name
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
		ServiceType:   measures.ServiceType(s.ServiceType),
		PointType:     s.PointType,
		MeterType:     s.MeterType,
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

func (s *Scheduler) Update(
	DistributorId string,
	ServiceType string,
	PointType string,
	MeterType []string,
	ProcessType string,
	Format string,
) error {
	if !scheduler.IsValidFormat(Format) {
		return ErrSchedulerInvalidFormat
	}
	s.DistributorId = DistributorId
	s.ServiceType = ServiceType
	s.PointType = PointType
	s.MeterType = MeterType
	s.ProcessType = ProcessType
	s.Format = Format

	return nil
}

func (s Scheduler) Clone() Scheduler {
	meterTypes := make([]string, 0, cap(s.MeterType))

	sc := Scheduler{
		ID:            s.ID,
		DistributorId: s.DistributorId,
		Name:          s.Name,
		SchedulerId:   s.SchedulerId,
		ServiceType:   s.ServiceType,
		PointType:     s.PointType,
		MeterType:     append(meterTypes, s.MeterType...),
		ProcessType:   s.ProcessType,
		Format:        s.Format,
	}

	return sc
}

func isValidDistributorID(distrubutorID string) bool {
	return len(distrubutorID) != 0
}

func NewScheduler(
	ID string,
	Name string,
	SchedulerId string,
	DistributorId string,
	ServiceType string,
	PointType string,
	MeterType []string,
	ProcessType string,
	Format string,
) (Scheduler, error) {

	if !scheduler.IsValidFormat(Format) {
		return Scheduler{}, scheduler.ErrSchedulerInvalidFormat
	}
	if !scheduler.IsValidName(Name) {
		return Scheduler{}, scheduler.ErrSchedulerInvalidName
	}
	if !isValidDistributorID(DistributorId) {
		return Scheduler{}, scheduler.ErrSchedulerIdFormat
	}

	return Scheduler{
		ID:            ID,
		Name:          Name,
		SchedulerId:   SchedulerId,
		DistributorId: DistributorId,
		ServiceType:   ServiceType,
		PointType:     PointType,
		MeterType:     MeterType,
		ProcessType:   ProcessType,
		Format:        Format,
	}, nil
}

type SearchScheduler struct {
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ProcessType   string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=BillingSchedulerRepository
type BillingSchedulerRepository interface {
	GetScheduler(ctx context.Context, id string) (Scheduler, error)
	ListScheduler(ctx context.Context, query db.Pagination) ([]Scheduler, int, error)
	DeleteScheduler(ctx context.Context, id string) error
	SearchScheduler(ctx context.Context, search SearchScheduler) ([]Scheduler, error)
	SaveScheduler(ctx context.Context, s Scheduler) error
}
