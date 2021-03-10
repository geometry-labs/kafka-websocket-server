# Kafka Websocket Server

This service can be used alongside any kafka cluster to expose selected topics to a websocket server.

Docker hub [image](https://hub.docker.com/r/pranavt61/kafka-websocket-server)

## Local build
```
go build -o main .
```

## Docker build
```
docker build . -t kafka-websocket-server:latest --target prod
docker run \
  -p "8080:8080" \
  -e KAFKA_WEBSOCKET_SERVER_TOPICS="blocks,transactions,logs"
  -e KAFKA_WEBSOCKET_SERVER_BROKER_URL="kafka:9092"
  -e KAFKA_WEBSOCKET_SERVER_PORT="8080"
  -e KAFKA_WEBSOCKET_SERVER_PREFIX="/ws"
  kafka-websocket-server:latest
```

## Enviroment Variables

| Name | Description | Default | Required |
|------|-------------|---------|----------|
| KAFKA_WEBSOCKET_SERVER_TOPICS | comma seperated list of topic names | NULL | True |
| KAFKA_WEBSOCKET_SERVER_BROKER_URL | location of broker | NULL | True |
| KAFKA_WEBSOCKET_SERVER_PORT | port to expose for websocket connections | "8080" | False |
| KAFKA_WEBSOCKET_SERVER_PREFIX | prefix for websocket endpoints | "" | False |
| KAFKA_WEBSOCKET_SERVER_HEALTH_PORT | port to expose for websocket health checks | "5001" | False |
