package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"task/config"
	"task/iternal/api/handler"
	"task/iternal/api/router"
	"task/iternal/database"
	"task/iternal/model"
	"task/iternal/service/shipment"
	"task/pkg/bloom"
)

var db *sql.DB
var conf *model.Config
var filter *bloom.Filter

func Run(configFilePath string) {
	conf = config.LoadConfig(configFilePath)
	db = database.OpenConnection(conf)
	defer database.CloseConnection(db)
	db.SetMaxOpenConns(100)
	setUpFilter()
	service := shipment.NewSQLShipmentService(db, filter)
	server := handler.NewShipmentApi(service)
	log.Println(fmt.Sprintf("server is started in %s:%s", conf.Server.Host, conf.Server.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port), router.CreateRouter(*server)))
}

func setUpFilter() {
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
	fillFilterFromSqlDB(&wg, &mutex)
	wg.Wait()
}

func fillFilterFromSqlDB(wg *sync.WaitGroup, mutex *sync.Mutex) {
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
