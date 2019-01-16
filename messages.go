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

type LicenseTracker struct {
	Started time.Time       `json:"time_started"`
	Ended   time.Time       `json:"time_ended"`
	Taxis   map[string]bool `json:"taxis"`
	Fraud   bool            `json:"fraud"`
}

// TaxiTrip stores the current trip status of a taxi
type TaxiTrip struct {
	TaxiID    string `json:"taxi_id"`
	LicenseID string `json:"license_id"`

	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
}
