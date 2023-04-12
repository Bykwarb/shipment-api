package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
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
	returnJSONResponse(w, AvailabilityResponse{Unavailable: availability})
}

func (api *shipmentServer) CreateShipmentHandler(w http.ResponseWriter, r *http.Request) {
	var req ShipmentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		returnJSONResponse(w, Response{Message: "failed to decode JSON"})
		return
	}
	log.Print(req)
	shipment := shipments.Shipment{Sender: req.Sender, Receiver: req.Receiver, Origin: req.Origin, Destination: req.Destination, Barcode: req.Barcode, IsDelivered: req.IsDelivered, CreatedAt: req.CreatedAt}
	err = api.service.Save(&shipment)
	if err != nil {
		returnJSONResponse(w, Response{Message: err.Error()})
		return
	}
	returnJSONResponse(w, Response{Message: "shipment successfully created"})
}

func (api *shipmentServer) GetShipmentByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		returnJSONResponse(w, Response{Message: "failed to parse id parameter"})
		return
	}
	shipment, err := api.service.GetById(id)
	if err != nil {
		returnJSONResponse(w, Response{Message: err.Error()})
		return
	}
	returnJSONResponse(w, shipment)
}

func (api *shipmentServer) GetShipmentByBarcodeHandler(w http.ResponseWriter, r *http.Request) {
	barcode := r.URL.Query().Get("barcode")
	if barcode == "" {
		returnJSONResponse(w, Response{Message: "missing barcode parameter"})
		return
	}
	shipment, err := api.service.GetByBarcode(barcode)
	if err != nil {
		returnJSONResponse(w, Response{Message: err.Error()})
		return
	}
	returnJSONResponse(w, shipment)
}

func (api *shipmentServer) DeleteShipmentByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		returnJSONResponse(w, Response{Message: "failed to parse id parameter"})
		return
	}
	err = api.service.DeleteById(id)
	if err != nil {
		returnJSONResponse(w, Response{Message: "failed to delete shipment"})
		return
	}
	returnJSONResponse(w, Response{Message: "shipment successfully deleted"})
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

type ShipmentRequest struct {
	Barcode     string    `json:"barcode"`
	Sender      string    `json:"sender"`
	Receiver    string    `json:"receiver"`
	IsDelivered bool      `json:"is_delivered"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	CreatedAt   time.Time `json:"created_at"`
}

type Response struct {
	Message string `json:"message"`
}

type AvailabilityResponse struct {
	Unavailable bool `json:"unavailable"`
}
