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
	barcode := mux.Vars(r)["barcode"]
	availability := api.service.CheckBarcodeAvailability(barcode)
	returnJSONResponse(w, availabilityResponse{availability})
}

func (api *shipmentServer) CreateShipmentHandler(w http.ResponseWriter, r *http.Request) {
	var req shipmentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	shipment := shipments.Shipment{Sender: req.sender, Receiver: req.receiver, Origin: req.origin, Destination: req.destination, Barcode: req.barcode}
	api.service.Save(&shipment)
	returnJSONResponse(w, response{message: "shipments successfully created"})
}

func (api *shipmentServer) GetShipmentByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	shipment, err := api.service.GetById(id)
	if err != nil {
		returnJSONResponse(w, err)
	} else {
		returnJSONResponse(w, shipment)
	}
}

func (api *shipmentServer) GetShipmentByBarcodeHandler(w http.ResponseWriter, r *http.Request) {
	barcode := r.URL.Query().Get("barcode")
	shipment, _ := api.service.GetByBarcode(barcode)
	returnJSONResponse(w, shipment)
}

func (api *shipmentServer) DeleteShipmentByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	err = api.service.DeleteById(id)
	if err != nil {
		returnJSONResponse(w, err)
	} else {
		returnJSONResponse(w, response{message: "shipment successfully deleted"})
	}
}

func returnJSONResponse(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type shipmentRequest struct {
	barcode       string    `json:"barcode"`
	sender        string    `json:"sender"`
	receiver      string    `json:"receiver"`
	isDelivered   bool      `json:"isDelivered"`
	origin        string    `json:"origin"`
	destination   string    `json:"destination"`
	departureDate time.Time `json:"departureDate"`
}

type response struct {
	message string `json:"message"`
}

type availabilityResponse struct {
	isAvailable bool `json:"isAvailable"`
}
