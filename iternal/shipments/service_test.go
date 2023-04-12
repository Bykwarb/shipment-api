package shipments

import (
	"database/sql"
	_ "database/sql/driver"
	"task/config"
	"task/iternal/bloom"
	"task/iternal/database"
	"testing"
)

var service *sqlShipmentService
var filter *bloom.Filter
var db *sql.DB

func init() {
	conf := config.LoadConfig("config.yml")
	db = database.OpenConnection(conf)
	arrSize, _ := bloom.CalculateArraySize(conf.Filter.ExpectedNumElements, conf.Filter.FalsePositiveProbability)
	filter = bloom.NewFilterWithDefaultHash(conf.Filter.ExpectedNumElements, arrSize)
	service = NewSQLShipmentService(db, filter)
}
func TestSqlShipmentService_GetById(t *testing.T) {
	shipment, _ := service.GetById(1)

}

func TestSqlShipmentService_GetByBarcode(t *testing.T) {
}

func TestSqlShipmentService_CheckBarcodeAvailability(t *testing.T) {

}

func TestSqlShipmentService_DeleteById(t *testing.T) {

}

func TestSqlShipmentService_Save(t *testing.T) {

}

func TestShipment_GenerateBarcode(t *testing.T) {

}
