# UI example

Shows the number of activity on the map based on the taxi data in web using leaflet.

It's a showcase for the view-iterator, super dirty and a questionable use case.

## Run it

* make sure Kafka and an emitter is running
* run `go run 05-cells/main.go`


## TODO
* the page should iterate on the initial load only, then web sockets to send cell-updates
