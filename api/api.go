package api

import (
	"task/iternal/shipments"
)

type shipmentApi struct {
	service *shipments.ShipmentService
}

func NewShipmentApi(service shipments.ShipmentService) *shipmentApi {
	return &shipmentApi{&service}
}
