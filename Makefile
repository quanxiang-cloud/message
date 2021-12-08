letter:
  dapr run -d /Users/went/Documents/dapr/message --app-id sub1 -p 8080 go run cmd/letter/main.go
manager:
  dapr run -d /Users/went/Documents/dapr/message --app-id sender -p 80 -- go run cmd/message/main.go --config ./configs/config.yml --pubsubName message-kafka-pubsub