package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"task/iternal/mock"
	"task/iternal/model"
)

func TestShipmentServer_CheckBarcodeAvailability(t *testing.T) {
	barcode := "TESTBARCODE123"
	expectedAvailability := true

	server := NewShipmentApi(&mock.ShipmentService{})

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

	var resp model.AvailabilityResponse

	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if resp.Unavailable != expectedAvailability {
		t.Errorf("handler returned wrong availability: got %v want %v", resp.Unavailable, expectedAvailability)
	}
}

func TestShipmentServer_GetShipmentByIdHandler(t *testing.T) {
	expectedResp := model.NewShipment("TEST", "TEST", "TEST", "TEST")
	expectedResp.Barcode = "TEST1003TEST1"
	expectedResp.Id = 1
	expectedResp.CreatedAt = time.Time{}

	server := NewShipmentApi(&mock.ShipmentService{})

	req, err := http.NewRequest("GET", "/api/v1/shipment/1", nil)

	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/shipment/{id}", server.GetShipmentById).Methods("GET")
	router.StrictSlash(true)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp model.Shipment
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if !reflect.DeepEqual(resp, *expectedResp) {
		t.Errorf("handler returned wrong shipment: got %v want %v", resp, expectedResp)
	}
}

func TestShipmentServer_GetShipmentByBarcodeHandler(t *testing.T) {
	expectedResp := *model.NewShipment("TEST", "TEST", "TEST", "TEST")
	expectedResp.Barcode = "TEST1003TEST1"
	expectedResp.CreatedAt = time.Time{}

	server := NewShipmentApi(&mock.ShipmentService{})

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/shipment?barcode=%s", "TEST1003TEST1"), nil)

	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(server.GetShipmentByBarcode)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp model.Shipment

	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if resp != expectedResp {
		t.Errorf("handler returned wrong shipment: got %v want %v", resp, expectedResp)
	}
}

func TestShipmentServer_CreateShipmentHandler(t *testing.T) {
	expectedReq := *model.NewShipment("TEST", "TEST", "TEST", "TEST")
	expectedReq.Barcode = "TEST1003TEST1"

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

	handler := http.HandlerFunc(NewShipmentApi(&mock.ShipmentService{}).CreateShipment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp model.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	expectedRes := model.Response{Message: "shipment successfully created"}

	if resp != expectedRes {
		t.Errorf("handler returned wrong shipment: got %v want %v", resp, expectedReq)
	}
}

func TestShipmentServer_DeleteShipmentByIdHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/shipment/%d", 1), nil)

	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/shipment/{id}", NewShipmentApi(&mock.ShipmentService{}).DeleteShipmentById).Methods("DELETE")
	router.StrictSlash(true)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp model.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	expectedRes := model.Response{Message: "shipment successfully deleted"}
	if resp != expectedRes {
		t.Errorf("handler returned wrong shipment: got %v want %v", resp, expectedRes)
	}
}
