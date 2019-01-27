# _Fraud Licenses_

Goal is to detect the usage of fraud taxi licenses and print some notification.
In this task we will mark licenses as fraud manually using
a simple helper tool in `cmd/stringproducer` that allows us to send arbitrary string messages to a topic.

## Tasks

* Processor 1 saves the LicenseConfig for each license
  * define the processor
  * implement the consume function so that it simply stores the incoming message as state (for the sake of simplicity)
* Processor 2 consumes trip-started events and notifies via print when a bad license is used.
  * define the processor
  * implement a notification (simple print) when every it notices a fraud license being used

## Run it
* make sure Kafka is running
* make sure you have an events emitter running (see previous examples)

* Run the processors with `go run 03-fraud-license/main.go`
* Emit something to the configure topic manually `go run cmd/stringemitter/main.go --topic configure-licenses --key=license-5 --value='{"fraud":true}'`
* If you write something invalid, your processor will die
  * Side task: tolerate an invalid message
