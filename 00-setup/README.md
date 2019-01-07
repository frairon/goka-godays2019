# Set up your working environment

* Golang https://golang.org/dl/
* Docker/Git
  * `sudo apt-get install docker git-core`
  * `brew install docker git-core`

* Setup a Gopath (if you don't have it yet)

```
mkdir -p $HOME/gocode/
cd $HOME/gocode
export $GOPATH=$HOME/gocode
```

* Get the Workshop Code

```
# need to have $GOPATH set
go get github.com/frairon/goka-godays2019
```


* Start Kafka locally

```
cd $GOPATH/src/github.com/frairon/goka-godays2019
# start the kafka locally
docker run --name=kafka-cluster1 -d  -p 2181:2181 -p 9092:9092 --env ADVERTISED_HOST=127.0.0.1 --env ADVERTISED_PORT=9092 spotify/kafka

```


## Now you're ready to start working with the code
