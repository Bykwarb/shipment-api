package abstraction

import "net/http"

type Api interface {
	CheckBarcodeAvailability(w http.ResponseWriter, r *http.Request)
	CreateShipment(w http.ResponseWriter, r *http.Request)
	GetShipmentById(w http.ResponseWriter, r *http.Request)
	GetShipmentByBarcode(w http.ResponseWriter, r *http.Request)
	DeleteShipmentById(w http.ResponseWriter, r *http.Request)
}
