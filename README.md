# Godays 2019 - Painless stream processing with Go and Kafka

Slides available [here](https://github.com/frairon/goka-godays2019/blob/master/slides.pdf).

## Agenda
* Checkout the [data](https://github.com/frairon/goka-godays2019/blob/master/testdata/README.md)
* Setup working environment
* clone the repo
* tasks
  1. Consume events using Sarama
  1. Track and view state of taxis
  1. alert on license fraud
  1. detect fraud
  1. extras (testing, monitoring, visualization of coordinate-cells)

## Set up your working environment

* Golang https://golang.org/dl/
* Docker/Git
  * Mac:
```bash
brew install git
brew cask install docker
open /Applications/Docker.app
```

  * Linux:
```bash
sudo apt-get install docker git
```

* Setup a Gopath (if you don't have it yet)
```bash
mkdir -p $HOME/gocode/
cd $HOME/gocode
export GOPATH=$HOME/gocode
```

* Get the Workshop Code
```bash
# need to have $GOPATH set
mkdir -p $GOPATH/src/github.com/frairon/goka-godays2019
git clone https://github.com/frairon/goka-godays2019.git
# OR
go get github.com/frairon/goka-godays2019
```

* Get the workshop data
```bash
# tiny data set is located in testdata/taxidata_tiny.csv
# 100k dataset can be loaded from https://storage.googleapis.com/lv-goka-godays2019/taxidata_100k.csv
make get-100

# [Optional]
# full dataset (~320MB) can be downloaded from https://storage.googleapis.com/lv-goka-godays2019/taxidata_complete.csv
make get-complete
```

* Start Kafka locally
```bash
make restart-kafka
```



## Now you're ready to start working with the code
