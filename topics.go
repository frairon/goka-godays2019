package godays

import "github.com/lovoo/goka"

const (
	// TripStartedTopic contains the stream for pickup-events
	TripStartedTopic goka.Stream = "trip-started"
	// TripEndedTopic contains the stream for dropoff-events
	TripEndedTopic goka.Stream = "trip-ended"

	// LicenseConfigTopic (for task 3 and 4) to configure licenses
	LicenseConfigTopic goka.Stream = "configure-licenses"
)
