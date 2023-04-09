package main

import (
	"fmt"
	"task/shipments"
)

func main() {
	c := shipments.GenerateBarcode("Tokyo", "Ghoul")
	fmt.Println(c)
}
