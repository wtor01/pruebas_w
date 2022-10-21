package festive_days

import (
	"github.com/satori/go.uuid"
	"time"
)

type FestiveDay struct {
	Id           uuid.UUID
	Date         time.Time
	Description  string
	GeographicId string
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedBy    string
	UpdatedAt    time.Time
}
