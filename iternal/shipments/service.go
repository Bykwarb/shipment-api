package shipments

import (
	"database/sql"
	"errors"
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
	db     *sql.DB
	filter *bloom.Filter
}

func NewSQLShipmentService(db *sql.DB, filter *bloom.Filter) *sqlShipmentService {
	return &sqlShipmentService{
		db:     db,
		filter: filter,
	}
}

func (service *sqlShipmentService) Save(shipment *Shipment) error {
	if service.filter.Check(shipment.Barcode) {
		return errors.New("barcode already exists")
	}

	if _, err := service.db.Exec("INSERT INTO shipments (barcode, sender, receiver, is_delivered, origin, destination, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		shipment.Barcode, shipment.Sender, shipment.Receiver, shipment.IsDelivered, shipment.Origin, shipment.Destination, shipment.CreatedAt); err != nil {
		log.Printf("failed to save shipment: %v", err)
		return fmt.Errorf("failed to save shipment: %v", err)
	}
	service.filter.AddToFilter(shipment.Barcode)
	log.Printf("shipment saved: %v", shipment)
	return nil
}

func (service *sqlShipmentService) GetById(id int) (*Shipment, error) {
	var shipment Shipment
	if err := service.db.QueryRow("SELECT id, barcode, sender, receiver, is_delivered, origin, destination, created_at FROM shipments WHERE id = $1", id).
		Scan(&shipment.Id, &shipment.Barcode, &shipment.Sender, &shipment.Receiver, &shipment.IsDelivered, &shipment.Origin, &shipment.Destination, &shipment.CreatedAt); err != nil {
		log.Printf("failed to get shipment by ID: %v", err)
		return nil, fmt.Errorf("failed to get shipment by ID: %v", err)
	}

	log.Printf("got shipment by ID: %v", shipment)
	return &shipment, nil
}

func (service *sqlShipmentService) GetByBarcode(barcode string) (*Shipment, error) {
	var s Shipment
	if err := service.db.QueryRow("SELECT s.id, s.barcode, s.sender, s.receiver, s.is_delivered, s.origin, s.destination, s.created_at FROM shipments s WHERE s.barcode = $1", barcode).
		Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.CreatedAt); err != nil {
		log.Printf("failed to get shipment: %v", err)
		return nil, fmt.Errorf("failed to get shipment: %v", err)
	}
	log.Printf("got shipment by barcode: %v", s)
	return &s, nil
}

func (service *sqlShipmentService) DeleteById(id int) error {
	if _, err := service.db.Exec("DELETE FROM shipments WHERE id = $1", id); err != nil {
		log.Printf("failed to delete shipment with id: %d, error: %v", id, err)
		return fmt.Errorf("failed to delete shipment with id: %d, error: %v", id, err)
	}

	log.Printf("shipment with id: %d has been deleted", id)
	return nil
}

func (service *sqlShipmentService) CheckBarcodeAvailability(barcode string) bool {
	log.Printf("check barcode: %s availability", barcode)
	return service.filter.Check(barcode)
}
