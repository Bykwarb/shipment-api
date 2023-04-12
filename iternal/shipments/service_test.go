package shipments

import (
	"database/sql"
	_ "database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"task/iternal/bloom"
	"testing"
	"time"
)

func TestSqlShipmentService_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "barcode", "sender", "receiver", "is_delivered", "origin", "destination", "created_at"}).
		AddRow(1, "TESTCODE1003", "Test", "Test", false, "Test", "Test", time.Time{})
	mock.ExpectQuery("SELECT id, barcode, sender, receiver, is_delivered, origin, destination, created_at FROM shipments WHERE id = .*").
		WillReturnRows(rows).WithArgs(1)
	service := NewSQLShipmentService(db, &bloom.Filter{})
	s, err := service.GetById(1)
	checkEquals(s, t)
}

func TestSqlShipmentService_GetByBarcode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "barcode", "sender", "receiver", "is_delivered", "origin", "destination", "created_at"}).
		AddRow(1, "TESTCODE1003", "Test", "Test", false, "Test", "Test", time.Time{})
	mock.ExpectQuery("SELECT id, barcode, sender, receiver, is_delivered, origin, destination, created_at FROM shipments WHERE barcode = .*").
		WillReturnRows(rows).WithArgs("TESTCODE1003")
	service := NewSQLShipmentService(db, &bloom.Filter{})
	s, err := service.GetByBarcode("TESTCODE1003")
	checkEquals(s, t)
}

func TestSqlShipmentService_CheckBarcodeAvailability(t *testing.T) {
	e := 100
	arrSize, _ := bloom.CalculateArraySize(e, 0.1)
	filter := bloom.NewFilterWithDefaultHash(e, arrSize)
	service := NewSQLShipmentService(&sql.DB{}, filter)
	available := service.CheckBarcodeUnavailable("TEST")
	if available {
		t.Errorf("Excpected: %t, got: %t", false, available)
	}
	filter.AddToFilter("TEST123456789")
	available = service.CheckBarcodeUnavailable("TEST123456789")
	if !available {
		t.Errorf("Excpected: %t, got: %t", true, available)
	}
	available = service.CheckBarcodeUnavailable("TEST1337228321")
	if available {
		t.Errorf("Excpected: %t, got: %t", false, available)
	}
}

func TestSqlShipmentService_DeleteById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM shipments WHERE id = .*").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	service := NewSQLShipmentService(db, &bloom.Filter{})
	err = service.DeleteById(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSqlShipmentService_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	shipment := &Shipment{
		Barcode:     "TESTBARCODE123",
		Sender:      "Test Sender",
		Receiver:    "Test Receiver",
		IsDelivered: false,
		Origin:      "Test Origin",
		Destination: "Test Destination",
		CreatedAt:   time.Now(),
	}

	filter := bloom.NewFilterWithDefaultHash(100, 5000)
	service := NewSQLShipmentService(db, filter)

	mock.ExpectExec("INSERT INTO shipments").WithArgs(
		shipment.Barcode,
		shipment.Sender,
		shipment.Receiver,
		shipment.IsDelivered,
		shipment.Origin,
		shipment.Destination,
		shipment.CreatedAt,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	// Test the method
	err = service.Save(shipment)
	if err != nil {
		t.Errorf("error was not expected while saving shipment: %v", err)
	}

	// Check the result
	if filter.Check(shipment.Barcode) == false {
		t.Errorf("shipment barcode was not added to the filter")
	}
}

func TestShipment_GenerateBarcode(t *testing.T) {
	shipment := &Shipment{
		Origin:      "New York",
		Destination: "Los Angeles",
	}

	shipment.GenerateBarcode()

	if len(shipment.Barcode) > 25 || len(shipment.Barcode) < 13 {
		t.Errorf("barcode length must be <= then 25 and >= 13, barcode: %s", shipment.Barcode)
	}

	// Check if prefix and suffix are correct
	expectedPrefix := "LS"
	if shipment.Barcode[len(shipment.Barcode)-2:] != expectedPrefix {
		t.Errorf("generated prefix is incorrect, expected %s but got %s", expectedPrefix, shipment.Barcode[len(shipment.Barcode)-2:])
	}

	expectedSuffix := "NK"
	if shipment.Barcode[:2] != expectedSuffix {
		t.Errorf("generated suffix is incorrect, expected %s but got %s", expectedSuffix, shipment.Barcode[:2])
	}
}

func checkEquals(s *Shipment, t *testing.T) {
	expectedShipment := &Shipment{
		Id:          1,
		Barcode:     "TESTCODE1003",
		Sender:      "Test",
		Receiver:    "Test",
		IsDelivered: false,
		Origin:      "Test",
		Destination: "Test",
		CreatedAt:   time.Time{}}

	if s.Id != expectedShipment.Id {
		t.Errorf("Expected Id: %d, got: %d", expectedShipment.Id, s.Id)
	}
	if s.Barcode != expectedShipment.Barcode {
		t.Errorf("Expected Barcode: %s, got: %s", expectedShipment.Barcode, s.Barcode)
	}
	if s.Sender != expectedShipment.Sender {
		t.Errorf("Expected Sender: %s, got: %s", expectedShipment.Sender, s.Sender)
	}
	if s.Receiver != expectedShipment.Receiver {
		t.Errorf("Expected Receiver: %s, got: %s", expectedShipment.Receiver, s.Receiver)
	}
	if s.IsDelivered != expectedShipment.IsDelivered {
		t.Errorf("Expected IsDelivered: %v, got: %v", expectedShipment.IsDelivered, s.IsDelivered)
	}
	if s.Origin != expectedShipment.Origin {
		t.Errorf("Expected Origin: %s, got: %s", expectedShipment.Origin, s.Origin)
	}
	if s.Destination != expectedShipment.Destination {
		t.Errorf("Expected Destination: %s, got: %s", expectedShipment.Destination, s.Destination)
	}
	if s.CreatedAt != expectedShipment.CreatedAt {
		t.Errorf("Expected CreatedAt: %v, got: %v", expectedShipment.CreatedAt, s.CreatedAt)
	}

}
