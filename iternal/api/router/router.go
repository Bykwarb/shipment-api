package router

import (
	"task/iternal/api/handler"

	"github.com/gorilla/mux"
)

func CreateRouter(a handler.ShipmentApi) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/api/v1/shipment/{id}", a.GetShipmentById).Methods("GET")
	router.HandleFunc("/api/v1/shipment/{id}", a.DeleteShipmentById).Methods("DELETE")
	router.HandleFunc("/api/v1/shipment", a.GetShipmentByBarcode).Methods("GET")
	router.HandleFunc("/api/v1/barcodes/{barcode}/availability", a.CheckBarcodeAvailability).Methods("GET")
	router.HandleFunc("/api/v1/shipment", a.CreateShipment).Methods("POST")
	return router
}
