package godays

import "encoding/json"

type TripStartedCodec struct{}

func (ts *TripStartedCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (ts *TripStartedCodec) Decode(data []byte) (interface{}, error) {
	var tripStarted TripStarted
	return &ts, json.Unmarshal(data, &tripStarted)
}

type TripEndedCodec struct{}

func (ts *TripEndedCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (ts *TripEndedCodec) Decode(data []byte) (interface{}, error) {
	var tripEnded TripEnded
	return &ts, json.Unmarshal(data, &tripEnded)
}
