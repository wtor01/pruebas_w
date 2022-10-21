package seasons

import (
	"github.com/satori/go.uuid"
	"time"
)

type Seasons struct {
	ID             uuid.UUID
	Name           string
	Description    string
	GeographicCode string
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedBy      string
	UpdatedAt      time.Time
}

type DayTypes struct {
	ID        string
	SeasonsId string
	Name      string
	Month     int
	IsFestive bool
	CreatedAt time.Time
	CreatedBy string
	UpdatedBy string
	UpdatedAt time.Time
}
