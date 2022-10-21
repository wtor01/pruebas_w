// Package scheduler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package scheduler

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for BillingMeasuresSchedulerUpdatableMeterType.
const (
	BillingMeasuresSchedulerUpdatableMeterTypeOTHER BillingMeasuresSchedulerUpdatableMeterType = "OTHER"

	BillingMeasuresSchedulerUpdatableMeterTypeTLG BillingMeasuresSchedulerUpdatableMeterType = "TLG"

	BillingMeasuresSchedulerUpdatableMeterTypeTLM BillingMeasuresSchedulerUpdatableMeterType = "TLM"
)

// Defines values for BillingMeasuresSchedulerUpdatablePointType.
const (
	BillingMeasuresSchedulerUpdatablePointTypeN1 BillingMeasuresSchedulerUpdatablePointType = "1"

	BillingMeasuresSchedulerUpdatablePointTypeN2 BillingMeasuresSchedulerUpdatablePointType = "2"

	BillingMeasuresSchedulerUpdatablePointTypeN3 BillingMeasuresSchedulerUpdatablePointType = "3"

	BillingMeasuresSchedulerUpdatablePointTypeN4 BillingMeasuresSchedulerUpdatablePointType = "4"

	BillingMeasuresSchedulerUpdatablePointTypeN5 BillingMeasuresSchedulerUpdatablePointType = "5"
)

// Defines values for BillingMeasuresSchedulerUpdatableProcessType.
const (
	BillingMeasuresSchedulerUpdatableProcessTypeDCNOTLG BillingMeasuresSchedulerUpdatableProcessType = "D-C NO TLG"

	BillingMeasuresSchedulerUpdatableProcessTypeDCTLG BillingMeasuresSchedulerUpdatableProcessType = "D-C TLG"

	BillingMeasuresSchedulerUpdatableProcessTypeGD BillingMeasuresSchedulerUpdatableProcessType = "G-D"
)

// Defines values for BillingMeasuresSchedulerUpdatableServiceType.
const (
	BillingMeasuresSchedulerUpdatableServiceTypeDC BillingMeasuresSchedulerUpdatableServiceType = "D-C"

	BillingMeasuresSchedulerUpdatableServiceTypeDD BillingMeasuresSchedulerUpdatableServiceType = "D-D"

	BillingMeasuresSchedulerUpdatableServiceTypeGD BillingMeasuresSchedulerUpdatableServiceType = "G-D"
)

// BillingMeasuresScheduler defines model for BillingMeasuresScheduler.
type BillingMeasuresScheduler struct {
	// Embedded fields due to inline allOf schema
	Id string `json:"id"`
	// Embedded struct due to allOf(#/components/schemas/BillingMeasuresSchedulerBase)
	BillingMeasuresSchedulerBase `yaml:",inline"`
}

// BillingMeasuresSchedulerBase defines model for BillingMeasuresSchedulerBase.
type BillingMeasuresSchedulerBase struct {
	// Embedded fields due to inline allOf schema
	Name string `json:"name"`
	// Embedded struct due to allOf(#/components/schemas/BillingMeasuresSchedulerUpdatable)
	BillingMeasuresSchedulerUpdatable `yaml:",inline"`
}

// BillingMeasuresSchedulerUpdatable defines model for BillingMeasuresSchedulerUpdatable.
type BillingMeasuresSchedulerUpdatable struct {
	DistributorId *string                                      `json:"distributor_id,omitempty"`
	MeterType     []BillingMeasuresSchedulerUpdatableMeterType `binding:"dive,oneof='TLG' 'TLM' 'OTHER'" json:"meter_type"`
	PointType     BillingMeasuresSchedulerUpdatablePointType   `json:"point_type"`
	ProcessType   BillingMeasuresSchedulerUpdatableProcessType `binding:"oneof='D-C TLG' 'D-C NO TLG' 'G-D'" json:"process_type"`
	Scheduler     string                                       `json:"scheduler"`
	ServiceType   BillingMeasuresSchedulerUpdatableServiceType `json:"service_type"`
}

// BillingMeasuresSchedulerUpdatableMeterType defines model for BillingMeasuresSchedulerUpdatable.MeterType.
type BillingMeasuresSchedulerUpdatableMeterType string

// BillingMeasuresSchedulerUpdatablePointType defines model for BillingMeasuresSchedulerUpdatable.PointType.
type BillingMeasuresSchedulerUpdatablePointType string

// BillingMeasuresSchedulerUpdatableProcessType defines model for BillingMeasuresSchedulerUpdatable.ProcessType.
type BillingMeasuresSchedulerUpdatableProcessType string

// BillingMeasuresSchedulerUpdatableServiceType defines model for BillingMeasuresSchedulerUpdatable.ServiceType.
type BillingMeasuresSchedulerUpdatableServiceType string

// Pagination defines model for Pagination.
type Pagination struct {
	Links struct {
		// url for request next list
		Next *string `json:"next,omitempty"`

		// url for request previous list
		Prev *string `json:"prev,omitempty"`

		// url for request current list
		Self string `json:"self"`
	} `json:"_links"`
	Count  int  `json:"count"`
	Limit  int  `json:"limit"`
	Offset *int `json:"offset,omitempty"`
	Size   int  `json:"size"`
}

// ListBillingMeasuresScheduler defines model for ListBillingMeasuresScheduler.
type ListBillingMeasuresScheduler struct {
	// Embedded struct due to allOf(#/components/schemas/Pagination)
	Pagination `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Results []BillingMeasuresScheduler `json:"results"`
}

// CreateBillingMeasuresScheduler defines model for CreateBillingMeasuresScheduler.
type CreateBillingMeasuresScheduler BillingMeasuresSchedulerBase

// PatchBillingMeasuresScheduler defines model for PatchBillingMeasuresScheduler.
type PatchBillingMeasuresScheduler BillingMeasuresSchedulerUpdatable

// ListBillingMeasuresSchedulerParams defines parameters for ListBillingMeasuresScheduler.
type ListBillingMeasuresSchedulerParams struct {
	// The number of items to skip before starting to collect the result set
	Offset *int `json:"offset,omitempty"`

	// The numbers of items to return
	Limit int `json:"limit"`
}

// CreateBillingMeasuresSchedulerJSONRequestBody defines body for CreateBillingMeasuresScheduler for application/json ContentType.
type CreateBillingMeasuresSchedulerJSONRequestBody CreateBillingMeasuresScheduler

// PatchBillingMeasuresSchedulerJSONRequestBody defines body for PatchBillingMeasuresScheduler for application/json ContentType.
type PatchBillingMeasuresSchedulerJSONRequestBody PatchBillingMeasuresScheduler