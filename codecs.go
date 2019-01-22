package godays

import "encoding/json"

// TripStartedCodec encodes/decodes TripStarted
type TripStartedCodec struct{}

// Encode encode message
func (ts *TripStartedCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes a message
func (ts *TripStartedCodec) Decode(data []byte) (interface{}, error) {
	var tripStarted TripStarted
	return &tripStarted, json.Unmarshal(data, &tripStarted)
}

// TripEndedCodec encodes/decodes TripEnded
type TripEndedCodec struct{}

// Encode encode message
func (ts *TripEndedCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes a message
func (ts *TripEndedCodec) Decode(data []byte) (interface{}, error) {
	var tripEnded TripEnded
	return &tripEnded, json.Unmarshal(data, &tripEnded)
}

// LicenseTrackerCodec encodes/decodes LicenseTracker
type LicenseTrackerCodec struct{}

// Encode encode message
func (ts *LicenseTrackerCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes a message
func (ts *LicenseTrackerCodec) Decode(data []byte) (interface{}, error) {
	var licenseTracker LicenseTracker
	return &licenseTracker, json.Unmarshal(data, &licenseTracker)
}

// TaxiStatusCodec encodes/decodes TaxiStatus
type TaxiStatusCodec struct{}

// Encode encode message
func (ts *TaxiStatusCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes a message
func (ts *TaxiStatusCodec) Decode(data []byte) (interface{}, error) {
	var TaxiStatus TaxiStatus
	return &TaxiStatus, json.Unmarshal(data, &TaxiStatus)
}

// LicenseConfigCodec encodes/decodes LicenseConfig
type LicenseConfigCodec struct{}

// Encode encode message
func (ts *LicenseConfigCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes a message
func (ts *LicenseConfigCodec) Decode(data []byte) (interface{}, error) {
	var licenseConfig LicenseConfig
	return &licenseConfig, json.Unmarshal(data, &licenseConfig)
}
