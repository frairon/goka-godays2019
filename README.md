# Goka
## Painless stream processing with Go and Kafka

## Agenda
* Setup
  * getting the dependencies
  * starting kafka locally
* (examples with pure Sarama)
* Goka concepts
  * starting an emitter
  * starting a producer
  * starting a view
  * some monitoring
  * testing our code


## Event Data
* a tiny data set is located in testdata/taxidata_tiny.csv
* a 100k dataset can be loaded from https://storage.googleapis.com/lv-goka-godays2019/taxidata_100k.csv
* the full dataset (~320MB) can be downloaded from https://storage.googleapis.com/lv-goka-godays2019/taxidata_complete.csv
