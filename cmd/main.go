package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
	"task/api"
	"task/config"
	"task/iternal/bloom"
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

var filter *bloom.Filter
var db *sql.DB
var conf *config.Config

func init() {
	conf = config.LoadConfig("config.yml")
	db = database.OpenConnection(conf)
	db.SetMaxOpenConns(100)
	if conf.Filter.Enabled {
		expectedNumElements := conf.Filter.ExpectedNumElements
		falsePositiveProbability := conf.Filter.FalsePositiveProbability
		filterArraySize, err := bloom.CalculateArraySize(expectedNumElements, falsePositiveProbability)
		if err != nil {
			panic(err)
		}
		filter = bloom.NewFilterWithDefaultHash(expectedNumElements, filterArraySize)
		var wg sync.WaitGroup
		var mutex sync.Mutex
		wg.Add(100)
		fillFilter(&wg, &mutex, filter)
		wg.Wait()

	}
}

func main() {
	defer database.CloseConnection(db)
	service := shipments.NewSQLShipmentService(db, filter)
	server := api.NewShipmentServer(service)
	log.Println(fmt.Sprintf("server is started in %s:%s", conf.Server.Host, conf.Server.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port), createRoute(server)))
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

func fillFilter(wg *sync.WaitGroup, mutex *sync.Mutex, filter *bloom.Filter) {
	num := conf.Filter.ExpectedNumElements
	numRanges := 100
	rangeSize := num / numRanges
	for i := 0; i < numRanges; i++ {
		start := i * rangeSize
		end := (i + 1) * rangeSize
		if i == numRanges-1 {
			end = num
		}
		go func() {
			defer wg.Done()
			rows, err := db.Query("SELECT barcode FROM shipments LIMIT $1 OFFSET $2", end, start)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				var barcode string
				err := rows.Scan(&barcode)
				if err != nil {
					log.Fatal(err)
				}
				mutex.Lock()
				filter.AddToFilter(barcode)
				mutex.Unlock()
			}
		}()
	}

}
