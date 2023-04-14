package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"task/iternal/abstraction"
	"task/iternal/model"

	"github.com/gorilla/mux"
)

type ShipmentApi struct {
	service abstraction.ShipmentService
}

func NewShipmentApi(service abstraction.ShipmentService) *ShipmentApi {
	return &ShipmentApi{service}
}

func (sa *ShipmentApi) CheckBarcodeAvailability(w http.ResponseWriter, r *http.Request) {
	barcode := mux.Vars(r)["barcode"]
	availability := sa.service.CheckBarcodeUnavailable(barcode)
	writeJSONResponse(w, model.AvailabilityResponse{Unavailable: availability})
}

func (sa *ShipmentApi) CreateShipment(w http.ResponseWriter, r *http.Request) {
	var req model.ShipmentRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeJSONResponse(w, model.Response{Message: "failed to decode JSON"})
		return
	}

	shipment := model.Shipment{
		Sender:      req.Sender,
		Receiver:    req.Receiver,
		Origin:      req.Origin,
		Destination: req.Destination,
		Barcode:     req.Barcode,
		IsDelivered: req.IsDelivered,
		CreatedAt:   req.CreatedAt,
	}

	err = sa.service.InsertShipment(&shipment)

	if err != nil {
		writeJSONResponse(w, model.Response{Message: err.Error()})
		return
	}

	writeJSONResponse(w, model.Response{Message: "shipment successfully created"})
}

func (sa *ShipmentApi) GetShipmentById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		writeJSONResponse(w, model.Response{Message: "failed to parse id parameter"})
		return
	}

	shipment, err := sa.service.SelectShipmentById(id)

	if err != nil {
		writeJSONResponse(w, model.Response{Message: err.Error()})
		return
	}

	writeJSONResponse(w, shipment)
}

func (sa *ShipmentApi) GetShipmentByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := r.URL.Query().Get("barcode")

	if barcode == "" {
		writeJSONResponse(w, model.Response{Message: "missing barcode parameter"})
		return
	}

	shipment, err := sa.service.SelectShipmentByBarcode(barcode)

	if err != nil {
		writeJSONResponse(w, model.Response{Message: err.Error()})
		return
	}

	writeJSONResponse(w, shipment)
}

func (sa *ShipmentApi) DeleteShipmentById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		writeJSONResponse(w, model.Response{Message: "failed to parse id parameter"})
		return
	}

	err = sa.service.DeleteShipmentById(id)

	if err != nil {
		writeJSONResponse(w, model.Response{Message: "failed to delete shipment"})
		return
	}

	writeJSONResponse(w, model.Response{Message: "shipment successfully deleted"})
}

func writeJSONResponse(w http.ResponseWriter, v interface{}) {
	json, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
