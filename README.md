# CloudEvents demo

## Getting started
```sh
docker-compose up -d

go run ./cmd/receiver

# new terminal
go run ./cmd/sender
```

Look at data received by receiver

## Clean up
```sh
docker-compose down
```