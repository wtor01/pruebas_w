package geographic

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"github.com/satori/go.uuid"
	"time"
)

//GeographicZones service struct
type GeographicZones struct {
	ID          uuid.UUID
	Code        string
	Description string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedBy   string
	UpdatedAt   time.Time
}

//Search search used in Get all
type Search struct {
	Q             string
	Limit         int
	Offset        *int
	Sort          map[string]string
	CurrentUser   auth.User
	DistributorID string
}
