# Assignment _Taxi Status_

Goal is to track a taxi's state and present that in a simple web interface.

## Tasks

* define a processor
  * add the input streams TripStartedTopic and TripStartedCodec (resp. for TripEnded)
  * add persistence definition using TaxiStatusCodec
* use the view to TripTrackerGroup
  * the view is already defined, just use it in the http getter


## Run it
Make sure Kafka is running, otherwise do `make restart-kafka` in the root.

* run the tracker `go run 02-taxistatus/main.go`
* run the emitter to get some events like `go run 00-emitter/main.go --input testdata/taxidata_100k.csv
* run `curl "localhost:8080/taxi?id=taxi-6"` to see if `taxi-6` is currently busy and what's its status.
* or just open http://localhost:8080/taxi?id=taxi-6
