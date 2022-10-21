package clients

import (
	"context"
)

type ValidationConfigExtraConfig struct {
	Id     string            `json:"id"`
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
}

type ValidationConfigParams struct {
	Type        string             `json:"type"`
	Validations []ValidationParams `json:"params"`
}

type ValidationParams struct {
	Id       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Keys     []string          `json:"keys"`
	Required bool              `json:"required"`
	Config   map[string]string `json:"config"`
}

type ValidationConfig struct {
	Id                string `json:"id"`
	ValidationMeasure struct {
		Id          string                 `json:"id"`
		Name        string                 `json:"name"`
		Action      string                 `json:"action"`
		Enabled     bool                   `json:"enabled"`
		MeasureType string                 `json:"measure_type"`
		Type        string                 `json:"type"`
		Code        string                 `json:"code"`
		Message     string                 `json:"message"`
		Description string                 `json:"description"`
		Params      ValidationConfigParams `json:"validation_measure"`
	}
	Action      string                        `json:"action"`
	Enabled     bool                          `json:"enabled"`
	ExtraConfig []ValidationConfigExtraConfig `json:"extra_config"`
}

//go:generate mockery --case=snake --outpkg=mocks --output=./clients_mocks --name=Validation
type Validation interface {
	GetValidationConfigList(ctx context.Context, distributorID string, valType string) ([]ValidationConfig, error)
}
