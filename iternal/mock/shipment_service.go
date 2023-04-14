package mock

import (
	"time"

	"task/iternal/model"
)

type ShipmentService struct {
}

func (s *ShipmentService) InsertShipment(shipment *model.Shipment) error {
	return nil
}
func (s *ShipmentService) SelectShipmentById(id int) (*model.Shipment, error) {
	shipment := model.NewShipment("TEST", "TEST", "TEST", "TEST")
	shipment.Barcode = "TEST1003TEST1"
	shipment.CreatedAt = time.Time{}
	shipment.Id = 1
	return shipment, nil
}
func (s *ShipmentService) SelectShipmentByBarcode(barcode string) (*model.Shipment, error) {
	shipment := model.NewShipment("TEST", "TEST", "TEST", "TEST")
	shipment.Barcode = "TEST1003TEST1"
	shipment.CreatedAt = time.Time{}
	return shipment, nil
}
func (s *ShipmentService) DeleteShipmentById(id int) error {
	return nil
}
func (s *ShipmentService) CheckBarcodeUnavailable(barcode string) bool {
	return true
}
