package api

import (
	"encoding/json"
	"net/http"
	"task/iternal/shipments"
)

type shipmentApi struct {
	service *shipments.ShipmentService
}

func NewShipmentApi(service shipments.ShipmentService) *shipmentApi {
	return &shipmentApi{&service}
}

func (api *shipmentApi) checkBarcodeAvailability(w http.ResponseWriter, r *http.Request) {

}

func (api *shipmentApi) shipmentHandler(w http.ResponseWriter, r *http.Request) {

}

func returnJson(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
