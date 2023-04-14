package shipment

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"task/iternal/model"
	"task/pkg/bloom"
)

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

func (sss *sqlShipmentService) InsertShipment(shipment *model.Shipment) error {
	if sss.filter.Check(shipment.Barcode) {
		return errors.New("barcode already exists")
	}

	if len(shipment.Barcode) > 25 || len(shipment.Barcode) < 13 {
		log.Printf("barcode length must be <= then 25 and >= 13, barcode: %s", shipment.Barcode)
		return fmt.Errorf("barcode length must be <= then 25 and >= 13")
	}

	if _, err := sss.db.Exec(
		"INSERT INTO shipment (barcode, sender, receiver, is_delivered, origin, destination, created_at) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7)",
		shipment.Barcode,
		shipment.Sender,
		shipment.Receiver,
		shipment.IsDelivered,
		shipment.Origin,
		shipment.Destination,
		shipment.CreatedAt); err != nil {
		log.Println("failed to save shipment")
		return errors.New("failed to save shipment")
	}

	sss.filter.AddToFilter(shipment.Barcode)

	log.Printf("shipment saved: %v", shipment)

	return nil
}

func (sss *sqlShipmentService) SelectShipmentById(id int) (*model.Shipment, error) {
	var shipment model.Shipment

	if err := sss.db.QueryRow("SELECT id, barcode, sender, receiver, is_delivered, origin, destination, created_at FROM shipments WHERE id = $1", id).
		Scan(&shipment.Id,
			&shipment.Barcode,
			&shipment.Sender,
			&shipment.Receiver,
			&shipment.IsDelivered,
			&shipment.Origin,
			&shipment.Destination,
			&shipment.CreatedAt); err != nil {

		log.Printf("failed to get shipment with Id: %d", id)
		return nil, fmt.Errorf("failed to get shipment with Id: %d", id)
	}

	log.Printf("got shipment by Id: %v", shipment)

	return &shipment, nil
}

func (sss *sqlShipmentService) SelectShipmentByBarcode(barcode string) (*model.Shipment, error) {
	var s model.Shipment
	if err := sss.db.QueryRow("SELECT id, barcode, sender, receiver, is_delivered, origin, destination, created_at FROM shipments WHERE barcode = $1", barcode).
		Scan(&s.Id, &s.Barcode, &s.Sender, &s.Receiver, &s.IsDelivered, &s.Origin, &s.Destination, &s.CreatedAt); err != nil {
		log.Printf("failed to get shipment with barcode: %s", barcode)
		return nil, fmt.Errorf("failed to get shipment with barcode: %s", barcode)
	}

	log.Printf("got shipment by barcode: %v", s)

	return &s, nil
}

func (sss *sqlShipmentService) DeleteShipmentById(id int) error {
	if _, err := sss.db.Exec("DELETE FROM shipment WHERE id = $1", id); err != nil {
		log.Printf("failed to delete shipment with id: %d, error: %v", id, err)
		return fmt.Errorf("failed to delete shipment with id: %d, error: %v", id, err)
	}

	log.Printf("shipment with id: %d has been deleted", id)
	return nil
}

func (sss *sqlShipmentService) CheckBarcodeUnavailable(barcode string) bool {
	log.Printf("check barcode: %s availability", barcode)
	if len(barcode) > 25 || len(barcode) < 13 {
		return false
	}
	return sss.filter.Check(barcode)
}
