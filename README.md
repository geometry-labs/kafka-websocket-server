# Kafka Websocket Server

This service can be used alongside any kafka cluster to expose selected topics to a websocket server.

## Local build
```
go build -o main .
```

## Docker build
```
docker build . -t kafka-websocket-server:latest
docker run \
  -p "8080:8080" \
  -e KAFKA_WEBSOCKET_SERVER_TOPICS="blocks,transactions,logs"
  -e KAFKA_WEBSOCKET_SERVER_BROKER_URL="kafka:9092"
  -e KAFKA_WEBSOCKET_SERVER_PORT="8080"
  kafka-websocket-server:latest
```

## Enviroment Variables

| Name | Description | Default | Required |
|------|-------------|---------|----------|
| KAFKA_WEBSOCKET_SERVER_TOPICS | comma seperated list of topic names | NULL | True |
| KAFKA_WEBSOCKET_SERVER_BROKER_URL | location of broker | NULL | True |
| KAFKA_WEBSOCKET_SERVER_PORT | port to expose for websocket connections | "8080" | False |
