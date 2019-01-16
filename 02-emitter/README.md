# Taxi Data Emitter

This pumps the raw taxi events into the two Kafka topics *trip-started* and *trip-ended*.

To get some timing, it allows to time lapse the event piping to avoid waiting 20 days (the whole data time frame) or be flooded with events.

```
# send events where 1 minute of events is 1 second in realtime
godays-godays2019  $ go run 02-emitter/main.go --time-lapse 60


```
