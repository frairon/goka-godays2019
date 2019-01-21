

restart-kafka:
	# remove disk caches
	-rm -r /tmp/goka
	# stop old versions of the container
	-docker stop kafka-cluster1
	-docker rm kafka-cluster1
	docker run --name=kafka-cluster1 -d  -p 2181:2181 -p 9092:9092 --env ADVERTISED_HOST=127.0.0.1 --env NUM_PARTITIONS=10 --env ADVERTISED_PORT=9092 spotify/kafka


get-100:
	cd testdata && wget https://storage.googleapis.com/lv-goka-godays2019/taxidata_100k.csv

get-complete:
	cd testdata && wget https://storage.googleapis.com/lv-goka-godays2019/taxidata_complete.csv
