package model

import "time"

type ShipmentRequest struct {
	Barcode     string    `json:"barcode"`
	Sender      string    `json:"sender"`
	Receiver    string    `json:"receiver"`
	IsDelivered bool      `json:"isDelivered"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Response struct {
	Message string `json:"message"`
}

type AvailabilityResponse struct {
	Unavailable bool `json:"unavailable"`
}
