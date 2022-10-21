package inventory

import (
	"context"
	"github.com/satori/go.uuid"
	"time"
)

const (
	TelegestionType string = "telegestion"
	TelemedidaType  string = "telemedida"
	OtherType       string = "otros"
)

type MeasureEquipment struct {
	ID                uuid.UUID
	SerialNumber      string
	Technology        string
	Type              string
	Brand             string
	Model             string
	ActiveConstant    float64
	ReactiveConstant  float64
	MaximeterConstant float64

	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	Distributor Distributor
}

type MeasureEquipmentRepository interface {
	AllMeasureEquipments(ctx context.Context, ids ...string) ([]MeasureEquipment, error)
	SaveMeasureEquipment(ctx context.Context, dt MeasureEquipment) (MeasureEquipment, error)
}
