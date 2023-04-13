package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"task/iternal/shipments"
	"testing"
	"time"
)

func TestShipmentServer_CheckBarcodeAvailability(t *testing.T) {
	barcode := "TESTBARCODE123"
	expectedAvailability := true

	service := &MockService{}
	server := NewShipmentServer(service)
	req, err := http.NewRequest("GET", fmt.Sprintf("api/v1/barcodes/%s/availability", barcode), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(server.CheckBarcodeAvailability)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp AvailabilityResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	if resp.Unavailable != expectedAvailability {
		t.Errorf("handler returned wrong availability: got %v want %v", resp.Unavailable, expectedAvailability)
	}
}

func TestShipmentServer_GetShipmentByIdHandler(t *testing.T) {
	expectedResp := *shipments.NewShipment("TEST", "TEST", "TEST", "TEST")
	expectedResp.Barcode = "TEST1003TEST1"
	expectedResp.Id = 1
	service := &MockService{}
	server := NewShipmentServer(service)
	id := 1
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/shipment/{%d}", id), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(server.GetShipmentByIdHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println(rr.Body)

	var resp shipments.Shipment
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	if resp != expectedResp {
		t.Errorf("handler returned wrong shipment: got %v want %v", resp, expectedResp)
	}
}

func TestShipmentServer_GetShipmentByBarcodeHandler(t *testing.T) {
	expectedResp := *shipments.NewShipment("TEST", "TEST", "TEST", "TEST")
	expectedResp.Barcode = "TEST1003TEST1"
	expectedResp.CreatedAt = time.Time{}
	service := &MockService{}
	server := NewShipmentServer(service)
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/shipment?barcode=%s", "TEST1003TEST1"), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(server.GetShipmentByBarcodeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp shipments.Shipment
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	if resp != expectedResp {
		t.Errorf("handler returned wrong shipment: got %v want %v", resp, expectedResp)
	}
}

func TestShipmentServer_CreateShipmentHandler(t *testing.T) {
	expectedReq := *shipments.NewShipment("TEST", "TEST", "TEST", "TEST")
	expectedReq.Barcode = "TEST1003TEST1"
	service := &MockService{}
	server := NewShipmentServer(service)
	jsonBody, err := json.Marshal(expectedReq)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req, err := http.NewRequest("POST", "api/v1/shipment", strings.NewReader(string(jsonBody)))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(server.CreateShipmentHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	expectedRes := Response{Message: "shipment successfully created"}
	if resp != expectedRes {
		t.Errorf("handler returned wrong shipment: got %v want %v", resp, expectedReq)
	}
}

func TestShipmentServer_DeleteShipmentByIdHandler(t *testing.T) {
	service := &MockService{}
	server := NewShipmentServer(service)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/shipment/%d", 1), nil)

	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(server.DeleteShipmentByIdHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	expectedRes := Response{Message: "shipment successfully deleted"}
	if resp != expectedRes {
		t.Errorf("handler returned wrong shipment: got %v want %v", resp, expectedRes)
	}
}

type MockService struct {
}

func (service *MockService) Save(shipment *shipments.Shipment) error {
	return nil
}
func (service *MockService) GetById(id int) (*shipments.Shipment, error) {
	shipment := shipments.NewShipment("TEST", "TEST", "TEST", "TEST")
	shipment.Barcode = "TEST1003TEST1"
	shipment.CreatedAt = time.Time{}
	shipment.Id = 1
	return shipment, nil
}
func (service *MockService) GetByBarcode(barcode string) (*shipments.Shipment, error) {
	shipment := shipments.NewShipment("TEST", "TEST", "TEST", "TEST")
	shipment.Barcode = "TEST1003TEST1"
	shipment.CreatedAt = time.Time{}
	return shipment, nil
}
func (service *MockService) DeleteById(id int) error {
	return nil
}
func (service *MockService) CheckBarcodeUnavailable(barcode string) bool {
	return true
}
