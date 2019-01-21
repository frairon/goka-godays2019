# Taxi Status

Consume a stream from Kafka and store it as state tracking the current state of a taxi

We have two components:

* Processor "taxiTrips": reads the startTrip/endTrip events and stores them in its state
* View: provide access to the state table of "taxiTrips" via REST interface
  * do `curl "localhost:8080/taxi?id=taxi-6"` to see if `taxi-6` is currently busy
