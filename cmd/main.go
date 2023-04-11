package main

import (
	"fmt"
	"task/config"
	"task/iternal/database"
	"task/iternal/shipments"
)

func init() {

}
func main() {
	c := config.LoadConfig("config.yml")
	db := database.OpenConnection(c)
	service := shipments.NewShipmentService(db)
	shipment := shipments.NewShipment("Bogdan", "Maybebaby", "Kyiv", "Tokyo")
	shipment.Barcode = "KV155784865744TO"
	err := service.SaveShipment(shipment)
	fmt.Println(err)
	defer database.CloseConnection(db)
}
