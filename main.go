package main

import (
	"fmt"
	"task/shipments"
)

func init() {

}
func main() {
	s := shipments.NewShipment("Bogdan", "Valera228", "BOMBASS", "Abebrus")
	s.GenerateBarcode()
	fmt.Println(s.Barcode)
	fmt.Println(s)
}
