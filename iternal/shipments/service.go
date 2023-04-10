package shipments

import (
	"database/sql"
	"fmt"
)

type ShipmentService interface {
	SaveShipment(shipment *Shipment)
	GetShipmentById(id int) *Shipment
	DeleteShipmentById(id int)
}

type SqlShipmentService struct {
	DB *sql.DB
}

func (service *SqlShipmentService) SaveShipment(shipment *Shipment) {

}

func (service *SqlShipmentService) GetShipmentById(id int) (*Shipment, error) {
	var s Shipment
	query := fmt.Sprintf("SELECT s.id, b.barcode, s.sender, s.receiver, s.is_delivered, s.origin, s.destination, s.created_at "+
		"FROM shipments s "+
		"JOIN barcodes b ON b.id = s.barcode_id "+
		"WHERE s.id = %d;", id)
	row := service.DB.QueryRow(query)
	err := row.Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.DepartureDate)
	return &s, err
}

func (service *SqlShipmentService) DeleteShipmentById() {

}
