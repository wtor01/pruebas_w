package http

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory/services"
)

type Controller struct {
	srv *services.InventoryServices
}

func NewController(srv *services.InventoryServices) *Controller {
	return &Controller{srv: srv}
}
