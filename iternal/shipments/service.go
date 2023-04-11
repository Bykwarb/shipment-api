package shipments

import (
	"database/sql"
)

type ShipmentService interface {
	SaveShipment(shipment *shipment) error
	GetShipmentById(id int) (*shipment, error)
	GetShipmentByBarcode(barcode string) (*shipment, error)
	DeleteShipmentById(id int) error
}

type sqlShipmentService struct {
	DB *sql.DB
}

func NewShipmentService(db *sql.DB) *sqlShipmentService {
	return &sqlShipmentService{DB: db}
}

func (service *sqlShipmentService) SaveShipment(shipment *shipment) error {

	_, err := service.DB.Exec(
		"INSERT INTO shipments (barcode, sender, receiver, is_delivered, origin, destination, created_at) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7)",
		shipment.Barcode, shipment.Sender, shipment.Receiver, shipment.IsDelivered, shipment.Origin, shipment.Destination, shipment.DepartureDate)
	if err != nil {
		return err
	}

	return nil
}

func (service *sqlShipmentService) GetShipmentById(id int) (*shipment, error) {
	var s shipment

	row := service.DB.QueryRow(
		"SELECT s.id, s.barcode, s.sender, s.receiver, s.is_delivered, s.origin, s.destination, s.created_at "+
			"FROM shipments s "+
			"WHERE s.id = $1", id)
	err := row.Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.DepartureDate)
	return &s, err
}

func (service *sqlShipmentService) GetShipmentByBarcode(barcode string) (*shipment, error) {
	var s shipment

	row := service.DB.QueryRow(
		"SELECT s.id, s.barcode, s.sender, s.receiver, s.is_delivered, s.origin, s.destination, s.created_at "+
			"FROM shipments s "+
			"WHERE s.barcode = $1", barcode)

	err := row.Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.DepartureDate)
	return &s, err
}

func (service *sqlShipmentService) DeleteShipmentById(id int) error {
	_, err := service.DB.Exec("DELETE FROM shipments WHERE shipments.id = $1", id)
	return err
}
