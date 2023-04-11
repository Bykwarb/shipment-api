package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"task/iternal/shipments"
	"time"
)

type shipmentServer struct {
	service shipments.ShipmentService
}

func NewShipmentServer(service shipments.ShipmentService) *shipmentServer {
	return &shipmentServer{service}
}

func (api *shipmentServer) CheckBarcodeAvailability(w http.ResponseWriter, r *http.Request) {

}

func (api *shipmentServer) CreateShipmentHandler(w http.ResponseWriter, r *http.Request) {
	var req shipmentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	shipment := shipments.NewShipment(req.Sender, req.Receiver, req.Origin, req.Destination)
	shipment.Barcode = req.Barcode
	api.service.SaveShipment(shipment)
	returnJsonResponse(w, "shipments successfully created")
}

func (api *shipmentServer) GetShipmentByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	shipment, err := api.service.GetShipmentById(id)
	if err != nil {
		returnJsonResponse(w, err)
	} else {
		returnJsonResponse(w, shipment)
	}
}

func (api *shipmentServer) GetShipmentByBarcodeHandler(w http.ResponseWriter, r *http.Request) {
	barcode := r.URL.Query().Get("barcode")
	shipment, _ := api.service.GetShipmentByBarcode(barcode)
	returnJsonResponse(w, shipment)
}

func (api *shipmentServer) DeleteShipmentByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	err = api.service.DeleteShipmentById(id)
	if err != nil {
		returnJsonResponse(w, err)
	} else {
		returnJsonResponse(w, "shipment successfully deleted")
	}
}

func returnJsonResponse(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type shipmentRequest struct {
	Barcode       string
	Sender        string
	Receiver      string
	IsDelivered   bool
	Origin        string
	Destination   string
	DepartureDate time.Time
}
