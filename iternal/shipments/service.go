package shipments

import (
	"database/sql"
	"fmt"
	"log"
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

func NewSQLShipmentService(db *sql.DB, filter *bloom.Filter) *sqlShipmentService {
	return &sqlShipmentService{
		DB:     db,
		Filter: filter,
	}
}

func (service *sqlShipmentService) Save(shipment *Shipment) error {
	availability := service.CheckBarcodeAvailability(shipment.Barcode)

	if availability == false {
		service.Filter.AddToFilter(shipment.Barcode)
		_, err := service.DB.Exec(
			"INSERT INTO shipments (barcode, sender, receiver, is_delivered, origin, destination, created_at) "+
				"VALUES ($1, $2, $3, $4, $5, $6, $7)",
			shipment.Barcode, shipment.Sender, shipment.Receiver, shipment.IsDelivered, shipment.Origin, shipment.Destination, shipment.CreatedAt)

		if err != nil {
			log.Printf("failed to save shipment: %v", err)
			return fmt.Errorf("failed to save shipment: %v", err)
		}

		log.Printf("shipment saved: %v", shipment)
		return nil
	}

	log.Printf("barcode %s already exists in the system", shipment.Barcode)
	return fmt.Errorf("barcode %s already exists in the system", shipment.Barcode)
}

func (service *sqlShipmentService) GetById(id int) (*Shipment, error) {
	var s Shipment

	row := service.DB.QueryRow(
		"SELECT s.id, s.barcode, s.sender, s.receiver, s.is_delivered, s.origin, s.destination, s.created_at "+
			"FROM shipments s "+
			"WHERE s.id = $1",
		id)
	err := row.Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.CreatedAt)
	if err != nil {
		log.Printf("failed to scan shipment: %v", err)
		return nil, fmt.Errorf("failed to scan shipment: %v", err)
	}

	log.Printf("got shipment by id: %v", s)
	return &s, nil
}

func (service *sqlShipmentService) GetByBarcode(barcode string) (*Shipment, error) {
	var s Shipment
	row := service.DB.QueryRow(
		"SELECT s.id, s.barcode, s.sender, s.receiver, s.is_delivered, s.origin, s.destination, s.created_at "+
			"FROM shipments s "+
			"WHERE s.barcode = $1",
		barcode)

	err := row.Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.CreatedAt)
	if err != nil {
		log.Printf("failed to get shipment: %v", err)
		return nil, fmt.Errorf("failed to get shipment: %v", err)
	}

	log.Printf("got shipment by barcode: %v", s)
	return &s, nil
}

func (service *sqlShipmentService) DeleteById(id int) error {
	_, err := service.DB.Exec("DELETE FROM shipments WHERE id = $1", id)
	if err != nil {
		log.Printf("failed to delete shipment with id: %d, error: %v", id, err)
		return fmt.Errorf("failed to delete shipment with id: %d, error: %v", id, err)
	}

	log.Printf("shipment with id: %d has been deleted", id)
	return nil
}

func (service *sqlShipmentService) CheckBarcodeAvailability(barcode string) bool {
	log.Printf("check barcode: %s availability", barcode)
	return service.Filter.Check(barcode)
}
