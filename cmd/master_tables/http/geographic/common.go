package geographic

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
)

// convert dto to response
func GeographicToResponse(d geographic.GeographicZones) GeographicZoneWithId {
	return GeographicZoneWithId{
		Id:          d.ID.String(),
		Code:        d.Code,
		Description: d.Description,
	}
}
