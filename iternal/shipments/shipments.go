package shipments

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Shipment struct {
	Id            int
	Barcode       string
	Sender        string
	Receiver      string
	IsDelivered   bool
	Origin        string
	Destination   string
	DepartureDate time.Time
}

func NewShipment(sender, receiver, from, to string) *Shipment {
	return &Shipment{
		Sender:        sender,
		Receiver:      receiver,
		Origin:        from,
		Destination:   to,
		IsDelivered:   false,
		DepartureDate: time.Now(),
	}
}

func (shipment *Shipment) GenerateBarcode() {
	suffix := strings.ToUpper(string(shipment.Origin[0]) + string(shipment.Origin[len(shipment.Origin)-1]))

	var body string
	for i := 0; i < rand.Intn(21-9+1)+9; i++ {
		body += strconv.Itoa(rand.Intn(9))
	}

	prefix := strings.ToUpper(string(shipment.Destination[0]) + string(shipment.Destination[len(shipment.Destination)-1]))
	shipment.Barcode = fmt.Sprintf("%s%s%s", suffix, body, prefix)
}
