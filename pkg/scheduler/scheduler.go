package scheduler

import (
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	schedulergcp "cloud.google.com/go/scheduler/apiv1"
	"context"
	"errors"
	"fmt"
	schedulerpb "google.golang.org/genproto/googleapis/cloud/scheduler/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidScheduler       = errors.New("invalid scheduler")
	ErrSchedulerExist         = errors.New("scheduler exist")
	ErrSchedulerInvalidFormat = errors.New("scheduler invalid format")
	ErrSchedulerInvalidName   = errors.New("scheduler invalid name [a-z|A-Z||\\d|_|-]{1,500}")
	ErrSchedulerIdFormat      = errors.New("invalid scheduler id")
)

type Scheduler interface {
	GetName() string
	GetDescription() string
	GetScheduler() string
	GetSchedulerID() string
	Marshall() ([]byte, error)
	EventType() string
}

//go:generate mockery --case=snake --outpkg=mocks --output=./mocks --name=Client
type Client interface {
	CreateJob(ctx context.Context, s Scheduler, topic string) (string, error)
	UpdateJob(ctx context.Context, s Scheduler, topic string) error
	Close() error
	DeleteJob(ctx context.Context, id string) error
}

type ClientCreator = func(ctx context.Context) (Client, error)

type GcpScheduler struct {
	ProjectID string
	Location  string
	TimeZone  string
	client    *schedulergcp.CloudSchedulerClient
}

func NewGcpScheduler(projectID string, location string, timeZone string) ClientCreator {
	return func(ctx context.Context) (Client, error) {
		clientScheduler, err := schedulergcp.NewCloudSchedulerClient(ctx)
		return &GcpScheduler{ProjectID: projectID, Location: location, TimeZone: timeZone, client: clientScheduler}, err
	}
}

func (g GcpScheduler) Close() error {
	return g.client.Close()
}
func (g GcpScheduler) DeleteJob(ctx context.Context, id string) error {
	return g.client.DeleteJob(ctx, &schedulerpb.DeleteJobRequest{
		Name: id,
	})
}

func (g GcpScheduler) CreateJob(ctx context.Context, entity Scheduler, topic string) (string, error) {

	d, _ := entity.Marshall()

	job, err := g.client.CreateJob(ctx, &schedulerpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", g.ProjectID, g.Location),
		Job: &schedulerpb.Job{
			Name:        fmt.Sprintf("projects/%s/locations/%s/jobs/%s", g.ProjectID, g.Location, entity.GetName()),
			Description: entity.GetDescription(),
			Target: &schedulerpb.Job_PubsubTarget{
				PubsubTarget: &schedulerpb.PubsubTarget{
					TopicName: fmt.Sprintf("projects/%s/topics/%s", g.ProjectID, topic),
					Data:      d,
					Attributes: map[string]string{
						event.EventTypeKey: entity.EventType(),
					},
				},
			},
			Schedule: entity.GetScheduler(),
			TimeZone: g.TimeZone,
			State:    schedulerpb.Job_ENABLED,
		},
	})

	if err != nil {
		return "", err
	}

	return job.Name, nil
}

func (g GcpScheduler) UpdateJob(ctx context.Context, entity Scheduler, topic string) error {

	d, _ := entity.Marshall()

	_, err := g.client.UpdateJob(ctx, &schedulerpb.UpdateJobRequest{
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"schedule",
				"time_zone",
				"description",
				"pubsub_target.topic_name",
				"pubsub_target.data",
				"pubsub_target.attributes",
			},
		},
		Job: &schedulerpb.Job{
			Name:        entity.GetSchedulerID(),
			Description: entity.GetDescription(),
			Target: &schedulerpb.Job_PubsubTarget{
				PubsubTarget: &schedulerpb.PubsubTarget{
					TopicName: fmt.Sprintf("projects/%s/topics/%s", g.ProjectID, topic),
					Data:      d,
					Attributes: map[string]string{
						event.EventTypeKey: entity.EventType(),
					},
				},
			},
			Schedule: entity.GetScheduler(),
			TimeZone: g.TimeZone,
			State:    schedulerpb.Job_ENABLED,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func IsValidName(name string) bool {
	if len(name) == 0 {
		return false
	}
	re := regexp.MustCompile(`[a-z|A-Z|\d|_|-]{1,500}`)
	str := re.FindString(name)

	return str == name
}

func IsValidFormat(str string) bool {
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
