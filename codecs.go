package godays

import "encoding/json"

type TripStartedCodec struct{}

func (ts *TripStartedCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (ts *TripStartedCodec) Decode(data []byte) (interface{}, error) {
	var tripStarted TripStarted
	return &tripStarted, json.Unmarshal(data, &tripStarted)
}

type TripEndedCodec struct{}

func (ts *TripEndedCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (ts *TripEndedCodec) Decode(data []byte) (interface{}, error) {
	var tripEnded TripEnded
	return &tripEnded, json.Unmarshal(data, &tripEnded)
}

type LicenseTrackerCodec struct{}

func (ts *LicenseTrackerCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (ts *LicenseTrackerCodec) Decode(data []byte) (interface{}, error) {
	var licenseTracker LicenseTracker
	return &licenseTracker, json.Unmarshal(data, &licenseTracker)
}
