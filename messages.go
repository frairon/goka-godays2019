package godays

import "time"

type TripStarted struct {
	Ts        time.Time `json:"time"`
	TaxiID    string    `json:"taxi_id"`
	LicenseID string    `json:"license_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type TripEnded struct {
	Ts        time.Time `json:"time"`
	TaxiID    string    `json:"taxi_id"`
	LicenseID string    `json:"license_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`

	Duration time.Duration `json:"duration"`
	Distance float64       `json:"distance"`

	Charge float64 `json:"charge"`
	Tip    float64 `json:"tip"`
}
