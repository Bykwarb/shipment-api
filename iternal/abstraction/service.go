package abstraction

import (
	"task/iternal/model"
)

type ShipmentService interface {
	InsertShipment(shipment *model.Shipment) error
	SelectShipmentById(id int) (*model.Shipment, error)
	SelectShipmentByBarcode(barcode string) (*model.Shipment, error)
	DeleteShipmentById(id int) error
	CheckBarcodeUnavailable(barcode string) bool
}
