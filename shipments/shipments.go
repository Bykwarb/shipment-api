package shipments

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Shipment struct {
	Barcode       string
	Sender        string
	Receiver      string
	IsDelivered   string
	From          string
	To            string
	DepartureDate time.Time
}

func GenerateBarcode(from, to string) string {
	suffix := strings.ToUpper(string(from[0]) + string(from[len(from)-1]))

	var body string
	for i := 0; i < rand.Intn(21-9+1)+9; i++ {
		body += strconv.Itoa(rand.Intn(9))
	}

	prefix := strings.ToUpper(string(to[0]) + string(to[len(to)-1]))

	return fmt.Sprintf("%s%s%s", suffix, body, prefix)
}
