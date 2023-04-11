package shipments

import (
	"database/sql"
	"errors"
	"task/iternal/bloom"
)

type ShipmentService interface {
	Save(shipment *Shipment) error
	GetById(id int) (*Shipment, error)
	GetByBarcode(barcode string) (*Shipment, error)
	DeleteById(id int) error
	CheckBarcodeAvailability(barcode string) bool
}

type sqlShipmentService struct {
	DB     *sql.DB
	Filter *bloom.Filter
}

func NewShipmentService(db *sql.DB, filter *bloom.Filter) *sqlShipmentService {
	return &sqlShipmentService{DB: db, Filter: filter}
}

func (service *sqlShipmentService) Save(shipment *Shipment) error {
	availability := service.Filter.Check(shipment.Barcode)
	if availability {
		service.Filter.AddToFilter(shipment.Barcode)
		_, err := service.DB.Exec(
			"INSERT INTO shipments (barcode, sender, receiver, is_delivered, origin, destination, created_at) "+
				"VALUES ($1, $2, $3, $4, $5, $6, $7)",
			shipment.Barcode, shipment.Sender, shipment.Receiver, shipment.IsDelivered, shipment.Origin, shipment.Destination, shipment.DepartureDate)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("barcodes not available")
}

func (service *sqlShipmentService) GetById(id int) (*Shipment, error) {
	var s Shipment

	row := service.DB.QueryRow(
		"SELECT s.id, s.barcode, s.sender, s.receiver, s.is_delivered, s.origin, s.destination, s.created_at "+
			"FROM shipments s "+
			"WHERE s.id = $1", id)
	err := row.Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.DepartureDate)
	return &s, err
}

func (service *sqlShipmentService) GetByBarcode(barcode string) (*Shipment, error) {
	var s Shipment
	row := service.DB.QueryRow(
		"SELECT s.id, s.barcode, s.sender, s.receiver, s.is_delivered, s.origin, s.destination, s.created_at "+
			"FROM shipments s "+
			"WHERE s.barcode = $1", barcode)

	err := row.Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.DepartureDate)
	return &s, err
}

func (service *sqlShipmentService) DeleteById(id int) error {
	_, err := service.DB.Exec("DELETE FROM shipments WHERE shipments.id = $1", id)
	return err
}

func (service *sqlShipmentService) CheckBarcodeAvailability(barcode string) bool {
	return service.Filter.Check(barcode)
}
