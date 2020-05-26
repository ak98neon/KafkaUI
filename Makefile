hello:
	echo "Hello"

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main ./src

run:
	go run bin/main

docker:
	docker build -t kafka_ui .
	docker tag kafka_ui akudria/kafka_ui:latest
	docker push akudria/kafka_ui:latest
