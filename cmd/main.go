package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"task/api"
	"task/config"
	"task/iternal/database"
	"task/iternal/shipments"
)

type Server interface {
	CheckBarcodeAvailability(w http.ResponseWriter, r *http.Request)
	CreateShipmentHandler(w http.ResponseWriter, r *http.Request)
	GetShipmentByIdHandler(w http.ResponseWriter, r *http.Request)
	GetShipmentByBarcodeHandler(w http.ResponseWriter, r *http.Request)
	DeleteShipmentByIdHandler(w http.ResponseWriter, r *http.Request)
}

func createRoute(server Server) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/api/v1/shipment/{id}", server.GetShipmentByIdHandler).Methods("GET")
	router.HandleFunc("/api/v1/shipment/{id}", server.DeleteShipmentByIdHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/shipment", server.GetShipmentByBarcodeHandler).Methods("GET")
	router.HandleFunc("/api/v1/barcodes/{barcode}/availability", server.CheckBarcodeAvailability).Methods("GET")
	router.HandleFunc("/api/v1/shipment", server.CreateShipmentHandler).Methods("POST")
	return router
}

func main() { // Создаем канал для принятия сигналов
	c := config.LoadConfig("config.yml")
	db := database.OpenConnection(c)
	defer database.CloseConnection(db)
	service := shipments.NewShipmentService(db)
	server := api.NewShipmentServer(service)
	log.Println(fmt.Sprintf("server is started in %s:%s", c.Server.Host, c.Server.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port), createRoute(server)))
}
