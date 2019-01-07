# Set up your working environment


```

mkdir $HOME/workshop
cd $HOME/workshop
export $GOPATH=`pwd`

# start the kafka locally
docker run -d --env ADVERTISED_PORT=9092 spotify/kafka

```
