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
	service := shipments.SqlShipmentService{DB: db}
	shipment, err := service.GetShipmentById(2)
	fmt.Println(err)
	fmt.Println(shipment)
	defer database.CloseConnection(db)
}
